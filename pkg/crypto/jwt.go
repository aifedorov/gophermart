package crypto

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/aifedorov/gophermart/internal/domain/user"
)

type JWTManager interface {
	GenerateToken(userID user.ID) (string, error)
	ValidateToken(tokenString string) (user.ID, error)
}

type Claims struct {
	UserID user.ID `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtManager struct {
	secretKey []byte
	tokenTTL  time.Duration
}

func NewJWTManager(secretKey string, tokenTTL time.Duration) JWTManager {
	return &jwtManager{
		secretKey: []byte(secretKey),
		tokenTTL:  tokenTTL,
	}
}

func (j *jwtManager) GenerateToken(userID user.ID) (string, error) {
	return "", nil
}

func (j *jwtManager) ValidateToken(tokenString string) (user.ID, error) {
	return user.ID(""), nil
}
