package controller

import (
	"log"
	"net/http"
	"song-match-backend/bootstrap"
	"song-match-backend/domain"
	"song-match-backend/internal/jsonutil"
)

type TrackController struct {
	TrackUsecase domain.TrackUseCase
	Env          *bootstrap.Env
}

func (tc *TrackController) FindTrack(w http.ResponseWriter, r *http.Request) {
	var request domain.FindTrackRequest

	track, err := tc.TrackUsecase.Match(r.Context(), request.Content)
	if err != nil {
		jsonutil.JsonResponse(w, http.StatusNotFound, domain.ErrorResponse{Message: "No matches found for this track"})
		return
	}

	loginResponse := domain.FindTrackResponse{
		Name: track.Name,
		Url:  track.Url,
	}

	err = jsonutil.JsonResponse(w, http.StatusOK, loginResponse)
	if err != nil {
		jsonutil.JsonResponse(w, http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
}

func (tc *TrackController) GetMany(w http.ResponseWriter, r *http.Request) {

	tracks, err := tc.TrackUsecase.GetMany(r.Context())
	if err != nil {
		log.Println(err)
		jsonutil.JsonResponse(w, http.StatusInternalServerError, domain.ErrorResponse{Message: "Unexpected error"})
		return
	}

	loginResponse := domain.GetTracksResponse{
		Tracks: []domain.FindTrackResponse{},
	}

	for _, track := range tracks {
		loginResponse.Tracks = append(loginResponse.Tracks, domain.FindTrackResponse{
			Name:    track.Name,
			Url:     track.Url,
			Matches: track.Matches,
		})
	}

	err = jsonutil.JsonResponse(w, http.StatusOK, loginResponse)
	if err != nil {
		jsonutil.JsonResponse(w, http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
}
