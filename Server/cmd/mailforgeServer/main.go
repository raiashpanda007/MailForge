package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raiashpanda007/MailForge/pkg/config"
)

func main() {
	cfg := config.MustLoad()
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
		Addr:    cfg.HTTPServer.Hostname + cfg.HTTPServer.Port,
		Handler: router,
	}

	fmt.Println("Server running at the server :: ", cfg.HTTPServer.Port)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("Unable to start the server :: %s", err.Error())
	}

}
