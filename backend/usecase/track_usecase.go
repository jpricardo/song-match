package usecase

import (
	"context"
	"fmt"
	"song-match-backend/domain"
	"song-match-backend/internal/audioutil"
	"song-match-backend/internal/ytbutil"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type trackUsecase struct {
	trackRepository domain.TrackRepository
	contextTimeout  time.Duration
}

func NewTrackUsecase(trackRepository domain.TrackRepository, timeout time.Duration) domain.TrackUseCase {
	return &trackUsecase{
		trackRepository: trackRepository,
		contextTimeout:  timeout,
	}
}

func (tu *trackUsecase) FindMatches(c context.Context, content []byte) ([]domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	samples, sampleRate, err := audioutil.DecodeAudio(content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio: %w", err)
	}

	fingerprints, err := audioutil.ExtractFingerprints(samples, sampleRate)
	if err != nil {
		return nil, fmt.Errorf("failed to extract fingerprints: %w", err)
	}

	sampleHashes := audioutil.GenerateHashes(fingerprints)

	var sampleHashValues []string
	sampleHashMap := make(map[string][]float64)

	for _, h := range sampleHashes {
		sampleHashValues = append(sampleHashValues, h.HashValue)
		sampleHashMap[h.HashValue] = append(sampleHashMap[h.HashValue], h.Time)
	}

	matchedDBHashes, err := tu.trackRepository.GetMatchingHashes(ctx, sampleHashValues)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch matching hashes: %w", err)
	}

	// Group the matched hashes by Track ID so we can build a histogram per track
	trackHistograms := make(map[primitive.ObjectID]map[int]int)

	for _, dbHash := range matchedDBHashes {
		if _, exists := trackHistograms[dbHash.TrackID]; !exists {
			trackHistograms[dbHash.TrackID] = make(map[int]int)
		}

		if sampleTimes, exists := sampleHashMap[dbHash.HashValue]; exists {
			for _, sampleTime := range sampleTimes {
				offsetBin := int((dbHash.Time - sampleTime) * 10)
				trackHistograms[dbHash.TrackID][offsetBin]++
			}
		}
	}

	// Score each track based on its histogram
	type trackScore struct {
		trackID primitive.ObjectID
		score   int
	}
	var scores []trackScore

	for trackID, offsets := range trackHistograms {
		maxScore := 0
		for _, count := range offsets {
			if count > maxScore {
				maxScore = count
			}
		}

		if maxScore > 5 {
			scores = append(scores, trackScore{trackID: trackID, score: maxScore})
		}
	}

	// Sort highest score first
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Fetch ONLY the tracks that successfully matched
	var matchedTracks []domain.Track
	for _, s := range scores {
		track, err := tu.trackRepository.GetByID(ctx, s.trackID.Hex())
		if err == nil {
			matchedTracks = append(matchedTracks, track)
		}
	}

	if len(matchedTracks) > 0 {
		fmt.Printf("Best match: %s\n", matchedTracks[0].Name)
	} else {
		fmt.Println("Couldn't find a match!")
	}

	return matchedTracks, nil
}

func (tu *trackUsecase) GetMany(c context.Context) ([]domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.Fetch(ctx)
}

func (tu *trackUsecase) GetByID(c context.Context, id string) (domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.GetByID(ctx, id)
}

func (tu *trackUsecase) AddTrack(c context.Context, url string) (*domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	wav, title, thumbnail, err := ytbutil.DownloadTrack(url)
	if err != nil {
		return nil, err
	}

	samples, sampleRate, err := audioutil.DecodeAudio(wav)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio: %w", err)
	}

	fingerprints, err := audioutil.ExtractFingerprints(samples, sampleRate)
	if err != nil {
		return nil, fmt.Errorf("failed to extract fingerprints: %w", err)
	}

	hashes := audioutil.GenerateHashes(fingerprints)

	t := &domain.Track{
		Name:         title,
		Url:          url,
		Thumbnail:    thumbnail,
		Fingerprints: fingerprints,
		Hashes:       hashes,
	}
	err = tu.trackRepository.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (tu *trackUsecase) DeleteByID(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.DeleteByID(ctx, id)
}
