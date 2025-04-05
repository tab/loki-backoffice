package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/app/models"
	"loki-backoffice/internal/app/models/dto"
	"loki-backoffice/internal/app/serializers"
	"loki-backoffice/internal/app/services"
	"loki-backoffice/pkg/logger"
)

type ScopesController interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type scopesController struct {
	scopes services.Scopes
	log    *logger.Logger
}

func NewScopesController(scopes services.Scopes, log *logger.Logger) ScopesController {
	return &scopesController{
		scopes: scopes,
		log:    log,
	}
}

//nolint:dupl
func (c *scopesController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pagination := services.NewPagination(r)
	rows, total, err := c.scopes.List(r.Context(), pagination)
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

	collection := make([]serializers.ScopeSerializer, 0, len(rows))

	for _, scope := range rows {
		collection = append(collection, serializers.ScopeSerializer{
			ID:          scope.ID,
			Name:        scope.Name,
			Description: scope.Description,
		})
	}

	response := serializers.PaginationResponse[serializers.ScopeSerializer]{
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
func (c *scopesController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Invalid UUID format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrInvalidArguments.Error()})
		return
	}

	record, err := c.scopes.FindById(r.Context(), id)

	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Failed to get scope")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.ScopeSerializer{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *scopesController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.ScopeRequest
	if err := params.Validate(r.Body); err != nil {
		c.log.Error().Err(err).Str("name", params.Name).Msg("Failed to create scope")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.scopes.Create(r.Context(), &models.Scope{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to create scope")

		switch {
		case errors.Is(err, errors.ErrInvalidArguments):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, errors.ErrFailedToCreateRecord):
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.ScopeSerializer{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

////nolint:dupl
func (c *scopesController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Invalid UUID format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrInvalidArguments.Error()})
		return
	}

	var params dto.ScopeRequest
	if err = params.Validate(r.Body); err != nil {
		c.log.Error().Err(err).Str("name", params.Name).Msg("Failed to create scope")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.scopes.Update(r.Context(), &models.Scope{
		ID:          id,
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Failed to update scope")

		switch {
		case errors.Is(err, errors.ErrInvalidArguments):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, errors.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, errors.ErrFailedToUpdateRecord):
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.ScopeSerializer{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *scopesController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Invalid UUID format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrInvalidArguments.Error()})
		return
	}

	_, err = c.scopes.Delete(r.Context(), id)
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Failed to delete scope")

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
