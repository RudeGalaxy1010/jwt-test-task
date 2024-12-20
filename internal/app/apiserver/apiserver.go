package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/app/tokenhandler"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/store/sqlstore"
)

func Start(config *Config) error {
	db, err := NewDB(config.Database_url)

	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	jwtHandler := tokenhandler.JwtNew([]byte(config.Secret_key))
	server := NewServer(store, *jwtHandler)

	return http.ListenAndServe(config.Address, server)
}

func NewDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
