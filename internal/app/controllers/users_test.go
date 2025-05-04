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
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/logger"
)

func Test_Users_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users, log)

	type result struct {
		response serializers.PaginationResponse[serializers.UserSerializer]
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
				users.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.User{
					{
						ID:             uuid.MustParse("10000000-1000-1000-1234-000000000001"),
						IdentityNumber: "PNOEE-60001017869",
						PersonalCode:   "60001017869",
						FirstName:      "EID2016",
						LastName:       "TESTNUMBER",
					},
					{
						ID:             uuid.MustParse("10000000-1000-1000-1234-000000000002"),
						IdentityNumber: "PNOEE-40504040001",
						PersonalCode:   "40504040001",
						FirstName:      "TESTNUMBER",
						LastName:       "ADULT",
					},
				}, uint64(2), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.UserSerializer]{
					Data: []serializers.UserSerializer{
						{
							ID:             uuid.MustParse("10000000-1000-1000-1234-000000000001"),
							IdentityNumber: "PNOEE-60001017869",
							PersonalCode:   "60001017869",
							FirstName:      "EID2016",
							LastName:       "TESTNUMBER",
						},
						{
							ID:             uuid.MustParse("10000000-1000-1000-1234-000000000002"),
							IdentityNumber: "PNOEE-40504040001",
							PersonalCode:   "40504040001",
							FirstName:      "TESTNUMBER",
							LastName:       "ADULT",
						},
					},
					Meta: serializers.PaginationMeta{
						Page:  1,
						Per:   25,
						Total: 2,
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
				users.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.UserSerializer]{
					Data: []serializers.UserSerializer{},
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
				users.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), errors.ErrInvalidArguments)
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
				users.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
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
				users.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/users", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/users", controller.List)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PaginationResponse[serializers.UserSerializer]
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Users_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users, log)

	id := uuid.MustParse("10000000-1000-1000-1234-000000000001")
	roleIds := []uuid.UUID{
		uuid.MustParse("10000000-1000-1000-1000-000000000001"),
		uuid.MustParse("10000000-1000-1000-1000-000000000002"),
	}
	scopeIds := []uuid.UUID{
		uuid.MustParse("10000000-1000-1000-2000-000000000001"),
		uuid.MustParse("10000000-1000-1000-2000-000000000001"),
	}

	type result struct {
		response serializers.UserSerializer
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
				users.EXPECT().FindById(gomock.Any(), id).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
					RoleIDs:        roleIds,
					ScopeIDs:       scopeIds,
				}, nil)
			},
			expected: result{
				response: serializers.UserSerializer{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
					RoleIDs:        roleIds,
					ScopeIDs:       scopeIds,
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Not found",
			before: func() {
				users.EXPECT().FindById(gomock.Any(), id).Return(&models.User{}, errors.ErrRecordNotFound)
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
				users.EXPECT().FindById(gomock.Any(), id).Return(&models.User{}, assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/users/10000000-1000-1000-1234-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/users/{id}", controller.Get)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.UserSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Users_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users, log)

	id := uuid.MustParse("10000000-1000-1000-1234-000000000001")

	type result struct {
		response serializers.UserSerializer
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
				users.EXPECT().Create(gomock.Any(), &models.User{
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}, nil)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
			expected: result{
				response: serializers.UserSerializer{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				},
				status: "201 Created",
				code:   http.StatusCreated,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				users.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty personal code"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Invalid arguments",
			before: func() {
				users.EXPECT().Create(gomock.Any(), &models.User{
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{}, errors.ErrInvalidArguments)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "invalid arguments"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name: "Bad request",
			before: func() {
				users.EXPECT().Create(gomock.Any(), &models.User{
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{}, errors.ErrFailedToCreateRecord)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
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
				users.EXPECT().Create(gomock.Any(), &models.User{
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{}, assert.AnError)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
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

			req := httptest.NewRequest(http.MethodPost, "/api/backoffice/users", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/backoffice/users", controller.Create)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.UserSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Users_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users, log)

	id := uuid.MustParse("10000000-1000-1000-1234-000000000001")

	type result struct {
		response serializers.UserSerializer
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
				users.EXPECT().Update(gomock.Any(), &models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}, nil)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
			expected: result{
				response: serializers.UserSerializer{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
			error: false,
		},
		{
			name: "Invalid params",
			before: func() {
				users.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "empty personal code"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
			error: true,
		},
		{
			name: "Invalid arguments",
			before: func() {
				users.EXPECT().Update(gomock.Any(), &models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{}, errors.ErrInvalidArguments)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
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
				users.EXPECT().Update(gomock.Any(), &models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{}, errors.ErrRecordNotFound)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: errors.ErrRecordNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
			},
		},
		{
			name: "Bad request",
			before: func() {
				users.EXPECT().Update(gomock.Any(), &models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{}, errors.ErrFailedToUpdateRecord)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
			expected: result{
				error:  serializers.ErrorSerializer{Error: "failed to update record"},
				status: "422 Unprocessable Entity",
				code:   http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Error",
			before: func() {
				users.EXPECT().Update(gomock.Any(), &models.User{
					ID:             id,
					IdentityNumber: "PNOEE-60001017869",
					PersonalCode:   "60001017869",
					FirstName:      "EID2016",
					LastName:       "TESTNUMBER",
				}).Return(&models.User{}, assert.AnError)
			},
			body: strings.NewReader(`{"identity_number": "PNOEE-60001017869", "personal_code": "60001017869", "first_name": "EID2016", "last_name": "TESTNUMBER"}`),
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

			req := httptest.NewRequest(http.MethodPut, "/api/backoffice/users/10000000-1000-1000-1234-000000000001", tt.body)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Put("/api/backoffice/users/{id}", controller.Update)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.UserSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, response)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Users_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	users := services.NewMockUsers(ctrl)
	controller := NewUsersController(users, log)

	id := uuid.MustParse("10000000-1000-1000-1234-000000000001")

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
				users.EXPECT().Delete(gomock.Any(), id).Return(true, nil)
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
				users.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(false, errors.ErrInvalidArguments)
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
				users.EXPECT().Delete(gomock.Any(), id).Return(false, errors.ErrRecordNotFound)
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
				users.EXPECT().Delete(gomock.Any(), id).Return(false, errors.ErrFailedToDeleteRecord)
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
				users.EXPECT().Delete(gomock.Any(), id).Return(false, assert.AnError)
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

			req := httptest.NewRequest(http.MethodDelete, "/api/backoffice/users/10000000-1000-1000-1234-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/api/backoffice/users/{id}", controller.Delete)
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
