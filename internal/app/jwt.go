package jwt

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokensPair struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func NewAccessToken(user *model.User, id, key string) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.Id,
		"ip":  user.IpAddress,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"jti": id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func NewRefreshToken(user *model.User) (string, error) {
	token := uuid.New().String()
	base64Token := base64.StdEncoding.EncodeToString([]byte(token))
	return base64Token, nil
}

func ValidateRefresh(access, refresh, key string) (string, error) {
	token, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("failed to parse token")
		}
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	ip := ""
	refreshId := ""

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ip = claims["ip"].(string)
		refreshId = claims["jti"].(string)
	}

	if refreshId != refresh {
		return "", errors.New("invalid token pair")
	}

	return ip, nil
}
