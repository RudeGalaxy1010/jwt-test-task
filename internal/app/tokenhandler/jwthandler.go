package tokenhandler

import (
	"errors"
	"time"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtHandler struct {
	key []byte
}

type UserClaims struct {
	Id           string
	IpAddress    string
	RefreshToken string
}

func JwtNew(key []byte) *JwtHandler {
	return &JwtHandler{
		key: key,
	}
}

func (handler *JwtHandler) NewAccessToken(user *model.User, ipAddress, refreshToken string) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.Id,
		"ip":  ipAddress,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"jti": refreshToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString(handler.key)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (handler *JwtHandler) NewRefreshToken() string {
	return uuid.New().String()
}

func (handler *JwtHandler) Decode(jwtToken string) (*UserClaims, error) {
	userClaims := &UserClaims{}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("failed to parse token")
		}

		return handler.key, nil
	})

	if err != nil {
		return userClaims, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userClaims.Id = claims["sub"].(string)
		userClaims.IpAddress = claims["ip"].(string)
		userClaims.RefreshToken = claims["jti"].(string)
	}

	return userClaims, nil
}

func (handler *JwtHandler) ValidateRefresh(jwtToken, refreshToken string) error {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("failed to parse token")
		}

		return handler.key, nil
	})

	if err != nil {
		return err
	}

	refreshId := ""

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		refreshId = claims["jti"].(string)
	}

	if refreshId != refreshToken {
		return errors.New("invalid token pair")
	}

	return nil
}
