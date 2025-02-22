package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/MohammadBohluli/social-app-go/pkg"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				pkg.UnAuthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				pkg.UnAuthorizedErrorResponse(w, r, err)
				return
			}

			// check the credentials
			username := app.config.auth.basic.username
			pass := app.config.auth.basic.password

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				pkg.UnAuthorizedErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
