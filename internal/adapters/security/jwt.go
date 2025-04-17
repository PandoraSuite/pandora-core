package security

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type JWTProvider struct {
	secret []byte
}

func (p *JWTProvider) GenerateToken(
	ctx context.Context, subject string,
) (*dto.TokenResponse, *errors.Error) {
	now := time.Now().UTC()
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
		return nil, errors.ErrTokenSigningFailed
	}

	return &dto.TokenResponse{
		Token:     tokenStr,
		TokenType: "Bearer",
		ExpiresIn: expTime,
	}, nil
}

func (p *JWTProvider) ValidateToken(
	ctx context.Context, token *dto.TokenRequest,
) (string, *errors.Error) {
	if token.Type != "Bearer" {
		return "", errors.ErrInvalidTokenType
	}

	t, err := jwt.Parse(token.Key, func(token *jwt.Token) (any, error) {
		return p.secret, nil
	})

	if err != nil || !t.Valid {
		return "", errors.ErrInvalidToken
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}
	return "", errors.ErrInvalidTokenData
}

func NewJWTProvider(secret []byte) *JWTProvider {
	return &JWTProvider{secret: secret}
}
