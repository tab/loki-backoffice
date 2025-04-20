package dto

import (
	"encoding/json"
	"io"
	"strings"

	"loki-backoffice/internal/app/errors"
)

type ScopeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (params *ScopeRequest) Validate(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(params); err != nil {
		return err
	}

	params.Name = strings.TrimSpace(params.Name)
	if params.Name == "" {
		return errors.ErrEmptyName
	}

	params.Description = strings.TrimSpace(params.Description)
	if params.Description == "" {
		return errors.ErrEmptyDescription
	}

	return nil
}
