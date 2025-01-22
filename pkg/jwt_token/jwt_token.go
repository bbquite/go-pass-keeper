package jwttoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	ExpiryHour time.Duration
	Secret     string
}

func NewJWTTokenManager(exp time.Duration, secret string) *JWTManager {
	return &JWTManager{
		ExpiryHour: exp,
		Secret:     secret,
	}
}

type JWT struct {
	Token string `json:"token"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint32
}

func (tm *JWTManager) CreateAccessToken(userid uint32) (accessToken string, err error) {
	claims := &Claims{
		UserID: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.ExpiryHour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(tm.Secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func (tm *JWTManager) IsAuthorized(requestToken string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.Secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (tm *JWTManager) ExtractIDFromToken(requestToken string) (uint32, error) {
	token, err := jwt.ParseWithClaims(requestToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID
		return userID, nil
	}

	return 0, fmt.Errorf("invalid Token")
}
