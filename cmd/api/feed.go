package main

import (
	"net/http"

	"github.com/MohammadBohluli/social-app-go/pkg"
	"github.com/MohammadBohluli/social-app-go/types"
)

func (app application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, types.ID(6))
	if err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

	if err := pkg.JsonResponse(w, http.StatusOK, feed); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}
