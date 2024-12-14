package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Id: "test-user",
	}
}
