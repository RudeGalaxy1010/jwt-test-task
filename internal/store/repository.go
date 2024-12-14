package store

import "github.com/RudeGalaxy1010/jwt-test-task/internal/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
}
