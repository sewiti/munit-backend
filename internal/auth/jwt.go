package auth

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

var secretKey ed25519.PrivateKey

var ErrExpiredToken = errors.New("expired token")

func MakeJWT(subject string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    hostname,
		Subject:   subject,
	})
	return token.SignedString(secretKey)
}

func VerifyJWT(token string) (subject string, err error) {
	var claims jwt.RegisteredClaims
	_, err = jwt.ParseWithClaims(token, &claims, jwtKeyFunc)
	if err != nil {
		return "", err
	}

	now := time.Now()
	if claims.ExpiresAt != nil {
		if claims.ExpiresAt.Time.Before(now) {
			return "", ErrExpiredToken
		}
		return claims.Subject, nil
	}
	if claims.IssuedAt != nil {
		if claims.IssuedAt.Add(tokenExpiry).Before(now) {
			return "", ErrExpiredToken
		}
		return claims.Subject, nil
	}
	return "", errors.New("token has no expiration")
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return secretKey.Public(), nil
}
