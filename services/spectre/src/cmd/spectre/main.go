package main

import (
	"spectre/internal/storage/sqlite"
	"spectre/pkg/logger"
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
	st, err := sqlite.New(dbPath)
	if err != nil {
		log.Panic(err.Error())
	}
	_ = st

	// init server
	// run
}
