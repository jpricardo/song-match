package usecase

import (
	"context"
	"song-match-backend/domain"
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

	// TODO - Track processing / lookup
	m := []domain.Track{}
	return m, nil
}

func (tu *trackUsecase) GetMany(c context.Context) ([]domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.Fetch(ctx)
}
