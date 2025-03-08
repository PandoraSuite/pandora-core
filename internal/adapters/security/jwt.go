package security

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret []byte
}

func (p *JWTProvider) GenerateToken(ctx context.Context, subject string) (*dto.AuthenticateResponse, error) {
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

	return &dto.AuthenticateResponse{
		Token:     tokenStr,
		TokenType: "Bearer",
		ExpiresAt: expTime,
	}, nil
}

func (p *JWTProvider) ValidateToken(ctx context.Context, token *dto.TokenRequest) (string, error) {
	if token.Type != "Bearer" {
		return "", domainErr.ErrInvalidTokenType
	}

	t, err := jwt.Parse(token.Key, func(token *jwt.Token) (any, error) {
		return p.secret, nil
	})

	if err != nil || !t.Valid {
		return "", domainErr.ErrInvalidToken
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}
	return "", domainErr.ErrInvalidTokenData
}

func NewJWTProvider(secret []byte) *JWTProvider {
	return &JWTProvider{secret: secret}
}
