package teststore_test

import (
	"testing"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	teststore "github.com/RudeGalaxy1010/jwt-test-task/internal/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	teststore := teststore.New()
	user := model.TestUser(t)
	assert.NoError(t, teststore.User().Create(user))
	assert.NotNil(t, user)
}

func TestUserRepository_Find(t *testing.T) {
	teststore := teststore.New()
	user := model.TestUser(t)
	teststore.User().Create(user)
	u, err := teststore.User().Find(user.Id)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
