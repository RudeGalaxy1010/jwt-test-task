package teststore

import (
	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/store"
)

type Store struct {
	UserRepository *UserRepository
}

func New() *Store {
	return &Store{}
}

func (store *Store) User() store.UserRepository {
	if store.UserRepository == nil {
		store.UserRepository = &UserRepository{
			store: store,
			users: make(map[string]*model.User),
		}
	}

	return store.UserRepository
}
