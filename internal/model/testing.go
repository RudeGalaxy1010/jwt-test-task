package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Id:        "test-user",
		IpAddress: "127.0.0.1",
	}
}
