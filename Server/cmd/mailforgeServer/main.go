package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	slog.Info("SERVER IS RUNNING ON ", "PORT NUMBER ", cfg.Port)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		err := server.ListenAndServe()

		if err != nil {
			log.Fatalf("Unable to start the server :: %s", err.Error())
		}
	}()

	//Blocking the thread

	<-done

	slog.Info("Shutting down the server ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Unable to shutdown the server :: ", slog.String("error", err.Error()))
	}

	slog.Info("Server shut down gracefull ")
}
