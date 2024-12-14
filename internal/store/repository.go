package store

import "github.com/RudeGalaxy1010/jwt-test-task/internal/model"

type UserRepository interface {
	Create(*model.User) error
	Find(string) (*model.User, error)
	UpdateRefreshToken(*model.User) error
}
