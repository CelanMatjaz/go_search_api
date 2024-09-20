package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	AccessTokenName  = "access_token"
	RefreshTokenName = "refresh_token"
)

const (
	AccessTokenDuration  = time.Hour * 4
	RefreshTokenDuration = time.Hour * 24 * 30
)

type JwtData struct {
	secret []byte
	issuer []byte
}

var JwtAuth JwtData

func InitJwt() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not provided as an env variable")
	}
	JwtAuth.secret = []byte(secret)

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		log.Fatal("JWT_ISSUER not provided as an env variable")
	}
	JwtAuth.issuer = []byte(issuer)
}

func CreateAccessToken(accountId int) (string, error) {
	currentTime := time.Now()

	claims := customAccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    string(JwtAuth.issuer),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(AccessTokenDuration)),
		},
		AccountId: accountId,
	}

	return createToken[*customAccessTokenClaims](&claims)
}

func CreateRefreshToken(accountId int, version int) (string, error) {
	currentTime := time.Now()

	claims := customRefreshTokenClaims{
		customAccessTokenClaims: customAccessTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    string(JwtAuth.issuer),
				IssuedAt:  jwt.NewNumericDate(currentTime),
				ExpiresAt: jwt.NewNumericDate(currentTime.Add(RefreshTokenDuration)),
			},
			AccountId: accountId,
		},
		RefreshTokenVersion: version,
	}

	return createToken[*customRefreshTokenClaims](&claims)
}

func createToken[C jwt.Claims](claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(JwtAuth.secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

type customAccessTokenClaims struct {
	jwt.RegisteredClaims
	AccountId int `json:"account_id"`
}

func (c *customAccessTokenClaims) Valid() error {
	return jwt.RegisteredClaims.Valid(c.RegisteredClaims)
}

func VerifyAccessToken(tokenString string) (int, bool) {
	if tokenString == "" {
		return 0, false
	}

	token, err := jwt.ParseWithClaims(tokenString, &customAccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtAuth.secret, nil
	})
	if err != nil {
		return 0, false
	}

	claims := token.Claims.(*customAccessTokenClaims)
	return claims.AccountId, true
}

type customRefreshTokenClaims struct {
	customAccessTokenClaims
	RefreshTokenVersion int `json:"refresh_token_version"`
}

func (c *customRefreshTokenClaims) Valid() error {
	return jwt.RegisteredClaims.Valid(c.RegisteredClaims)
}

func VerifyRefreshToken(tokenString string) (int, int, bool) {
	if tokenString == "" {
		return 0, 0, false
	}

	token, err := jwt.ParseWithClaims(tokenString, &customRefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtAuth.secret, nil
	})
	if err != nil {
		return 0, 0, false
	}

	claims := token.Claims.(*customRefreshTokenClaims)
	return claims.AccountId, claims.RefreshTokenVersion, true
}
