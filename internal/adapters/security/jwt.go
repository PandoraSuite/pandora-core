package security

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports"
)

type jwtProvider struct {
	secret []byte
}

func (p *jwtProvider) GenerateAccessToken(
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

	return p.signToken(claims, expTime)
}

func (p *jwtProvider) GenerateScopedJWT(
	ctx context.Context, subject, scope string,
) (*dto.TokenResponse, errors.Error) {
	now := time.Now()
	expTime := now.Add(1 * time.Minute)

	claims := jwt.MapClaims{
		"iss":   "pandora-core",
		"sub":   subject,
		"scope": scope,
		"exp":   expTime.Unix(),
		"nbf":   now.Unix(),
		"iat":   now.Unix(),
	}

	return p.signToken(claims, expTime)
}

func (p *jwtProvider) signToken(
	claims jwt.Claims, exp time.Time,
) (*dto.TokenResponse, errors.Error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(p.secret)
	if err != nil {
		return nil, errors.NewInternal("failed to sign token", err)
	}

	return &dto.TokenResponse{
		ExpiresIn:   exp,
		AccessToken: tokenStr,
	}, nil
}

func (p *jwtProvider) ValidateAccessToken(
	ctx context.Context, token string,
) (string, errors.Error) {
	t, err := p.validate(token)
	if err != nil {
		return "", err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}
	return "", errors.NewUnauthorized("Invalid access token claims", nil)
}

func (p *jwtProvider) ValidateScopedToken(
	ctx context.Context, token, expectedScope string,
) errors.Error {
	t, err := p.validate(token)
	if err != nil {
		return err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return errors.NewUnauthorized("Invalid token claims", nil)
	}

	scope, ok := claims["scope"].(string)
	if !ok || scope != expectedScope {
		return errors.NewForbidden("Token does not grant required scope", nil)
	}

	return nil
}

func (p *jwtProvider) validate(token string) (*jwt.Token, errors.Error) {
	t, err := jwt.Parse(
		token,
		func(token *jwt.Token) (any, error) { return p.secret, nil },
	)

	if err != nil || !t.Valid {
		return nil, errors.NewUnauthorized("Invalid access token", err)
	}

	return t, nil
}

func NewJWTProvider(secret []byte) ports.TokenProvider {
	return &jwtProvider{secret: secret}
}
