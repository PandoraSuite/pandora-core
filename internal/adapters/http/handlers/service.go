package handlers

import "github.com/MAD-py/pandora-core/internal/ports/inbound"

type APIKeyHandler struct {
	apiKeyUsecase inbound.APIKeyPort
}
