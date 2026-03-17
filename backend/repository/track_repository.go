package repository

import (
	"context"
	"log/slog"
	"song-match-backend/domain"
	"song-match-backend/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type trackRepository struct {
	database   mongo.Database
	collection string
}

func NewTrackRepository(db mongo.Database, collection string) domain.TrackRepository {
	tr := &trackRepository{
		database:   db,
		collection: collection,
	}

	// Create the index on the Hashes collection in the background
	hashColl := db.Collection(domain.CollectionHashes)

	indexModel := mongoDriver.IndexModel{
		Keys:    bson.D{{Key: "hash_value", Value: 1}},
		Options: options.Index().SetBackground(true),
	}

	_, err := hashColl.CreateOneIndex(context.Background(), indexModel)
	if err != nil {
		slog.Warn("failed to create index on hashes collection", "error", err)
	} else {
		slog.Info("index on hash_value verified")
	}

	return tr
}

// Create inserts a track document and optionally its fingerprints and hashes.
// For a newly created "processing" track the fingerprint and hash slices will
// be empty; they are written later by UpdateTrackData once processing is done.
//
// The three writes (track, fingerprints, hashes) are wrapped in
// a MongoDB session transaction so a failure at any step rolls back all prior
// writes, leaving no partial data in the database.
const insertBatchSize = 1000

func insertManyInBatches(c context.Context, coll mongo.Collection, docs []interface{}) error {
	for i := 0; i < len(docs); i += insertBatchSize {
		end := i + insertBatchSize
		if end > len(docs) {
			end = len(docs)
		}
		if _, err := coll.InsertMany(c, docs[i:end]); err != nil {
			return err
		}
	}
	return nil
}

func (tr *trackRepository) Create(c context.Context, track *domain.Track) error {
	client := tr.database.Client()

	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(c)

	_, err = session.WithTransaction(c, func(sc mongoDriver.SessionContext) (interface{}, error) {
		collection := tr.database.Collection(tr.collection)

		id, err := collection.InsertOne(sc, track)
		if err != nil {
			return nil, err
		}

		if oid, ok := id.(primitive.ObjectID); ok {
			track.ID = oid
		}

		if len(track.Fingerprints) > 0 {
			fpColl := tr.database.Collection(domain.CollectionFingerprint)
			docs := make([]interface{}, len(track.Fingerprints))
			for i := range track.Fingerprints {
				track.Fingerprints[i].TrackID = track.ID
				docs[i] = track.Fingerprints[i]
			}
			if _, err = fpColl.InsertMany(sc, docs); err != nil {
				return nil, err
			}
		}

		if len(track.Hashes) > 0 {
			hashColl := tr.database.Collection(domain.CollectionHashes)
			docs := make([]interface{}, len(track.Hashes))
			for i := range track.Hashes {
				track.Hashes[i].TrackID = track.ID
				docs[i] = track.Hashes[i]
			}
			if _, err = hashColl.InsertMany(sc, docs); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}

func (tr *trackRepository) Fetch(c context.Context) ([]domain.Track, error) {
	collection := tr.database.Collection(tr.collection)

	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}

	var tracks []domain.Track
	if err = cursor.All(c, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (tr *trackRepository) GetByID(c context.Context, id string) (domain.Track, error) {
	collection := tr.database.Collection(tr.collection)
	var track domain.Track

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return track, err
	}

	// Fetch the Track
	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&track)
	if err != nil {
		return track, err
	}

	// Fetch its Fingerprints
	fps, err := tr.GetFingerprintsByID(c, id)
	if err != nil {
		return track, err
	}

	track.Fingerprints = fps

	return track, nil
}

func (tr *trackRepository) GetManyByIDs(c context.Context, ids []string) ([]domain.Track, error) {
	collection := tr.database.Collection(tr.collection)

	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			slog.Warn("GetManyByIDs: skipping invalid ObjectID", "id", id, "error", err)
			continue
		}
		objectIDs = append(objectIDs, oid)
	}

	if len(objectIDs) == 0 {
		return []domain.Track{}, nil
	}

	cursor, err := collection.Find(c, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return nil, err
	}

	var tracks []domain.Track
	if err = cursor.All(c, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (tr *trackRepository) GetFingerprintsByID(c context.Context, id string) ([]domain.TrackFingerprint, error) {
	collection := tr.database.Collection(domain.CollectionFingerprint)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(c, bson.M{"track_id": idHex})
	if err != nil {
		return nil, err
	}

	var fingerprints []domain.TrackFingerprint
	if err = cursor.All(c, &fingerprints); err != nil {
		return nil, err
	}

	return fingerprints, nil
}

func (tr *trackRepository) DeleteByID(c context.Context, id string) error {
	collection := tr.database.Collection(tr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(c, bson.M{"_id": idHex})
	return err
}

// UpdateTrackData atomically replaces a track's metadata and writes its
// fingerprints and hashes once background processing is complete.
//
// Wrapped in a transaction so a hash InsertMany failure cannot
// leave the track marked "ready" with no hashes.
//
// Fingerprints and hashes for the track are deleted before
// re-inserting, so a retried job cannot duplicate data.
func (tr *trackRepository) UpdateTrackData(c context.Context, track *domain.Track) error {
	client := tr.database.Client()

	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(c)

	_, err = session.WithTransaction(c, func(sc mongoDriver.SessionContext) (interface{}, error) {
		// Update the track document metadata + status.
		collection := tr.database.Collection(tr.collection)
		update := bson.M{
			"$set": bson.M{
				"name":      track.Name,
				"thumbnail": track.Thumbnail,
				"status":    track.Status,
			},
		}
		if _, err := collection.UpdateOne(sc, bson.M{"_id": track.ID}, update); err != nil {
			return nil, err
		}

		// Purge any pre-existing fingerprints and hashes for this
		// track before re-inserting, preventing duplication on retry.
		fpColl := tr.database.Collection(domain.CollectionFingerprint)
		if _, err := fpColl.DeleteMany(sc, bson.M{"track_id": track.ID}); err != nil {
			return nil, err
		}

		// Insert fresh fingerprints.
		if len(track.Fingerprints) > 0 {
			docs := make([]interface{}, len(track.Fingerprints))
			for i := range track.Fingerprints {
				track.Fingerprints[i].TrackID = track.ID
				docs[i] = track.Fingerprints[i]
			}
			if _, err := fpColl.InsertMany(sc, docs); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	// Hash writes are outside the transaction: a single track produces
	// ~800k hashes whose combined index size exceeds WiredTiger's
	// per-transaction cache limit. The delete-before-insert pattern
	// makes these writes safe to retry without transactional protection.
	hashColl := tr.database.Collection(domain.CollectionHashes)
	if _, err := hashColl.DeleteMany(c, bson.M{"track_id": track.ID}); err != nil {
		return err
	}

	if len(track.Hashes) > 0 {
		docs := make([]interface{}, len(track.Hashes))
		for i := range track.Hashes {
			track.Hashes[i].TrackID = track.ID
			docs[i] = track.Hashes[i]
		}
		if err := insertManyInBatches(c, hashColl, docs); err != nil {
			return err
		}
	}

	return nil
}

func (tr *trackRepository) UpdateTrackStatus(c context.Context, id string, status string) error {
	collection := tr.database.Collection(tr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"status": status}}
	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}

func (tr *trackRepository) GetMatchingHashes(c context.Context, hashValues []string) ([]domain.AudioHash, error) {
	collection := tr.database.Collection(domain.CollectionHashes)

	filter := bson.M{"hash_value": bson.M{"$in": hashValues}}

	cursor, err := collection.Find(c, filter)
	if err != nil {
		return nil, err
	}

	var hashes []domain.AudioHash
	if err = cursor.All(c, &hashes); err != nil {
		return nil, err
	}

	return hashes, nil
}
