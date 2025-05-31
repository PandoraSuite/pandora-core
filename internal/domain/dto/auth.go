package dto

import "time"

// ... Requests ...

type Authenticate struct {
	Username string `name:"username" validate:"required"`
	Password string `name:"password" validate:"required"`
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
