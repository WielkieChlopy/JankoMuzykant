package main

import (
	v1 "backend/api/v1"
	"backend/db"
	"backend/router"
	"backend/store"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db := db.NewDb()
	us := store.NewUserStore(db)

	h, err := v1.NewHandler(us)
	if err != nil {
		log.Fatal(err)
	}

	e := router.New()

	v1 := e.Group("/api/v1")
	h.Register(v1)

	e.Logger.Fatal(e.Start(":8080"))
}
