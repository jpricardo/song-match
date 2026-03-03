package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FindTrackMatchesRequest struct {
	Content []byte `json:"content"`
}

type TrackDTO struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Matches int    `json:"matches"`
}

type FindTrackMatchesResponse struct {
	Matches []TrackDTO `json:"matches"`
}

type GetTracksResponse struct {
	Tracks []TrackDTO `json:"tracks"`
}

type TrackUseCase interface {
	FindMatches(c context.Context, content []byte) ([]Track, error)
	GetMany(c context.Context) ([]Track, error)
}

const (
	CollectionTrack = "tracks"
)

type Track struct {
	ID      primitive.ObjectID `bson:"_id"`
	Name    string             `bson:"name"`
	Url     string             `bson:"url"`
	Matches int                `bson:"matches"`
}

type TrackRepository interface {
	Create(c context.Context, track *Track) error
	Fetch(c context.Context) ([]Track, error)
	GetByID(c context.Context, id string) (Track, error)
}
