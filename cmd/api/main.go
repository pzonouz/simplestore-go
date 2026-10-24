package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	db := DBInit()
	defer db.Close()
	mux := chi.NewMux()
	mux.Use(middleware.Logger)
	Mount(db, mux)

	log.Println("Started Server on port 8000")
	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err.Error())
	}
}
