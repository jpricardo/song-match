package route

import (
	"song-match-backend/api/controller"
	"song-match-backend/bootstrap"
	"song-match-backend/domain"
	"song-match-backend/mongo"
	"song-match-backend/repository"
	"song-match-backend/usecase"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewTrackRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *chi.Mux) *controller.TrackController {
	tr := repository.NewTrackRepository(db, domain.CollectionTrack)
	tc := &controller.TrackController{
		TrackUsecase: usecase.NewTrackUsecase(tr, timeout),
		Env:          env,
	}

	group.Post("/tracks", tc.AddTrack)
	group.Get("/tracks", tc.GetMany)
	group.Get("/tracks/{trackId}", tc.GetById)
	group.Delete("/tracks/{trackId}", tc.DeleteById)
	group.Post("/matches", tc.FindMatches)

	return tc
}
