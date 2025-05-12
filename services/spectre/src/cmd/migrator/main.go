package main

import (
	"flag"
	"spectre/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migsPath = "file://../db/migrations"
	dbPath   = "sqlite3://../db/spec.db"

	logLevel = "debug"
)

func main() {
	log := logger.New(logLevel)

	action := flag.String("action", "up", "action: up or down or force")
	vers := flag.Int("version", -1, "version to migrate (force)")

	flag.Parse()
	log.Debugf("run migrations with action: %s, version: %d", *action, *vers)

	if !actionIsValid(*action) {
		log.Fatalf("invalid action: %s", *action)
	}

	m, err := migrate.New(migsPath, dbPath)
	if err != nil {
		log.Fatalf("cannot make migrations: %v", err)
	}

	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("cannot migrate up: %v", err)
		}
		log.Info("migrated up")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("cannot migrate down: %v", err)
		}
		log.Info("migrated down")
	case "force":
		v := *vers
		if v == -1 {
			log.Fatal("specify version to force")
		}
		if err := m.Force(v); err != nil {
			log.Fatalf("cannot force migrate: %v", err)
		}
		log.Infof("force migrated to version: %d", v)
	default:
		log.Fatalf("unknown action: %s", *action)
	}
}

func actionIsValid(action string) bool {
	switch action {
	case "up", "down":
		return true
	default:
		return false
	}
}
