package usecase

import (
	"context"
	"fmt"
	"song-match-backend/domain"
	"song-match-backend/internal/audioutil"
	"song-match-backend/internal/ytbutil"
	"sort"
	"time"
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

	// Generate Hashes for the incoming sample
	sampleHashes := audioutil.GenerateHashes(fingerprints)

	// Map format: [HashValue] -> Array of times it occurred
	sampleHashMap := make(map[string][]float64)
	for _, h := range sampleHashes {
		sampleHashMap[h.HashValue] = append(sampleHashMap[h.HashValue], h.Time)
	}

	tracks, err := tu.trackRepository.Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tracks: %w", err)
	}

	type trackScore struct {
		track domain.Track
		score int
	}
	var scores []trackScore

	// Score each track using a Time-Offset Histogram
	for _, track := range tracks {
		fps, err := tu.trackRepository.GetFingerprintsByID(ctx, track.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch fingerprints: %w", err)
		}

		if len(fps) == 0 {
			continue
		}

		dbHashes := audioutil.GenerateHashes(fps)
		offsets := make(map[int]int)

		for _, dbHash := range dbHashes {
			if sampleTimes, exists := sampleHashMap[dbHash.HashValue]; exists {
				// Calculate how far apart they are to ensure the audio lines up
				for _, sampleTime := range sampleTimes {

					// Multiply by 10 to group offsets into 100ms bins.
					offsetBin := int((dbHash.Time - sampleTime) * 10)
					offsets[offsetBin]++
				}
			}
		}

		// The track's final score is the height of the tallest peak in the histogram
		maxScore := 0
		for _, count := range offsets {
			if count > maxScore {
				maxScore = count
			}
		}

		// Filter out random noise. A real match usually has a score of 15+.
		if maxScore > 5 {
			scores = append(scores, trackScore{track: track, score: maxScore})
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	var matchedTracks []domain.Track
	for _, s := range scores {
		matchedTracks = append(matchedTracks, s.track)
	}

	fmt.Printf("Best match: %s\n", matchedTracks[0].Name)

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

	t := &domain.Track{
		Name:         title,
		Url:          url,
		Thumbnail:    thumbnail,
		Matches:      0,
		Fingerprints: fingerprints,
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
