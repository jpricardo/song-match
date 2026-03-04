package usecase

import (
	"context"
	"fmt"
	"song-match-backend/domain"
	"song-match-backend/internal/audioutil"
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
	_, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	samples, sampleRate, err := audioutil.DecodeAudio(content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio: %w", err)
	}

	fingerprints, err := audioutil.ExtractFingerprints(samples, sampleRate)
	if err != nil {
		return nil, fmt.Errorf("failed to extract fingerprints: %w", err)
	}

	fmt.Printf("Extracted %d fingerprints at %d Hz\n", len(fingerprints), sampleRate)

	// TODO - Track processing / lookup
	m := []domain.Track{}
	return m, nil
}

func (tu *trackUsecase) GetMany(c context.Context) ([]domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.Fetch(ctx)
}
