package sqlstore

import (
	"database/sql"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/store"
)

type UserRepository struct {
	store *Store
}

func (repository *UserRepository) Create(u *model.User) error {
	return repository.store.db.QueryRow(
		"INSERT INTO users (id, ipAddress) VALUES ($1, $2) RETURNING id",
		u.Id,
		u.IpAddress,
	).Scan(&u.Id)
}

func (repository *UserRepository) Find(id int) (*model.User, error) {
	user := &model.User{}
	if err := repository.store.db.QueryRow(
		"SELECT id, ipAddress FROM users WHERE id = $1",
		id,
	).Scan(
		&user.Id,
		&user.IpAddress,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}
