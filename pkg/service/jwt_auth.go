package service

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/golang-jwt/jwt/v4"
)

var JwtClient JwtAuth

type JwtAuth struct {
	secret []byte
}

func (a *JwtAuth) InitJwtAuth() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not provided as an env variable")
	}
	a.secret = []byte(secret)

	return nil
}

func (a *JwtAuth) CreateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		UserIdKey: userId,
	})

	signed, err := token.SignedString(a.secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

type customClaims struct {
	UserId int `json:"user_id"`
}

func (c *customClaims) Valid() error {
	if c.UserId == 0 {
		return types.UserIdNotProvidedErr
	}

	return nil
}

func (a *JwtAuth) VerifyToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("authorization")
	if authHeader == "" {
		return 0, types.MissingRequiredHeaderErr
	}

	isValid := strings.HasPrefix(authHeader, "Bearer ")
	if !isValid {
		return 0, types.WronglyFormattedAuthHeaderErr
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})

	claims := token.Claims.(*customClaims)

	return claims.UserId, err

}
