package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/logger"
)

type ContextKey string

const (
	UserIDKey  ContextKey = "user_id"
	CookieName            = "JWT"
)

const (
	tokenExp   = time.Hour * 1
	cookiesExp = time.Hour * 1
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type JWTMiddleware struct {
	secretKey string
}

func NewJWTMiddleware(secretKey string) *JWTMiddleware {
	return &JWTMiddleware{
		secretKey: secretKey,
	}
}

func (m *JWTMiddleware) CheckJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if errors.Is(err, http.ErrNoCookie) {
			logger.Log.Info("auth: no auth cookies", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userID, err := parseUserID(cookie.Value, m.secretKey)
		if err != nil {
			logger.Log.Info("auth: failed to get cookie", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetNewAuthCookies(userID, secretKey string, w http.ResponseWriter) {
	token, err := buildJWTString(userID, secretKey)
	if err != nil {
		logger.Log.Error("auth: failed to build JWT token", zap.Error(err))
		return
	}

	cookie := http.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  time.Now().Add(cookiesExp),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)
}

func GetUserID(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", errors.New("user id not found")
	}
	return userID, nil
}

func parseUserID(tokenString, secretKey string) (string, error) {
	if tokenString == "" {
		logger.Log.Error("auth: empty token")
		return "", errors.New("auth: token is empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("auth: unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		logger.Log.Info("auth: error parsing token", zap.Error(err))
		return "", errors.New("auth: invalid token")
	}

	if !token.Valid {
		logger.Log.Info("auth: invalid token")
		return "", errors.New("auth: invalid token")
	}
	return claims.UserID, nil
}

func buildJWTString(userID, secretKey string) (string, error) {
	logger.Log.Debug("auth: building JWT token with user_id", zap.String("user_id", userID))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gophermart",
			Audience:  []string{"gophermart-users"},
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.Log.Error("auth: failed to sign JWT token", zap.Error(err))
		return "", err
	}
	return tokenString, nil
}
