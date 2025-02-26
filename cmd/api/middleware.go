package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MohammadBohluli/social-app-go/pkg"
	"github.com/MohammadBohluli/social-app-go/types"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			pkg.UnAuthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			pkg.UnAuthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			pkg.UnAuthorizedErrorResponse(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			pkg.UnAuthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, types.ID(userID))
		if err != nil {
			pkg.UnAuthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			// read the auth header
// 			authHeader := r.Header.Get("Authorization")
// 			if authHeader == "" {
// 				pkg.UnAuthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
// 				return
// 			}

// 			// parse it -> get the base64
// 			parts := strings.Split(authHeader, " ")
// 			if len(parts) != 2 || parts[0] != "Basic" {
// 				return
// 			}

// 			// decode it
// 			decoded, err := base64.StdEncoding.DecodeString(parts[1])
// 			if err != nil {
// 				pkg.UnAuthorizedErrorResponse(w, r, err)
// 				return
// 			}

// 			// check the credentials
// 			username := app.config.auth.basic.username
// 			pass := app.config.auth.basic.password

// 			creds := strings.SplitN(string(decoded), ":", 2)
// 			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
// 				pkg.UnAuthorizedErrorResponse(w, r, fmt.Errorf("invalid credentials"))
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
