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

func (tu *trackUsecase) Match(c context.Context, content []byte) (domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	// TODO - Track processing / lookup
	id := ""
	// TODO - Matches++
	return tu.trackRepository.GetByID(ctx, id)
}

func (tu *trackUsecase) GetMany(c context.Context) ([]domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.Fetch(ctx)
}
