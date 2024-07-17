package service_auth

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"login"`
	jwt.StandardClaims
}

func CookieJWTGet(log *slog.Logger, cookie *http.Cookie, jwtKey []byte) (*Claims, int, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			log.Error("Token signature error")
			return nil, http.StatusUnauthorized, err
		}
		log.Error("Token parsing error")
		return nil, http.StatusInternalServerError, err
	}
	if !tkn.Valid {
		log.Error("Token validation error")
		return nil, http.StatusUnauthorized, err
	}
	log.Debug("get token from cookie success")
	return claims, http.StatusOK, nil
}

func CookieJWTCreate(log *slog.Logger,
	claims *Claims,
	jwtExpiresAt *time.Time,
	jwtKey []byte) (*http.Cookie, int, error) {

	if jwtExpiresAt != nil {
		claims.StandardClaims.ExpiresAt = jwtExpiresAt.Unix()
	}
	log.Debug("Claim JWT", slog.Any("claims", claims))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Error("Token signature error")
		return nil, http.StatusUnauthorized, err
	}
	log.Debug("JWT", slog.String("tokenString", tokenString))

	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	if jwtExpiresAt != nil {
		cookie.Expires = *jwtExpiresAt
	} else {
		cookie.MaxAge = 0
	}
	log.Debug("Create token-cookie success")
	return &cookie, http.StatusOK, nil
}
