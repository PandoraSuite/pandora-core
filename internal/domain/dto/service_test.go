package dto

import (
	"testing"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

func TestServiceFilterValidation(t *testing.T) {
	tests := []struct {
		name       string
		dto        ServiceFilter
		wantErr    bool
		wantLocErr string
	}{
		{
			name:    "EmptyStatus",
			dto:     ServiceFilter{},
			wantErr: false,
		},
		{
			name:    "ValidStatus",
			dto:     ServiceFilter{Status: enums.ServiceStatusEnabled},
			wantErr: false,
		},
		{
			name:       "InvalidStatus",
			dto:        ServiceFilter{Status: "invalid_status"},
			wantErr:    true,
			wantLocErr: "status",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := v.ValidateStruct(test.dto, map[string]string{})

			if !test.wantErr {
				if err != nil {
					t.Errorf("got %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Error("got nil, want error")
				return
			}

			if errors.CodeValidationFailed != err.Code() {
				t.Errorf(
					"got %s code, want %s",
					err.Code(),
					errors.CodeValidationFailed,
				)
				return
			}

			vErr, ok := err.(*errors.AttributeError)
			if !ok {
				t.Errorf("got %T error, want AttributeError", err)
				return
			}

			if vErr.Loc() != test.wantLocErr {
				t.Errorf("got %s loc, want %s", vErr.Loc(), test.wantLocErr)
				return
			}
		})
	}
}

func TestServiceCreateValidation(t *testing.T) {
	tests := []struct {
		name       string
		dto        ServiceCreate
		wantErr    bool
		wantLocErr string
	}{
		{
			name:    "ValidCreate",
			dto:     ServiceCreate{Name: "Service", Version: "1.0.0"},
			wantErr: false,
		},
		{
			name:       "MissingName",
			dto:        ServiceCreate{Version: "1.0.0"},
			wantErr:    true,
			wantLocErr: "name",
		},
		{
			name:       "MissingVersion",
			dto:        ServiceCreate{Name: "Service"},
			wantErr:    true,
			wantLocErr: "version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := v.ValidateStruct(test.dto, map[string]string{})

			if !test.wantErr {
				if err != nil {
					t.Errorf("got %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Error("got nil, want error")
				return
			}

			if errors.CodeValidationFailed != err.Code() {
				t.Errorf(
					"got %s code, want %s",
					err.Code(),
					errors.CodeValidationFailed,
				)
				return
			}

			vErr, ok := err.(*errors.AttributeError)
			if !ok {
				t.Errorf("got %T error, want AttributeError", err)
				return
			}

			if vErr.Loc() != test.wantLocErr {
				t.Errorf("got %s loc, want %s", vErr.Loc(), test.wantLocErr)
				return
			}
		})
	}
}
