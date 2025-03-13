package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MohammadBohluli/social-app-go/docs"
	"github.com/MohammadBohluli/social-app-go/internal/auth"
	"github.com/MohammadBohluli/social-app-go/internal/mailer"
	"github.com/MohammadBohluli/social-app-go/internal/ratelimiter"
	"github.com/MohammadBohluli/social-app-go/internal/store"
	"github.com/MohammadBohluli/social-app-go/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	cacheStorage  cache.Storage
	rateLimiter   ratelimiter.Limiter
}

type redisConfig struct {
	host     string
	port     int
	password string
	db       int
	enabled  bool
}
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
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
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}

func (app application) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.RateLimiterMiddleware)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {

		r.Get("/health-check", app.healthCheckHandler)

		// swagger
		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.checkPostOwnerShip("admin", app.deletePostHandler))
				r.Patch("/", app.checkPostOwnerShip("moderator", app.updatePostHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {

			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)

			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)

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

	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- server.Shutdown(ctx)
	}()

	app.logger.Infof("✅ Starting server on http://localhost%s", server.Addr)

	return http.ListenAndServe(server.Addr, server.Handler)
}
