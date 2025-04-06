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
	"loki-backoffice/internal/config/logger"
)

type UsersController interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type usersController struct {
	users services.Users
	log   *logger.Logger
}

func NewUsersController(users services.Users, log *logger.Logger) UsersController {
	return &usersController{
		users: users,
		log:   log,
	}
}

//nolint:dupl
func (c *usersController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pagination := services.NewPagination(r)
	rows, total, err := c.users.List(r.Context(), pagination)
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

	collection := make([]serializers.UserSerializer, 0, len(rows))

	for _, user := range rows {
		collection = append(collection, serializers.UserSerializer{
			ID:             user.ID,
			IdentityNumber: user.IdentityNumber,
			PersonalCode:   user.PersonalCode,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
		})
	}

	response := serializers.PaginationResponse[serializers.UserSerializer]{
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
func (c *usersController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Invalid UUID format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrInvalidArguments.Error()})
		return
	}

	record, err := c.users.FindById(r.Context(), id)

	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Failed to get user")

		switch {
		case errors.Is(err, errors.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.UserSerializer{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
		RoleIDs:        record.RoleIDs,
		ScopeIDs:       record.ScopeIDs,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *usersController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.UserRequest
	if err := params.Validate(r.Body); err != nil {
		c.log.Error().Err(err).Str("identity_number", params.IdentityNumber).Msg("Failed to create user")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.users.Create(r.Context(), &models.User{
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		RoleIDs:        params.RoleIDs,
		ScopeIDs:       params.ScopeIDs,
	})
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to create user")

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

	response := serializers.UserSerializer{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

////nolint:dupl
func (c *usersController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Invalid UUID format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrInvalidArguments.Error()})
		return
	}

	var params dto.UserRequest
	if err = params.Validate(r.Body); err != nil {
		c.log.Error().Err(err).Str("identity_number", params.IdentityNumber).Msg("Failed to create user")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.users.Update(r.Context(), &models.User{
		ID:             id,
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		RoleIDs:        params.RoleIDs,
		ScopeIDs:       params.ScopeIDs,
	})
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Failed to update user")

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

	response := serializers.UserSerializer{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *usersController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Invalid UUID format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrInvalidArguments.Error()})
		return
	}

	_, err = c.users.Delete(r.Context(), id)
	if err != nil {
		c.log.Error().Err(err).Str("id", id.String()).Msg("Failed to delete user")

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
