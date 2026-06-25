package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("chave_secreta_barberflow_super_segura")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateTokens(email string) (string, string, error) {

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(signedToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}
