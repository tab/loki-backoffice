package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/app/models"
	"loki-backoffice/internal/app/serializers"
	"loki-backoffice/internal/app/services"
	"loki-backoffice/pkg/logger"
)

func Test_Roles_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	controller := NewRolesController(roles, log)

	type result struct {
		response serializers.PaginationResponse[serializers.RoleSerializer]
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		before   func()
		expected result
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.Role{
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
						Name:        "admin",
						Description: "Admin role",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
						Name:        "manager",
						Description: "Manager role",
					},
					{
						ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
						Name:        "user",
						Description: "User role",
					},
				}, uint64(3), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.RoleSerializer]{
					Data: []serializers.RoleSerializer{
						{
							ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
							Name:        "admin",
							Description: "Admin role",
						},
						{
							ID:          uuid.MustParse("10000000-1000-1000-1000-000000000002"),
							Name:        "manager",
							Description: "Manager role",
						},
						{
							ID:          uuid.MustParse("10000000-1000-1000-1000-000000000003"),
							Name:        "user",
							Description: "User role",
						},
					},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: 3,
					},
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Empty",
			before: func() {
				roles.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.RoleSerializer]{
					Data: []serializers.RoleSerializer{},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: 0,
					},
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				roles.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), errors.ErrInvalidArguments)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "invalid arguments"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Bad request",
			before: func() {
				roles.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "failed to fetch results"},
				status: "503 Service Unavailable",
				code:   http.StatusServiceUnavailable,
			},
			error: true,
		},
		{
			name: "Error",
			before: func() {
				roles.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), assert.AnError)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/roles", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/roles", controller.List)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PaginationResponse[serializers.RoleSerializer]
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Roles_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	controller := NewRolesController(roles, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

	type result struct {
		response serializers.RoleSerializer
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		before   func()
		expected result
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().FindById(gomock.Any(), uuid.MustParse("10000000-1000-1000-1000-000000000001")).Return(&models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				}, nil)
			},
			expected: result{
				response: serializers.RoleSerializer{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Not found",
			before: func() {
				roles.EXPECT().FindById(gomock.Any(), id).Return(&models.Role{}, errors.ErrRecordNotFound)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrRecordNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
		},
		{
			name: "Error",
			before: func() {
				roles.EXPECT().FindById(gomock.Any(), id).Return(&models.Role{}, assert.AnError)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/roles/10000000-1000-1000-1000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/roles/{id}", controller.Get)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.RoleSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Roles_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	controller := NewRolesController(roles, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

	type result struct {
		response serializers.RoleSerializer
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		before   func()
		body     io.Reader
		expected result
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().Create(gomock.Any(), &models.Role{
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				}, nil)
			},
			body: strings.NewReader(`{"name": "admin", "description" :"Admin role"}`),
			expected: result{
				response: serializers.RoleSerializer{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				},
				status: "201 Created",
				code:   http.StatusCreated,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				roles.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"name": "admin"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty description"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Invalid arguments",
			before: func() {
				roles.EXPECT().Create(gomock.Any(), &models.Role{
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{}, errors.ErrInvalidArguments)
			},
			body: strings.NewReader(`{"name": "admin", "description": "Admin role"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "invalid arguments"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name: "Bad request",
			before: func() {
				roles.EXPECT().Create(gomock.Any(), &models.Role{
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{}, errors.ErrFailedToCreateRecord)
			},
			body: strings.NewReader(`{"name": "admin", "description": "Admin role"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "failed to create record"},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
			error: true,
		},
		{
			name: "Error",
			before: func() {
				roles.EXPECT().Create(gomock.Any(), &models.Role{
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{}, assert.AnError)
			},
			body: strings.NewReader(`{"name": "admin", "description": "Admin role"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodPost, "/api/backoffice/roles", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/backoffice/roles", controller.Create)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.RoleSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Roles_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	controller := NewRolesController(roles, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

	type result struct {
		response serializers.RoleSerializer
		error    serializers.ErrorSerializer
		status   string
		code     int
	}

	tests := []struct {
		name     string
		before   func()
		body     io.Reader
		expected result
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().Update(gomock.Any(), &models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				}, nil)
			},
			body: strings.NewReader(`{"name": "admin", "description": "Admin role"}`),
			expected: result{
				response: serializers.RoleSerializer{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				roles.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"name": "admin"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty description"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Invalid arguments",
			before: func() {
				roles.EXPECT().Update(gomock.Any(), &models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{}, errors.ErrInvalidArguments)
			},
			body: strings.NewReader(`{"name": "admin", "description": "Admin role"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "invalid arguments"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Not found",
			before: func() {
				roles.EXPECT().Update(gomock.Any(), &models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{}, errors.ErrRecordNotFound)
			},
			body: strings.NewReader(`{"name": "admin", "description": "Admin role"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrRecordNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
		},
		{
			name: "Bad request",
			before: func() {
				roles.EXPECT().Update(gomock.Any(), &models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{}, errors.ErrFailedToUpdateRecord)
			},
			body: strings.NewReader(`{"name": "admin", "description": "Admin role"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "failed to update record"},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Error",
			before: func() {
				roles.EXPECT().Update(gomock.Any(), &models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Admin role",
				}).Return(&models.Role{}, assert.AnError)
			},
			body: strings.NewReader(`{"name": "admin", "description": "Admin role"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodPut, "/api/backoffice/roles/10000000-1000-1000-1000-000000000001", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Put("/api/backoffice/roles/{id}", controller.Update)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.RoleSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Roles_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roles := services.NewMockRoles(ctrl)
	log := logger.NewLogger()
	controller := NewRolesController(roles, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

	type result struct {
		error  serializers.ErrorSerializer
		status string
		code   int
	}

	tests := []struct {
		name     string
		before   func()
		expected result
		error    bool
	}{
		{
			name: "Success",
			before: func() {
				roles.EXPECT().Delete(gomock.Any(), id).Return(true, nil)
			},
			expected: result{
				status: "204 No Content",
				code:   http.StatusNoContent,
			},
			error: false,
		},
		{
			name: "Invalid arguments",
			before: func() {
				roles.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(false, errors.ErrInvalidArguments)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "invalid arguments"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Not found",
			before: func() {
				roles.EXPECT().Delete(gomock.Any(), id).Return(false, errors.ErrRecordNotFound)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrRecordNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
		},
		{
			name: "Bad request",
			before: func() {
				roles.EXPECT().Delete(gomock.Any(), id).Return(false, errors.ErrFailedToDeleteRecord)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: "failed to delete record"},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Error",
			before: func() {
				roles.EXPECT().Delete(gomock.Any(), id).Return(false, assert.AnError)
			},
			expected: result{
				error:  serializers.ErrorSerializer{Error: assert.AnError.Error()},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
			},
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodDelete, "/api/backoffice/roles/10000000-1000-1000-1000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/api/backoffice/roles/{id}", controller.Delete)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}
