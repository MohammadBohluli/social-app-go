package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_578 //1mg
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJsonError(w http.ResponseWriter, statusCode int, message string) error {

	type envelope struct {
		Error string `json:"error"`
	}

	return WriteJson(w, statusCode, envelope{Error: message})
}
