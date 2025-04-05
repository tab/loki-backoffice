package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/app/serializers"
	"loki-backoffice/internal/app/services"
	"loki-backoffice/pkg/logger"
)

type TokensController interface {
	List(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type tokensController struct {
	tokens services.Tokens
	log    *logger.Logger
}

func NewTokensController(tokens services.Tokens, log *logger.Logger) TokensController {
	return &tokensController{
		tokens: tokens,
		log:    log,
	}
}

//nolint:dupl
func (c *tokensController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pagination := services.NewPagination(r)
	rows, total, err := c.tokens.List(r.Context(), pagination)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrInvalidArguments):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, errors.ErrFailedToFetchResults):
			w.WriteHeader(http.StatusServiceUnavailable)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	collection := make([]serializers.TokenSerializer, 0, len(rows))

	for _, token := range rows {
		collection = append(collection, serializers.TokenSerializer{
			ID:        token.ID,
			UserId:    token.UserId,
			Type:      token.Type,
			Value:     token.Value,
			ExpiresAt: token.ExpiresAt,
		})
	}

	response := serializers.PaginationResponse[serializers.TokenSerializer]{
		Data: collection,
		Meta: serializers.PaginationMeta{
			Page:  pagination.Page,
			Per:   pagination.PerPage,
			Total: total,
		},
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *tokensController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Invalid UUID format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrInvalidArguments.Error()})
		return
	}

	_, err = c.tokens.Delete(r.Context(), id)
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Failed to delete token")

		switch {
		case errors.Is(err, errors.ErrInvalidArguments):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, errors.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, errors.ErrFailedToDeleteRecord):
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
