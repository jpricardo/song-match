package controller

import (
	"io"
	"log/slog"
	"net/http"
	"song-match-backend/bootstrap"
	"song-match-backend/domain"
	"song-match-backend/internal/jsonutil"

	"github.com/go-chi/chi/v5"
)

// maxAudioBodyBytes caps incoming WAV bodies at 50 MB.
// A 3-minute stereo 44100 Hz 16-bit WAV is ~30 MB, so this is generous
// while still protecting against unbounded memory allocation.
const maxAudioBodyBytes = 50 * 1024 * 1024

type TrackController struct {
	TrackUsecase domain.TrackUseCase
	Env          *bootstrap.Env
}

func (tc *TrackController) FindMatches(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAudioBodyBytes)

	content, err := io.ReadAll(r.Body)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusRequestEntityTooLarge, "audio payload too large or unreadable")
		return
	}
	defer r.Body.Close()

	slog.Info("match request received", "bytes", len(content))

	matches, err := tc.TrackUsecase.FindMatches(r.Context(), content)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rd := domain.FindTrackMatchesResponse{
		Matches: []domain.TrackDTO{},
	}

	for _, track := range matches {
		rd.Matches = append(rd.Matches, domain.TrackDTO{
			ID:           track.ID,
			Name:         track.Name,
			Url:          track.Url,
			Thumbnail:    track.Thumbnail,
			Status:       track.Status,
			Fingerprints: []domain.FingerprintDTO{},
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
		slog.Error("GetMany failed", "error", err)
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
			Status:       track.Status,
			Fingerprints: fp,
		})
	}

	err = jsonutil.JsonSuccessResponse(w, http.StatusOK, rd)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (tc *TrackController) GetById(w http.ResponseWriter, r *http.Request) {
	trackId := chi.URLParam(r, "trackId")

	track, err := tc.TrackUsecase.GetByID(r.Context(), trackId)
	if err != nil {
		slog.Error("GetByID failed", "track_id", trackId, "error", err)
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	fp := []domain.FingerprintDTO{}
	for _, fingerprint := range track.Fingerprints {
		fp = append(fp, domain.FingerprintDTO{Timestamp: fingerprint.Timestamp, Peaks: fingerprint.Peaks})
	}

	rd := domain.TrackDTO{
		ID:           track.ID,
		Name:         track.Name,
		Url:          track.Url,
		Thumbnail:    track.Thumbnail,
		Status:       track.Status,
		Fingerprints: fp,
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
		slog.Error("AddTrack: failed to parse payload", "error", err)
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	track, err := tc.TrackUsecase.AddTrack(r.Context(), payload.Url)
	if err != nil {
		slog.Error("AddTrack failed", "url", payload.Url, "error", err)
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
		Status:       track.Status,
		Fingerprints: fp,
	}

	err = jsonutil.JsonSuccessResponse(w, http.StatusAccepted, rd)
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (tc *TrackController) DeleteById(w http.ResponseWriter, r *http.Request) {
	trackId := chi.URLParam(r, "trackId")

	err := tc.TrackUsecase.DeleteByID(r.Context(), trackId)
	if err != nil {
		slog.Error("DeleteByID failed", "track_id", trackId, "error", err)
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	err = jsonutil.JsonSuccessResponse(w, http.StatusOK, "")
	if err != nil {
		jsonutil.JsonErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
