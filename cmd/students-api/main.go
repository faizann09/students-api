package main

import (
	"context"
	//"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faizann09/students-api/config"
    "github.com/faizann09/students-api/internal/http/handlers/students"
	
)

func main() {
	cfg := config.MustLoad()

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

		slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr))


	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start the server ")
		}
	}()

	<-done

	slog.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil{
				slog.Error("Failed to Shutdown Server", slog.String("error",err.Error()))

	}
}
