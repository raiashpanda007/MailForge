package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second))
	router.Get("/api/MailForger", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Welcome to mail forge server side"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server running at the server :: ", ":8080")
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("Unable to start the server :: %s", err.Error())
	}

}
