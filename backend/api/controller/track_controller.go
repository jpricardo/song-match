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

func (tc *TrackController) FindMatches(w http.ResponseWriter, r *http.Request) {
	var request domain.FindTrackMatchesRequest

	matches, err := tc.TrackUsecase.FindMatches(r.Context(), request.Content)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusNotFound, "No matches found for this track")
		return
	}

	rd := domain.FindTrackMatchesResponse{
		Matches: []domain.TrackDTO{},
	}

	for _, track := range matches {
		rd.Matches = append(rd.Matches, domain.TrackDTO{Name: track.Name, Url: track.Url})

	}

	err = jsonutil.JsonSuccessResponse(w, http.StatusOK, rd)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (tc *TrackController) GetMany(w http.ResponseWriter, r *http.Request) {

	tracks, err := tc.TrackUsecase.GetMany(r.Context())
	if err != nil {
		log.Println(err)
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	rd := domain.GetTracksResponse{
		Tracks: []domain.TrackDTO{},
	}

	for _, track := range tracks {
		rd.Tracks = append(rd.Tracks, domain.TrackDTO{
			Name:    track.Name,
			Url:     track.Url,
			Matches: track.Matches,
		})
	}

	err = jsonutil.JsonSuccessResponse(w, http.StatusOK, rd)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
