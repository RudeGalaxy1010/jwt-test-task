package teststore

import (
	"fmt"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (repository *UserRepository) Create(u *model.User) error {
	u.Id = len(repository.users) + 1
	fmt.Println(u.Id)
	repository.users[u.Id] = u
	return nil
}

func (repository *UserRepository) Find(id int) (*model.User, error) {
	u, ok := repository.users[id]

	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}