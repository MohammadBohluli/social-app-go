package main

import (
	"time"

	"github.com/MohammadBohluli/social-app-go/internal/auth"
	"github.com/MohammadBohluli/social-app-go/internal/db"
	"github.com/MohammadBohluli/social-app-go/internal/mailer"
	"github.com/MohammadBohluli/social-app-go/internal/store"
	"github.com/MohammadBohluli/social-app-go/internal/store/cache"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Gopher Social API
//	@description	API for gopher Social
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	// http://localhost:8000/v1/swagger/index.html
	cfg := config{
		addr:        ":8000",
		apiUrl:      "localhost:8000",
		frontEndURL: "http://localhost:4000",
		evn:         "development",
		redisCfg: redisConfig{
			host:     "localhost",
			port:     6379,
			password: "",
			db:       0,
			enabled:  false,
		},
		auth: authConfig{
			basic: basicConfig{
				username: "admin",
				password: "admin",
			},
			token: tokenConfig{
				secret: "my_token",
				exp:    time.Hour * 24 * 3, //3day
			},
		},
		db: dbConfig{
			addr:         "postgres://myusername:mypassword1234@localhost/social?sslmode=disable",
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
		mail: mailConfig{
			exp:       time.Hour * 23 * 3,
			fromEmail: "GopherSocial",
			sendGrid: sendGridConfig{
				apiKey: "API_KEY",
			},
		}, // 3 day

	}

	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("✅ Database is Connected")

	rdb := cache.New(cfg.redisCfg.host, cfg.redisCfg.port, cfg.redisCfg.password, cfg.redisCfg.db)

	store := store.NewPostgresStorage(db)

	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, "gopherSocial", "gopherSocial")
	app := application{
		config:        cfg,
		cacheStorage:  cache.NewRedisStorage(rdb),
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := app.RegisterRoutes()

	logger.Fatal(app.start(mux))
}
