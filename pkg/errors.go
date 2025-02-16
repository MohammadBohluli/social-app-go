package pkg

import (
	"net/http"

	"log"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌Internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌Bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusBadRequest, err.Error())
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌Not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusNotFound, "Not found")
}
