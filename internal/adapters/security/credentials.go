package security

import (
	"context"
	"encoding/json"
	"os"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`

	ForcePasswordReset bool `json:"force_password_reset"`
}

type CredentialsRepository struct {
	credentialsFile string

	credentials *credentials
}

func (r *CredentialsRepository) FindCredentials(
	ctx context.Context, username string,
) (*entities.Credentials, error) {
	if username != r.credentials.Username {
		return nil, domainErr.ErrCredentialsNotFound
	}

	return &entities.Credentials{
		Username:           r.credentials.Username,
		HashedPassword:     r.credentials.Password,
		ForcePasswordReset: r.credentials.ForcePasswordReset,
	}, nil
}

func (r *CredentialsRepository) ChangePassword(
	ctx context.Context, credentials *entities.Credentials,
) error {
	if credentials.Username != r.credentials.Username {
		return domainErr.ErrCredentialsNotFound
	}

	oldPassword := r.credentials.Password
	r.credentials.Password = credentials.HashedPassword

	if r.credentials.ForcePasswordReset {
		r.credentials.ForcePasswordReset = false
	}

	if err := r.saveCredentials(); err != nil {
		r.credentials.Password = oldPassword
		return err
	}
	return nil
}

func (r *CredentialsRepository) saveCredentials() error {
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

func NewCredentialsRepository(credentialsFile string) (*CredentialsRepository, error) {
	file, err := os.Open(credentialsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var credentials credentials
	if err := json.NewDecoder(file).Decode(&credentials); err != nil {
		return nil, err
	}

	return &CredentialsRepository{
		credentials:     &credentials,
		credentialsFile: credentialsFile,
	}, nil
}
