package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type Credentials struct {
	Username string `name:"username" validate:"required"`
	Password string `name:"password" validate:"required"`
}

type Authenticate struct {
	*Credentials
}

type Reauthenticate struct {
	*Credentials
	Action enums.SensitiveAction `name:"action" validate:"required,enums=REVEAL_API_KEY"`
}

type ChangePassword struct {
	Username        string `name:"username" validate:"required"`
	NewPassword     string `name:"new_password" validate:"required,min=12,eqfield=ConfirmPassword"`
	ConfirmPassword string `name:"confirm_password" validate:"required"`
}

// ... Responses ...

type TokenResponse struct {
	AccessToken string    `name:"access_token"`
	ExpiresIn   time.Time `name:"expires_in"`
}

type AuthenticateResponse struct {
	*TokenResponse
	ForcePasswordReset bool `name:"force_password_reset"`
}

type ReauthenticateResponse struct {
	*TokenResponse
}
