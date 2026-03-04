package controller

import (
	"fmt"
	"io"
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

	content, err := io.ReadAll(r.Body)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received %d bytes\n", len(content))

	matches, err := tc.TrackUsecase.FindMatches(r.Context(), content)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rd := domain.FindTrackMatchesResponse{
		Matches: []domain.TrackDTO{},
	}

	for _, track := range matches {
		fp := []domain.FingerprintDTO{}

		for _, fingerprint := range track.Fingerprints {
			fp = append(fp, domain.FingerprintDTO{Timestamp: fingerprint.Timestamp, Peaks: fingerprint.Peaks})
		}

		rd.Matches = append(rd.Matches, domain.TrackDTO{
			ID:           track.ID,
			Name:         track.Name,
			Url:          track.Url,
			Thumbnail:    track.Thumbnail,
			Matches:      track.Matches,
			Fingerprints: fp,
		})
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
		fp := []domain.FingerprintDTO{}

		for _, fingerprint := range track.Fingerprints {
			fp = append(fp, domain.FingerprintDTO{Timestamp: fingerprint.Timestamp, Peaks: fingerprint.Peaks})
		}

		rd.Tracks = append(rd.Tracks, domain.TrackDTO{
			ID:           track.ID,
			Name:         track.Name,
			Url:          track.Url,
			Thumbnail:    track.Thumbnail,
			Matches:      track.Matches,
			Fingerprints: fp,
		})
	}

	err = jsonutil.JsonSuccessResponse(w, http.StatusOK, rd)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (tc *TrackController) AddTrack(w http.ResponseWriter, r *http.Request) {
	var payload domain.AddTrackPayload

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		log.Println(err)
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	track, err := tc.TrackUsecase.AddTrack(r.Context(), payload.Url)
	if err != nil {
		log.Println(err)
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	fp := []domain.FingerprintDTO{}

	for _, fingerprint := range track.Fingerprints {
		fp = append(fp, domain.FingerprintDTO{Timestamp: fingerprint.Timestamp, Peaks: fingerprint.Peaks})
	}

	rd := domain.AddTrackResponse{
		ID:           track.ID,
		Name:         track.Name,
		Url:          track.Url,
		Thumbnail:    track.Thumbnail,
		Matches:      track.Matches,
		Fingerprints: fp,
	}

	err = jsonutil.JsonSuccessResponse(w, http.StatusOK, rd)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
