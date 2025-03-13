package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MohammadBohluli/social-app-go/internal/store"
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
		user, err := app.getUser(ctx, types.ID(userID))
		if err != nil {
			pkg.UnAuthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *application) checkPostOwnerShip(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := 1 //getPostFromContext(r)

		if post == int(user.ID) {
			next.ServeHTTP(w, r)
			return

		}

		allowed, err := app.checkRolePrecedence(r.Context(), user, role)
		if err != nil {
			pkg.InternalServerError(w, r, err)
			return
		}

		if !allowed {
			pkg.ForbiddenErrorResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return types.ID(user.Role.Level) >= types.ID(role.Level), nil
}

func (app application) getUser(ctx context.Context, userID types.ID) (*store.User, error) {
	user, err := app.cacheStorage.Users.Get(ctx, types.ID(userID))
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err := app.store.Users.GetByID(ctx, types.ID(userID))
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}

	}

	return user, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				pkg.RateLimitExceededErrorResponse(w, r, retryAfter.String())
				return
			}
		}
		next.ServeHTTP(w, r)
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
