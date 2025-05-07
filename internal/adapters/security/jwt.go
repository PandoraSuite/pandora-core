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

func (p *JWTProvider) Generate(
	ctx context.Context, subject string,
) (*dto.TokenResponse, errors.Error) {
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
		return nil, errors.NewInternal("failed to sign access token", err)
	}

	return &dto.TokenResponse{
		TokenType:   "Bearer",
		ExpiresIn:   expTime,
		AccessToken: tokenStr,
	}, nil
}

func (p *JWTProvider) Validate(
	ctx context.Context, token *dto.TokenValidation,
) (string, errors.Error) {
	if token.TokenType != "Bearer" {
		return "", errors.NewUnauthorized("invalid access token type, expected 'Bearer'")
	}

	t, err := jwt.Parse(
		token.AccessToken,
		func(token *jwt.Token) (any, error) {
			return p.secret, nil
		},
	)

	if err != nil || !t.Valid {
		return "", errors.NewUnauthorized("invalid access token")
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}
	return "", errors.NewUnauthorized("invalid access token claims")
}

func NewJWTProvider(secret []byte) *JWTProvider {
	return &JWTProvider{secret: secret}
}
