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

type CreatePostRequest struct {
	Content string   `json:"content"`
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
	Content *string `json:"content"`
	Title   *string `json:"title"`
}

func (app application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var postReq CreatePostRequest
	if err := pkg.ReadJson(w, r, &postReq); err != nil {
		pkg.BadRequestError(w, r, err)
		return
	}
	// fake data
	userID := 1

	post := store.Post{
		UserID:  types.ID(userID),
		Title:   postReq.Title,
		Content: postReq.Content,
		Tags:    postReq.Tags,
	}

	if err := app.store.Posts.Create(r.Context(), &post); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

	if err := pkg.JsonResponse(w, http.StatusCreated, post); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}

func (app application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	// fake data
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

	post, err := app.store.Posts.GetByID(ctx, types.ID(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			pkg.NotFoundError(w, r, err)
		default:
			pkg.InternalServerError(w, r, err)
		}
		return
	}

	comments, err := app.store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := pkg.JsonResponse(w, http.StatusOK, post); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}

func (app application) updatePostHandler(w http.ResponseWriter, r *http.Request) {

	var postReq UpdatePostRequest
	if err := pkg.ReadJson(w, r, &postReq); err != nil {
		pkg.BadRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

	post, err := app.store.Posts.GetByID(ctx, types.ID(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			pkg.NotFoundError(w, r, err)
		default:
			pkg.InternalServerError(w, r, err)
		}
		return
	}

	if postReq.Content != nil {
		post.Content = *postReq.Content
	}
	if postReq.Title != nil {
		post.Title = *postReq.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

	if err := pkg.JsonResponse(w, http.StatusCreated, post); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}

func (app application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

	if err := app.store.Posts.Delete(ctx, types.ID(id)); err != nil {

		switch {
		case errors.Is(err, store.ErrorNotFound):
			pkg.NotFoundError(w, r, err)
		default:
			pkg.InternalServerError(w, r, err)
		}
		return

	}

	w.WriteHeader(http.StatusNoContent)
}
