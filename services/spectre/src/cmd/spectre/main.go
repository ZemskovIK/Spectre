package main

import (
	"net/http"
	"os"
	"os/signal"
	"spectre/api"
	"spectre/internal/storage/sqlite"
	"spectre/pkg/logger"
	"syscall"
)

const (
	logLevel = "DEBUG"
	dbPath   = "../db/spec.db"
)

func main() {
	// init logger
	log := logger.New(logLevel)
	_ = log

	// init config (?)

	// init storage
	st, err := sqlite.New(dbPath, log)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Info("init storage")

	// init router
	r := api.NewRouter(st, log)
	log.Info("init router")

	// run
	log.Info("running server")
	go func() {
		if err := http.ListenAndServe(":5000", r); err != nil {
			log.Fatalf("cannot run server on :5000 : %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Infof("shutdown :)")
}
