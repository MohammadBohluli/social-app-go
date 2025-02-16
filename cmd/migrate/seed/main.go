package main

import (
	"log"

	"github.com/MohammadBohluli/social-app-go/internal/db"
	"github.com/MohammadBohluli/social-app-go/internal/store"
)

func main() {
	addr := "postgres://myusername:mypassword1234@localhost/social?sslmode=disable"
	conn, err := db.NewDB(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewPostgresStorage(conn)
	db.Seed(store)
}
