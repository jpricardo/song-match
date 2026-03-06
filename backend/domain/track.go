package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FingerprintDTO struct {
	Timestamp float64 `json:"timestamp"`
	Peaks     []int   `json:"peaks"`
}

type TrackDTO struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Url          string             `json:"url"`
	Thumbnail    string             `json:"thumbnail"`
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
	CollectionTrack       = "tracks"
	CollectionFingerprint = "fingerprints"
	CollectionHashes      = "hashes"
)

type TrackFingerprint struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TrackID   primitive.ObjectID `bson:"track_id"`
	Timestamp float64            `bson:"timestamp"`
	Peaks     []int              `bson:"peaks"`
}

type AudioHash struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TrackID   primitive.ObjectID `bson:"track_id"`
	HashValue string             `bson:"hash_value"`
	Time      float64            `bson:"time"`
}

type Track struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Url          string             `bson:"url"`
	Thumbnail    string             `bson:"thumbnail"`
	Fingerprints []TrackFingerprint `bson:"-"`
	Hashes       []AudioHash        `bson:"-"`
}

type TrackRepository interface {
	Create(c context.Context, track *Track) error
	Fetch(c context.Context) ([]Track, error)
	DeleteByID(c context.Context, id string) error
	GetByID(c context.Context, id string) (Track, error)
	GetFingerprintsByID(c context.Context, id string) ([]TrackFingerprint, error)
	GetMatchingHashes(c context.Context, hashValues []string) ([]AudioHash, error)
}
