package security

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret []byte
}

func (p *JWTProvider) GenerateToken(ctx context.Context, subject string) (*entities.Token, error) {
	now := time.Now()
	expTime := now.Add(time.Hour)

	claims := jwt.MapClaims{
		"iss": "pandora-core",
		"sub": subject,
		"exp": expTime.Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(p.secret)
	if err != nil {
		return nil, domainErr.ErrTokenSigningFailed
	}

	return &entities.Token{
		Token:     tokenStr,
		TokenType: "Bearer",
		ExpiresAt: expTime,
	}, nil
}

func (p *JWTProvider) ValidateToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return p.secret, nil
	})

	if err != nil || !token.Valid {
		return "", domainErr.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}
	return "", domainErr.ErrInvalidTokenData
}

func NewJWTProvider(secret []byte) *JWTProvider {
	return &JWTProvider{secret: secret}
}
