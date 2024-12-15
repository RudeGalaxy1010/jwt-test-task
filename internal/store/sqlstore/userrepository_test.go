package sqlstore_test

import (
	"testing"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	sqlstore "github.com/RudeGalaxy1010/jwt-test-task/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown("users")

	sqlstore := sqlstore.New(db)
	user := model.TestUser(t)
	assert.NoError(t, sqlstore.User().Create(user))
	assert.NotNil(t, user)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown("users")

	sqlstore := sqlstore.New(db)
	u := model.TestUser(t)

	sqlstore.User().Create(&model.User{
		Id: u.Id,
	})
	user, err := sqlstore.User().Find(u.Id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
