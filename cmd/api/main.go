package main

import (
	"log"

	"github.com/MohammadBohluli/social-app-go/internal/db"
	"github.com/MohammadBohluli/social-app-go/internal/store"
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
		addr:   ":8000",
		apiUrl: "localhost:8000",
		db: dbConfig{
			addr:         "postgres://myusername:mypassword1234@localhost/social?sslmode=disable",
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
	}

	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("✅ Database is Connected")

	store := store.NewPostgresStorage(db)
	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.RegisterRoutes()

	log.Fatal(app.start(mux))
}
