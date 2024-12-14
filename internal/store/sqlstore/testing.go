package sqlstore

import (
	"database/sql"
	"strings"
	"testing"
)

func TestDB(t *testing.T, databaseUrl string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("postgres", databaseUrl)

	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(s ...string) {
		if len(s) > 0 {
			db.Exec("TRUNCATE %s CASCADE;", strings.Join(s, ", "))
		}

		db.Close()
	}
}
