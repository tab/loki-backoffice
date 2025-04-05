package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func Test_Tokens_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokens := services.NewMockTokens(ctrl)
	log := logger.NewLogger()
	controller := NewTokensController(tokens, log)

	accessTokenExp := time.Now().Add(models.AccessTokenExp)
	refreshTokenExp := time.Now().Add(models.RefreshTokenExp)

	type result struct {
		response serializers.PaginationResponse[serializers.TokenSerializer]
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
				tokens.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.Token{
					{
						ID:        uuid.MustParse("10000000-1000-1000-6000-000000000001"),
						UserId:    uuid.MustParse("10000000-1000-1000-1234-000000000001"),
						Type:      models.AccessTokenType,
						Value:     "access-token-value",
						ExpiresAt: accessTokenExp,
					},
					{
						ID:        uuid.MustParse("10000000-1000-1000-6000-000000000002"),
						UserId:    uuid.MustParse("10000000-1000-1000-1234-000000000002"),
						Type:      models.RefreshTokenType,
						Value:     "refresh-token-value",
						ExpiresAt: refreshTokenExp,
					},
				}, uint64(2), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.TokenSerializer]{
					Data: []serializers.TokenSerializer{
						{
							ID:        uuid.MustParse("10000000-1000-1000-6000-000000000001"),
							UserId:    uuid.MustParse("10000000-1000-1000-1234-000000000001"),
							Type:      models.AccessTokenType,
							Value:     "access-token-value",
							ExpiresAt: accessTokenExp,
						},
						{
							ID:        uuid.MustParse("10000000-1000-1000-6000-000000000002"),
							UserId:    uuid.MustParse("10000000-1000-1000-1234-000000000002"),
							Type:      models.RefreshTokenType,
							Value:     "refresh-token-value",
							ExpiresAt: refreshTokenExp,
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
				tokens.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), nil)
			},
			expected: result{
				response: serializers.PaginationResponse[serializers.TokenSerializer]{
					Data: []serializers.TokenSerializer{},
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
				tokens.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), errors.ErrInvalidArguments)
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
				tokens.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), errors.ErrFailedToFetchResults)
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
				tokens.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, uint64(0), assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/backoffice/tokens", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/backoffice/tokens", controller.List)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.error {
				var response serializers.ErrorSerializer
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error, response)
			} else {
				var response serializers.PaginationResponse[serializers.TokenSerializer]
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				for i, item := range response.Data {
					assert.Equal(t, tt.expected.response.Data[i].ID, item.ID)
					assert.Equal(t, tt.expected.response.Data[i].UserId, item.UserId)
					assert.Equal(t, tt.expected.response.Data[i].Type, item.Type)
					assert.Equal(t, tt.expected.response.Data[i].Value, item.Value)
					assert.Equal(t, tt.expected.response.Data[i].ExpiresAt.Unix(), item.ExpiresAt.Unix())
				}
				assert.Equal(t, tt.expected.response.Meta, response.Meta)
			}

			assert.Equal(t, tt.expected.code, resp.StatusCode)
			assert.Equal(t, tt.expected.status, resp.Status)
		})
	}
}

func Test_Tokens_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokens := services.NewMockTokens(ctrl)
	log := logger.NewLogger()
	controller := NewTokensController(tokens, log)

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
				tokens.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-6000-000000000001")).Return(true, nil)
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
				tokens.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(false, errors.ErrInvalidArguments)
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
				tokens.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-6000-000000000001")).Return(false, errors.ErrRecordNotFound)
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
				tokens.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-6000-000000000001")).Return(false, errors.ErrFailedToDeleteRecord)
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
				tokens.EXPECT().Delete(gomock.Any(), uuid.MustParse("10000000-1000-1000-6000-000000000001")).Return(false, assert.AnError)
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

			req := httptest.NewRequest(http.MethodDelete, "/api/backoffice/tokens/10000000-1000-1000-6000-000000000001", nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Delete("/api/backoffice/tokens/{id}", controller.Delete)
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
