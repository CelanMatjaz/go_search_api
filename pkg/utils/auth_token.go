package utils

import (
	"log"
	"os"
	"time"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/golang-jwt/jwt/v4"
)

var JwtClient JwtAuth

var UserIdKey = "user_id"

var (
	AccessTokenName  = "access_token"
	RefreshTokenName = "refresh_token"
)

const (
	AccessTokenDuration  = time.Second * 10
	RefreshTokenDuration = time.Hour * 24 * 30
)

type JwtAuth struct {
	secret []byte
	issuer []byte
}

func (a *JwtAuth) InitJwtAuth() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not provided as an env variable")
	}
	a.secret = []byte(secret)

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		log.Fatal("JWT_ISSUER not provided as an env variable")
	}
	a.issuer = []byte(issuer)

	return nil
}

func (a *JwtAuth) CreateAccessToken(userId int) (string, error) {
	currentTime := time.Now()

	claims := customAccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    string(a.issuer),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(AccessTokenDuration)),
		},
		UserId: userId,
	}

	return createToken[*customAccessTokenClaims](a, &claims)
}

func (a *JwtAuth) CreateRefreshToken(userId int, version int) (string, error) {
	currentTime := time.Now()

	claims := customRefreshTokenClaims{
		customAccessTokenClaims: customAccessTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    string(a.issuer),
				IssuedAt:  jwt.NewNumericDate(currentTime),
				ExpiresAt: jwt.NewNumericDate(currentTime.Add(RefreshTokenDuration)),
			},
			UserId: userId,
		},
		RefreshTokenVersion: version,
	}

	return createToken[*customRefreshTokenClaims](a, &claims)
}

func createToken[C jwt.Claims](a *JwtAuth, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(a.secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

type customAccessTokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func (c *customAccessTokenClaims) Valid() error {
	return jwt.RegisteredClaims.Valid(c.RegisteredClaims)
}

func (a *JwtAuth) VerifyAccessToken(tokenString string) (int, error) {
	if tokenString == "" {
		return 0, types.InvalidTokenErr
	}

	token, err := jwt.ParseWithClaims(tokenString, &customAccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})
	if err != nil {
		return 0, types.UnauthenticatedErr
	}

	claims := token.Claims.(*customAccessTokenClaims)
	return claims.UserId, err
}

type customRefreshTokenClaims struct {
	customAccessTokenClaims
	RefreshTokenVersion int `json:"refresh_token_version"`
}

func (c *customRefreshTokenClaims) Valid() error {
	return jwt.RegisteredClaims.Valid(c.RegisteredClaims)
}

func (a *JwtAuth) VerifyRefreshToken(tokenString string) (int, int, error) {
	if tokenString == "" {
		return 0, 0, types.InvalidTokenErr
	}

	token, err := jwt.ParseWithClaims(tokenString, &customRefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})
	if err != nil {
		return 0, 0, types.UnauthenticatedErr
	}

	claims := token.Claims.(*customRefreshTokenClaims)
	return claims.UserId, claims.RefreshTokenVersion, err
}
