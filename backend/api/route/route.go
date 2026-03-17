package route

import (
	"net/http"
	"song-match-backend/api/controller"
	"song-match-backend/bootstrap"
	"song-match-backend/mongo"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database) (http.Handler, *controller.TrackController) {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.RequestID)
	mux.Use(slogRequestLogger())
	mux.Use(middleware.Heartbeat("/ping"))

	tc := NewTrackRouter(env, timeout, db, mux)

	return mux, tc
}
