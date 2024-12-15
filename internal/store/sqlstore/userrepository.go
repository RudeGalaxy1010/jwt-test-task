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
		"INSERT INTO users (id) VALUES ($1, $2) RETURNING id",
		u.Id,
	).Scan(&u.Id)
}

func (repository *UserRepository) Find(id string) (*model.User, error) {
	user := &model.User{}
	if err := repository.store.db.QueryRow(
		"SELECT id, refresh FROM users WHERE id = $1",
		id,
	).Scan(
		&user.Id,
		&user.Refresh,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

func (repository *UserRepository) UpdateRefreshToken(user *model.User) error {
	id := ""

	return repository.store.db.QueryRow(
		"UPDATE users SET refresh = $1 WHERE id = $2 RETURNING id",
		user.Refresh,
		user.Id,
	).Scan(&id)
}
