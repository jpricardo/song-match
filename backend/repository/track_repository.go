package repository

import (
	"context"
	"log"
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
		log.Printf("Warning: Failed to create index on hashes collection: %v\n", err)
	} else {
		log.Println("Successfully verified index on hash_value.")
	}

	return tr
}

func (tr *trackRepository) Create(c context.Context, track *domain.Track) error {
	collection := tr.database.Collection(tr.collection)

	// Insert the Track
	id, err := collection.InsertOne(c, track)
	if err != nil {
		return err
	}

	if oid, ok := id.(primitive.ObjectID); ok {
		track.ID = oid
	}

	// Bulk Insert the Fingerprints into their own collection
	if len(track.Fingerprints) > 0 {
		fpColl := tr.database.Collection(domain.CollectionFingerprint)

		var docs []interface{}
		for i := range track.Fingerprints {
			track.Fingerprints[i].TrackID = track.ID
			docs = append(docs, track.Fingerprints[i])
		}

		_, err = fpColl.InsertMany(c, docs)
		if err != nil {
			return err
		}
	}

	// The same for the hashes
	if len(track.Hashes) > 0 {
		hashColl := tr.database.Collection(domain.CollectionHashes)

		var docs []interface{}
		for i := range track.Hashes {
			track.Hashes[i].TrackID = track.ID
			docs = append(docs, track.Hashes[i])
		}

		_, err = hashColl.InsertMany(c, docs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tr *trackRepository) Fetch(c context.Context) ([]domain.Track, error) {
	collection := tr.database.Collection(tr.collection)

	opts := options.Find()
	cursor, err := collection.Find(c, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var tracks []domain.Track
	err = cursor.All(c, &tracks)
	if err != nil || tracks == nil {
		return []domain.Track{}, err
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
	err = cursor.All(c, &fingerprints)
	if err != nil || fingerprints == nil {
		return []domain.TrackFingerprint{}, err
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
	return nil
}

func (tr *trackRepository) GetMatchingHashes(c context.Context, hashValues []string) ([]domain.AudioHash, error) {
	collection := tr.database.Collection(domain.CollectionHashes)

	filter := bson.M{"hash_value": bson.M{"$in": hashValues}}

	cursor, err := collection.Find(c, filter)
	if err != nil {
		return nil, err
	}

	var hashes []domain.AudioHash
	err = cursor.All(c, &hashes)
	if err != nil || hashes == nil {
		return []domain.AudioHash{}, err
	}

	return hashes, nil
}

func (tr *trackRepository) UpdateTrackData(c context.Context, track *domain.Track) error {
	collection := tr.database.Collection(tr.collection)

	// Update the track document itself (name, thumbnail, status)
	update := bson.M{
		"$set": bson.M{
			"name":      track.Name,
			"thumbnail": track.Thumbnail,
			"status":    track.Status,
		},
	}
	_, err := collection.UpdateOne(c, bson.M{"_id": track.ID}, update)
	if err != nil {
		return err
	}

	// Bulk Insert the Fingerprints
	if len(track.Fingerprints) > 0 {
		fpColl := tr.database.Collection(domain.CollectionFingerprint)
		var docs []interface{}
		for i := range track.Fingerprints {
			track.Fingerprints[i].TrackID = track.ID
			docs = append(docs, track.Fingerprints[i])
		}
		_, err = fpColl.InsertMany(c, docs)
		if err != nil {
			return err
		}
	}

	// Bulk Insert the Hashes
	if len(track.Hashes) > 0 {
		hashColl := tr.database.Collection(domain.CollectionHashes)
		var docs []interface{}
		for i := range track.Hashes {
			track.Hashes[i].TrackID = track.ID
			docs = append(docs, track.Hashes[i])
		}
		_, err = hashColl.InsertMany(c, docs)
		if err != nil {
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
