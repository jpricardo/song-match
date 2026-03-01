package jsonutil

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, status int, data interface{}) error {
	out, err := json.Marshal(data)
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
