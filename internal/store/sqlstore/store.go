package sqlstore

import (
	"database/sql"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db             *sql.DB
	UserRepository *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) User() store.UserRepository {
	if store.UserRepository != nil {
		store.UserRepository = &UserRepository{
			store: store,
		}
	}

	return store.UserRepository
}
