package pkg

import (
	"net/http"

	"log"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌Internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func ForbiddenErrorResponse(w http.ResponseWriter, r *http.Request) {
	// must be inject logger
	log.Printf("❌forbidden error: %s path: %s error: %s", r.Method, r.URL.Path, "error")
	WriteJsonError(w, http.StatusForbidden, "forbidden")
}

func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌Bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusBadRequest, err.Error())
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌Not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusNotFound, "Not found")
}

func ConflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌conflict error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusConflict, err.Error())
}

func UnAuthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("unauthorized error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func UnAuthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("unauthorized basic error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted, charset="UTF-8"`)
	WriteJsonError(w, http.StatusUnauthorized, "unauthorized basic")
}

func RateLimitExceededErrorResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	log.Printf("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	WriteJsonError(w, http.StatusUnauthorized, "rate limit exceeded, retry after: ")
}
