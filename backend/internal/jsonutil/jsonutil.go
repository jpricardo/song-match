package jsonutil

import (
	"encoding/json"
	"net/http"
	"song-match-backend/domain"
)

func WriteJson(w http.ResponseWriter, status int, response domain.Response) error {
	out, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func JsonSuccessResponse(w http.ResponseWriter, status int, data interface{}) error {
	r := domain.Response{Success: true, Data: data}
	return WriteJson(w, status, r)
}

func JsonErrorResponse(w http.ResponseWriter, status int, message string) error {
	r := domain.Response{Success: false, Message: message}
	return WriteJson(w, status, r)
}
