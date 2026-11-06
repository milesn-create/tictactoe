package appservice

import (
	"fmt"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainservice"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtProvider struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJwtProvider(accessSecret, refreshSecret []byte, accessTTL, refreshTTL time.Duration) domainservice.JwtProvider {
	return &jwtProvider{accessSecret: accessSecret, refreshSecret: refreshSecret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}
func (jp *jwtProvider) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{Subject: userID.String(), ExpiresAt: jwt.NewNumericDate(time.Now().Add(jp.accessTTL)), IssuedAt: jwt.NewNumericDate(time.Now())}
	tokenNotSigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenNotSigned.SignedString(jp.accessSecret)

}
func (jp *jwtProvider) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{Subject: userID.String(), ExpiresAt: jwt.NewNumericDate(time.Now().Add(jp.refreshTTL)), IssuedAt: jwt.NewNumericDate(time.Now())}
	tokenNotSigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenNotSigned.SignedString(jp.refreshSecret)
}

func (jp *jwtProvider) ValidateAccessToken(token string) error {
	claims := jwt.RegisteredClaims{}

	parsed, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return jp.accessSecret, nil
	})
	if err != nil {
		return fmt.Errorf("parse access token string: %w", err)
	}
	if !parsed.Valid {
		return domainerror.ErrInvalidToken
	}
	return nil
}
func (jp *jwtProvider) ValidateRefreshToken(token string) error {
	claims := jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return jp.refreshSecret, nil
	})
	if err != nil {
		return fmt.Errorf("parse refresh token string: %w", err)
	}
	if !parsed.Valid {
		return domainerror.ErrInvalidToken
	}
	return nil
}
func (jp *jwtProvider) GetUUIDFromToken(token string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return jp.accessSecret, nil
	})
	if err != nil {
		parsed, err = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
			return jp.refreshSecret, nil
		})
		if err != nil {
			return uuid.Nil, fmt.Errorf("get uuid from token: %w", err)
		}
	}
	if !parsed.Valid {
		return uuid.Nil, domainerror.ErrInvalidToken
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}
	return userID, nil

}
