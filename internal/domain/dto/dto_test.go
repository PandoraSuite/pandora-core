package dto

import (
	"os"
	"testing"

	"github.com/MAD-py/pandora-core/internal/validator"
)

var v validator.Validator

func TestMain(m *testing.M) {
	v = validator.NewValidator()
	os.Exit(m.Run())
}
