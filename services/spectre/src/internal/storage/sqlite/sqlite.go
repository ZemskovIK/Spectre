package sqlite

import (
	"database/sql"
	st "spectre/internal/storage"
	"spectre/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

const (
	GLOC_SQL = "src/internal/storage/sqlite/sqlite.go/"

	UNK_NAME = ""
)

type sqliteDB struct {
	db  *sql.DB
	log *logger.Logger
}

func NewStorage(dbPath string, log *logger.Logger) (st.Storage, error) {
	loc := GLOC_SQL + "NewStorage()"
	log.Debugf("%s: opening sqlite db at path: %s", loc, dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Errorf("%s: failed to open sqlite db: %v", loc, err)
		return nil, errCannotConnectSQLite(loc, err)
	}

	return &sqliteDB{
		db:  db,
		log: log,
	}, nil
}
