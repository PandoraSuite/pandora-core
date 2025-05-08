package security

import (
	"context"
	"encoding/json"
	"os"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`

	ForcePasswordReset bool `json:"force_password_reset"`
}

type CredentialsRepository interface {
	GetByUsername(ctx context.Context, username string) (*entities.Credentials, errors.Error)
	ChangePassword(ctx context.Context, credentials *entities.Credentials) errors.Error
}

type credentialsRepository struct {
	credentialsFile string

	credentials *credentials
}

func (r *credentialsRepository) GetByUsername(
	ctx context.Context, username string,
) (*entities.Credentials, errors.Error) {
	if username != r.credentials.Username {
		return nil, errors.NewNotFound("credentials not found")
	}

	return &entities.Credentials{
		Username:           r.credentials.Username,
		HashedPassword:     r.credentials.Password,
		ForcePasswordReset: r.credentials.ForcePasswordReset,
	}, nil
}

func (r *credentialsRepository) ChangePassword(
	ctx context.Context, credentials *entities.Credentials,
) errors.Error {
	if credentials.Username != r.credentials.Username {
		return errors.NewNotFound("credentials not found")
	}

	oldPassword := r.credentials.Password
	r.credentials.Password = credentials.HashedPassword

	if r.credentials.ForcePasswordReset {
		r.credentials.ForcePasswordReset = false
	}

	if err := r.saveCredentials(); err != nil {
		r.credentials.Password = oldPassword
		return errors.NewInternal("failed to save credentials", err)
	}
	return nil
}

func (r *credentialsRepository) saveCredentials() error {
	data, err := json.Marshal(r.credentials)
	if err != nil {
		return err
	}

	err = os.WriteFile(r.credentialsFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func NewCredentialsRepository(credentialsFile string) (CredentialsRepository, error) {
	file, err := os.Open(credentialsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var credentials credentials
	if err := json.NewDecoder(file).Decode(&credentials); err != nil {
		return nil, err
	}

	return &credentialsRepository{
		credentials:     &credentials,
		credentialsFile: credentialsFile,
	}, nil
}
