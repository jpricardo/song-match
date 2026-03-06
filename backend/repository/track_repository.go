package repository

import (
	"context"
	"song-match-backend/domain"
	"song-match-backend/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type trackRepository struct {
	database   mongo.Database
	collection string
}

func NewTrackRepository(db mongo.Database, collection string) domain.TrackRepository {
	return &trackRepository{
		database:   db,
		collection: collection,
	}
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
