package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackFingerprint struct {
	Timestamp float64
	Peaks     []int
}

type FingerprintDTO struct {
	Timestamp float64 `json:"timestamp"`
	Peaks     []int   `json:"peaks"`
}

type TrackDTO struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Url          string             `json:"url"`
	Thumbnail    string             `json:"thumbnail,omitempty"`
	Matches      int                `json:"matches"`
	Fingerprints []FingerprintDTO   `json:"fingerprints"`
}

type FindTrackMatchesResponse struct {
	Matches []TrackDTO `json:"matches"`
}

type GetTracksResponse struct {
	Tracks []TrackDTO `json:"tracks"`
}

type AddTrackPayload struct {
	Url string `json:"url"`
}

type AddTrackResponse TrackDTO

type TrackUseCase interface {
	FindMatches(c context.Context, content []byte) ([]Track, error)
	GetMany(c context.Context) ([]Track, error)
	GetByID(c context.Context, id string) (Track, error)
	DeleteByID(c context.Context, id string) error
	AddTrack(c context.Context, url string) (*Track, error)
}

const (
	CollectionTrack = "tracks"
)

type Track struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Url          string             `bson:"url"`
	Thumbnail    string             `bson:"thumbnail,omitempty"`
	Matches      int                `bson:"matches"`
	Fingerprints []TrackFingerprint `bson:"fingerprints"`
}

type TrackRepository interface {
	Create(c context.Context, track *Track) error
	Fetch(c context.Context) ([]Track, error)
	DeleteByID(c context.Context, id string) error
	GetByID(c context.Context, id string) (Track, error)
}
