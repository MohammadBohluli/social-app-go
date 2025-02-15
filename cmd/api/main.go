package main

import (
	"log"

	"github.com/MohammadBohluli/social-app-go/internal/db"
	"github.com/MohammadBohluli/social-app-go/internal/store"
)

func main() {

	cfg := config{
		addr: ":8000",
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
