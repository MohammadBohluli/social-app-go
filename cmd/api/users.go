package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/MohammadBohluli/social-app-go/internal/store"
	"github.com/MohammadBohluli/social-app-go/pkg"
	"github.com/MohammadBohluli/social-app-go/types"
	"github.com/go-chi/chi/v5"
)

func (app application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
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

	if err := pkg.JsonResponse(w, http.StatusOK, user); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}
