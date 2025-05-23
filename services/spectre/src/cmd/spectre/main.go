package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	server "spectre/internal/srv"
	"spectre/internal/storage/sqlite"
	"spectre/pkg/logger"
	"syscall"
	"time"
)

const (
	logLevel = "DEBUG"
	dbPath   = "../db/spec.db"
)

func main() {
	// init logger
	log := logger.New(logLevel)

	// init config (?)

	// init storage
	st, err := sqlite.NewStorage(dbPath, log)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Info("init storage")

	// init router
	r := server.NewRouter(st, log)
	r.Use(server.AuthMiddleware)
	r.Use(server.CORSMiddleware)
	r.Use(server.JSONRespMiddleware)
	log.Info("init router")

	// init server
	srv := &http.Server{
		Addr:    ":5000",
		Handler: r,
	}
	log.Info("init server")

	// run
	go func() {
		log.Info("running server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("cannot run server on :5000 : %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Info("server exited gracefully :)")
}
