package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

// ... Requests ...

type Authenticate struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (a *Authenticate) ToDomain() *dto.Authenticate {
	return &dto.Authenticate{
		Username: a.Username,
		Password: a.Password,
	}
}

type ChangePassword struct {
	Username        string `json:"username" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func (c *ChangePassword) ToDomain() *dto.ChangePassword {
	return &dto.ChangePassword{
		Username:        c.Username,
		NewPassword:     c.NewPassword,
		ConfirmPassword: c.ConfirmPassword,
	}
}

type TokenValidation struct {
	TokenType   string
	AccessToken string
}

func (t *TokenValidation) ToDomain() *dto.TokenValidation {
	return &dto.TokenValidation{
		TokenType:   t.TokenType,
		AccessToken: t.AccessToken,
	}
}

// ... Responses ...

type AuthenticateResponse struct {
	TokenType          string    `json:"token_type"`
	ExpiresIn          time.Time `json:"expires_in" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	AccessToken        string    `json:"access_token"`
	ForcePasswordReset bool      `json:"force_password_reset"`
}

func AuthenticateResponseFromDomain(auth *dto.AuthenticateResponse) *AuthenticateResponse {
	return &AuthenticateResponse{
		TokenType:          auth.TokenType,
		ExpiresIn:          auth.ExpiresIn,
		AccessToken:        auth.AccessToken,
		ForcePasswordReset: auth.ForcePasswordReset,
	}
}
