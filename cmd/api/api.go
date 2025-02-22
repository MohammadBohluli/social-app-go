package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MohammadBohluli/social-app-go/docs"
	"github.com/MohammadBohluli/social-app-go/internal/mailer"
	"github.com/MohammadBohluli/social-app-go/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type authConfig struct {
	basic basicConfig
}

type basicConfig struct {
	username string
	password string
}
type mailConfig struct {
	exp       time.Duration
	sendGrid  sendGridConfig
	fromEmail string
}

type sendGridConfig struct {
	apiKey string
}
type config struct {
	addr        string
	db          dbConfig
	apiUrl      string
	mail        mailConfig
	frontEndURL string
	evn         string
	auth        authConfig
}

func (app application) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {

		r.With(app.BasicAuthMiddleware()).Get("/health-check", app.healthCheckHandler)

		// swagger
		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Get("/{postID}", app.getPostHandler)
			r.Delete("/{postID}", app.deletePostHandler)
			r.Patch("/{postID}", app.updatePostHandler)
		})

		r.Route("/users", func(r chi.Router) {

			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)

			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)

		})
	})

	return r
}

func (app *application) start(mux http.Handler) error {

	// documents swagger
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiUrl
	docs.SwaggerInfo.BasePath = "/v1"

	server := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infof("✅ Starting server on http://localhost%s", server.Addr)

	return http.ListenAndServe(server.Addr, server.Handler)
}
