package apiserver

import "errors"

var (
	ErrInsufficientRequest = errors.New("request not fully qualified")
	ErrUserCreationFailed  = errors.New("failed to create user")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserUpdateFailed    = errors.New("user update failed")
	ErrTokenCreationFailed = errors.New("failed to create token")
	ErrTokenParseFailed    = errors.New("failed to parse token")
	ErrTokenInvalid        = errors.New("token is invalid")
)
