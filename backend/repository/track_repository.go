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

	_, err := collection.InsertOne(c, track)

	return err
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
	if tracks == nil {
		return []domain.Track{}, err
	}

	return tracks, err
}

func (tr *trackRepository) GetByID(c context.Context, id string) (domain.Track, error) {
	collection := tr.database.Collection(tr.collection)

	var track domain.Track

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return track, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&track)
	return track, err
}
