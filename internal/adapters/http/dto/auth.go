package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type Authenticate struct {
	Username string `form:"username" validate:"required"`

	Password string `form:"password" validate:"required" format:"password" minLength:"12"`
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
	Username string `json:"-" swaggerignore:"true"`

	NewPassword string `json:"new_password" validate:"required" format:"password" minLength:"12"`

	ConfirmPassword string `json:"confirm_password" validate:"required" format:"password" minLength:"12"`
}

func (c *ChangePassword) ToDomain() *dto.ChangePassword {
	return &dto.ChangePassword{
		Username:        c.Username,
		NewPassword:     c.NewPassword,
		ConfirmPassword: c.ConfirmPassword,
	}
}

type Reauthenticate struct {
	Username string `json:"-" swaggerignore:"true"`

	Action string `json:"action" validate:"required" enums:"REVEAL_API_KEY"`

	Password string `json:"password" validate:"required" format:"password" minLength:"12"`
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
	TokenType string `json:"token_type"`

	ExpiresIn time.Time `json:"expires_in" format:"date-time" extensions:"x-timezone=utc"`

	AccessToken string `json:"access_token" format:"jwt"`
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
