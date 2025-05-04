package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/app/models"
	"loki-backoffice/internal/app/rpcs"
	proto "loki-backoffice/internal/app/rpcs/proto/sso/v1"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/logger"
)

func Test_Tokens_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	mockClient := rpcs.NewMockTokenServiceClient(ctrl)
	service := NewTokens(mockClient, log)

	accessTokenExp := time.Now().Add(models.AccessTokenExp)
	refreshTokenExp := time.Now().Add(models.RefreshTokenExp)

	tests := []struct {
		name     string
		before   func()
		expected []models.Token
		total    uint64
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(&proto.ListTokensResponse{
					Data: []*proto.Token{
						{
							Id:        "10000000-1000-1000-6000-000000000001",
							UserId:    "10000000-1000-1000-1234-000000000001",
							Type:      models.AccessTokenType,
							Value:     "access-token-value",
							ExpiresAt: timestamppb.New(accessTokenExp),
						},
						{
							Id:        "10000000-1000-1000-6000-000000000002",
							UserId:    "10000000-1000-1000-1234-000000000002",
							Type:      models.RefreshTokenType,
							Value:     "refresh-token-value",
							ExpiresAt: timestamppb.New(refreshTokenExp),
						},
					},
					Meta: &proto.PaginationMeta{
						Page:  1,
						Per:   10,
						Total: 2,
					},
				}, nil)
			},
			expected: []models.Token{
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
			total: 2,
			error: nil,
		},
		{
			name: "Empty",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(&proto.ListTokensResponse{
					Data: []*proto.Token{},
					Meta: &proto.PaginationMeta{
						Page:  1,
						Per:   10,
						Total: 0,
					},
				}, nil)
			},
			expected: []models.Token{},
			total:    0,
			error:    nil,
		},
		{
			name: "InvalidArgument status code",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: nil,
			total:    0,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "Unavailable status code",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(nil, status.Error(codes.Unavailable, "service unavailable"))
			},
			expected: nil,
			total:    0,
			error:    errors.ErrFailedToFetchResults,
		},
		{
			name: "Internal status code",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: nil,
			total:    0,
			error:    errors.ErrFailedToFetchResults,
		},
		{
			name: "Error",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			total:    0,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, total, err := service.List(ctx, &Pagination{
				Page:    uint64(1),
				PerPage: uint64(10),
			})

			if tt.error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.error, err)
				assert.Nil(t, result)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)

				for i, expected := range tt.expected {
					assert.Equal(t, expected.ID, result[i].ID)
					assert.Equal(t, expected.UserId, result[i].UserId)
					assert.Equal(t, expected.Type, result[i].Type)
					assert.Equal(t, expected.Value, result[i].Value)
					assert.Equal(t, expected.ExpiresAt.Unix(), result[i].ExpiresAt.Unix())
				}
				assert.Equal(t, tt.total, total)
			}
		})
	}
}

func Test_Tokens_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	mockClient := rpcs.NewMockTokenServiceClient(ctrl)
	service := NewTokens(mockClient, log)

	id := uuid.MustParse("10000000-1000-1000-6000-000000000001")

	tests := []struct {
		name     string
		before   func()
		expected bool
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteTokenRequest{
					Id: "10000000-1000-1000-6000-000000000001",
				}).Return(&emptypb.Empty{}, nil)
			},
			expected: true,
			error:    nil,
		},
		{
			name: "InvalidArgument status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteTokenRequest{
					Id: "10000000-1000-1000-6000-000000000001",
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: false,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "NotFound status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteTokenRequest{
					Id: "10000000-1000-1000-6000-000000000001",
				}).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expected: false,
			error:    errors.ErrRecordNotFound,
		},
		{
			name: "Internal status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteTokenRequest{
					Id: "10000000-1000-1000-6000-000000000001",
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: false,
			error:    errors.ErrFailedToDeleteRecord,
		},
		{
			name: "Error",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteTokenRequest{
					Id: "10000000-1000-1000-6000-000000000001",
				}).Return(nil, assert.AnError)
			},
			expected: false,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			ok, err := service.Delete(ctx, id)

			if tt.error != nil {
				assert.Equal(t, tt.expected, ok)
				assert.Error(t, err)
				assert.Equal(t, tt.error, err)
			} else {
				assert.Equal(t, tt.expected, ok)
				assert.NoError(t, err)
			}
		})
	}
}
