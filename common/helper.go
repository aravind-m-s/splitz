package common

import (
	"fmt"
	"splitz/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTStruct struct {
	cnf *config.EnvModel
}

func NewHelper(cnf *config.EnvModel) *JWTStruct {
	return &JWTStruct{cnf: cnf}
}

func (c *JWTStruct) GenerateJWT(id uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)

	claims := &Claims{
		UserID: fmt.Sprintf("%v", id),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(c.cnf.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *JWTStruct) VerifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(c.cnf.JWTSecret), nil
	})

	return token, err
}

func (c *JWTStruct) GetFromToken(tokenString string, key string) (string, error) {
	token, err := c.VerifyJWT(tokenString)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims[key].(string)
		if !ok {
			return "", fmt.Errorf("%s claim is not present in the token", key)
		}
		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}
