package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type Authenticate struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (a *Authenticate) ToDomain() *dto.Authenticate {
	return &dto.Authenticate{
		Credentials: &dto.Credentials{
			Username: a.Username,
			Password: a.Password,
		},
	}
}

type ChangePassword struct {
	Username        string `json:"-" swaggerignore:"true"`
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

type Reauthenticate struct {
	Action   string `json:"action" binding:"required" enums:"REVEAL_API_KEY"`
	Username string `json:"-" swaggerignore:"true"`
	Password string `json:"password" binding:"required"`
}

func (r *Reauthenticate) ToDomain() *dto.Reauthenticate {
	return &dto.Reauthenticate{
		Action: enums.SensitiveAction(r.Action),
		Credentials: &dto.Credentials{
			Username: r.Username,
			Password: r.Password,
		},
	}
}

// ... Responses ...

type TokenReponse struct {
	TokenType   string    `json:"token_type"`
	ExpiresIn   time.Time `json:"expires_in" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	AccessToken string    `json:"access_token"`
}

type AuthenticateResponse struct {
	*TokenReponse
	ForcePasswordReset bool `json:"force_password_reset"`
}

func AuthenticateResponseFromDomain(auth *dto.AuthenticateResponse) *AuthenticateResponse {
	return &AuthenticateResponse{
		TokenReponse: &TokenReponse{
			TokenType:   "Bearer",
			ExpiresIn:   auth.ExpiresIn,
			AccessToken: auth.AccessToken,
		},
		ForcePasswordReset: auth.ForcePasswordReset,
	}
}

type ReauthenticateResponse struct {
	*TokenReponse
}

func ReauthenticateResponseFromDomain(auth *dto.ReauthenticateResponse) *ReauthenticateResponse {
	return &ReauthenticateResponse{
		TokenReponse: &TokenReponse{
			TokenType:   "Bearer",
			ExpiresIn:   auth.ExpiresIn,
			AccessToken: auth.AccessToken,
		},
	}
}
