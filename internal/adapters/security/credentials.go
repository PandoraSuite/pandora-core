package security

import (
	"context"
	"encoding/json"
	"os"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`

	ForcePasswordReset bool `json:"force_password_reset"`
}

type credentialsRepository struct {
	credentialsFile string

	credentials *credentials
}

func (r *credentialsRepository) GetByUsername(
	ctx context.Context, username string,
) (*entities.Credentials, errors.PersistenceError) {
	if username != r.credentials.Username {
		return nil, errors.NewPersistenceNotFoundError(
			"Credentials", nil,
		)
	}

	return &entities.Credentials{
		Username:           r.credentials.Username,
		HashedPassword:     r.credentials.Password,
		ForcePasswordReset: r.credentials.ForcePasswordReset,
	}, nil
}

func (r *credentialsRepository) ChangePassword(
	ctx context.Context, credentials *entities.Credentials,
) errors.PersistenceError {
	if credentials.Username != r.credentials.Username {
		return errors.NewPersistenceNotFoundError(
			"Credentials", nil,
		)
	}

	oldPassword := r.credentials.Password
	r.credentials.Password = credentials.HashedPassword

	if r.credentials.ForcePasswordReset {
		r.credentials.ForcePasswordReset = false
	}

	if err := r.saveCredentials(); err != nil {
		r.credentials.Password = oldPassword
		return errors.NewPersistenceConnectionError(
			"failed to write to disk", err,
		)
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

func NewCredentialsRepository(credentialsFile string) ports.CredentialsRepository {
	file, err := os.Open(credentialsFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var credentials credentials
	if err := json.NewDecoder(file).Decode(&credentials); err != nil {
		panic(err)
	}

	return &credentialsRepository{
		credentials:     &credentials,
		credentialsFile: credentialsFile,
	}
}
