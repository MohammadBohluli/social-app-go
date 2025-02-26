package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MohammadBohluli/social-app-go/internal/store"
	"github.com/MohammadBohluli/social-app-go/pkg"
	"github.com/MohammadBohluli/social-app-go/types"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

type FollowUserRequest struct {
	UserID types.ID `json:"user_id"`
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	store.User
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromContext(r)

	if err := pkg.JsonResponse(w, http.StatusOK, user); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	token := chi.URLParam(r, "token")
	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			pkg.BadRequestError(w, r, err)
		default:
			pkg.InternalServerError(w, r, err)
		}
		return
	}

	if err := pkg.JsonResponse(w, http.StatusNoContent, ""); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)
	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		pkg.BadRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.Follow(ctx, followerUser.ID, types.ID(followedID)); err != nil {

		switch err {
		case store.ErrorConflict:
			pkg.ConflictErrorResponse(w, r, err)
			return
		default:
			pkg.InternalServerError(w, r, err)
			return

		}
	}

}

// unfollowUser godoc
//
//	@Summary		unfollows a user
//	@Description	unfollows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	followerUser := getUserFromContext(r)
	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		pkg.BadRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.UnFollow(ctx, followerUser.ID, types.ID(unfollowedID)); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

}

func (app application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(id, 10, 64)
		fmt.Println("id: ", id)
		fmt.Println("userID: ", userID)
		if err != nil {
			pkg.BadRequestError(w, r, err)
			return
		}

		user, err := app.store.Users.GetByID(ctx, types.ID(userID))
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				pkg.BadRequestError(w, r, err)
				return
			default:
				pkg.InternalServerError(w, r, err)
				return
			}

		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
