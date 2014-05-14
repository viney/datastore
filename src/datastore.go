package datastore

import (
	"database/sql"
	"io"
	"os"

	_ "goracle"
)

func Open() (*sql.DB, error) {
	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
	// open
	db, err := sql.Open("goracle", "viney/admin@test")
	if err != nil {
		return nil, err
	}

	// suc
	return db, nil
}

type Closer interface {
	io.Closer
}

func Close(db Closer) {
	if db != nil {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}
}

func Rollback(tx *sql.Tx) {
	err := tx.Rollback()

	// roll back error is a serious error
	if err != nil {
		panic(err)
	}
}
