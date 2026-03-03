package jsonutil

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"song-match-backend/domain"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	const maxBytes = 1048576 // 1mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("invalid JSON value")
	}

	return nil
}

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
