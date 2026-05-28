package jwt

import (
	"fmt"
	"seo-backend/internal/helper"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Generator struct {
	secret string
	expiry time.Duration
}

func NewGenerator(secret string, expiry string) (*Generator, error) {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		duration = 24 * time.Hour // default 24 hours
	}

	return &Generator{
		secret: secret,
		expiry: duration,
	}, nil
}

// GenerateToken generates a JWT token
func (g *Generator) GenerateToken(userID, teamID, email, name, role string, perms []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"name":    name,
		"role":    role,
		"exp":     helper.ParseWIBTime(time.Now().Format(time.RFC3339)).Add(g.expiry).Unix(),
		"iat":     helper.ParseWIBTime(time.Now().Format(time.RFC3339)).Unix(),
	}

	if teamID != "" {
		claims["team_id"] = teamID
	}

	if len(perms) > 0 {
		claims["permissions"] = perms
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secret))
}

// ValidateToken validates a JWT token
func (g *Generator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(g.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
