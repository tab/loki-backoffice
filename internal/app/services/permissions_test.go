package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/app/models"
	"loki-backoffice/internal/app/rpcs"
	proto "loki-backoffice/internal/app/rpcs/proto/sso/v1"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/logger"
)

func Test_Permissions_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	mockClient := rpcs.NewMockPermissionServiceClient(ctrl)
	service := NewPermissions(mockClient, log)

	tests := []struct {
		name     string
		before   func()
		expected []models.Permission
		total    uint64
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(&proto.ListPermissionsResponse{
					Data: []*proto.Permission{
						{
							Id:          "10000000-1000-1000-3000-000000000001",
							Name:        "read:self",
							Description: "Read own data",
						},
						{
							Id:          "10000000-1000-1000-3000-000000000002",
							Name:        "write:self",
							Description: "Write own data",
						},
					},
					Meta: &proto.PaginationMeta{
						Page:  1,
						Per:   10,
						Total: 2,
					},
				}, nil)
			},
			expected: []models.Permission{
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					Name:        "read:self",
					Description: "Read own data",
				},
				{
					ID:          uuid.MustParse("10000000-1000-1000-3000-000000000002"),
					Name:        "write:self",
					Description: "Write own data",
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
				}).Return(&proto.ListPermissionsResponse{
					Data: []*proto.Permission{},
					Meta: &proto.PaginationMeta{
						Page:  1,
						Per:   10,
						Total: 0,
					},
				}, nil)
			},
			expected: []models.Permission{},
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
				assert.Equal(t, tt.expected, result)
				assert.Equal(t, tt.total, total)
			}
		})
	}
}

func Test_Permissions_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	mockClient := rpcs.NewMockPermissionServiceClient(ctrl)
	service := NewPermissions(mockClient, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

	tests := []struct {
		name     string
		before   func()
		expected *models.Permission
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetPermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(&proto.GetPermissionResponse{
					Data: &proto.Permission{
						Id:          "10000000-1000-1000-3000-000000000001",
						Name:        "read:self",
						Description: "Read own data",
					},
				}, nil)
			},
			expected: &models.Permission{
				ID:          id,
				Name:        "read:self",
				Description: "Read own data",
			},
			error: nil,
		},
		{
			name: "InvalidArgument status code",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetPermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: nil,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "NotFound status code",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetPermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expected: nil,
			error:    errors.ErrRecordNotFound,
		},
		{
			name: "Internal status code",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetPermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: nil,
			error:    errors.ErrFailedToFetchResults,
		},
		{
			name: "Error",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetPermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.FindById(ctx, id)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.error, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Permissions_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	mockClient := rpcs.NewMockPermissionServiceClient(ctrl)
	service := NewPermissions(mockClient, log)

	params := &models.Permission{
		Name:        models.ReadSelfType,
		Description: "Read own data",
	}

	tests := []struct {
		name     string
		before   func()
		expected *models.Permission
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreatePermissionRequest{
					Name:        "read:self",
					Description: "Read own data",
				}).Return(&proto.CreatePermissionResponse{
					Data: &proto.Permission{
						Id:          "10000000-1000-1000-3000-000000000001",
						Name:        "read:self",
						Description: "Read own data",
					},
				}, nil)
			},
			expected: &models.Permission{
				ID:          uuid.MustParse("10000000-1000-1000-3000-000000000001"),
				Name:        "read:self",
				Description: "Read own data",
			},
			error: nil,
		},
		{
			name: "InvalidArgument status code",
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreatePermissionRequest{
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: nil,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "Internal status code",
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreatePermissionRequest{
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: nil,
			error:    errors.ErrFailedToCreateRecord,
		},
		{
			name: "Error",
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreatePermissionRequest{
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Create(ctx, params)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.error, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Permissions_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	mockClient := rpcs.NewMockPermissionServiceClient(ctrl)
	service := NewPermissions(mockClient, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")
	params := &models.Permission{
		ID:          id,
		Name:        models.ReadSelfType,
		Description: "Read own data",
	}

	tests := []struct {
		name     string
		before   func()
		expected *models.Permission
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdatePermissionRequest{
					Id:          "10000000-1000-1000-3000-000000000001",
					Name:        "read:self",
					Description: "Read own data",
				}).Return(&proto.UpdatePermissionResponse{
					Data: &proto.Permission{
						Id:          "10000000-1000-1000-3000-000000000001",
						Name:        "read:self",
						Description: "Read own data",
					},
				}, nil)
			},
			expected: &models.Permission{
				ID:          id,
				Name:        "read:self",
				Description: "Read own data",
			},
			error: nil,
		},
		{
			name: "InvalidArgument status code",
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdatePermissionRequest{
					Id:          "10000000-1000-1000-3000-000000000001",
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: nil,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "NotFound status code",
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdatePermissionRequest{
					Id:          "10000000-1000-1000-3000-000000000001",
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expected: nil,
			error:    errors.ErrRecordNotFound,
		},
		{
			name: "Internal status code",
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdatePermissionRequest{
					Id:          "10000000-1000-1000-3000-000000000001",
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: nil,
			error:    errors.ErrFailedToUpdateRecord,
		},
		{
			name: "Error",
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdatePermissionRequest{
					Id:          "10000000-1000-1000-3000-000000000001",
					Name:        "read:self",
					Description: "Read own data",
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, params)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.error, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Permissions_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	ctx := context.Background()
	mockClient := rpcs.NewMockPermissionServiceClient(ctrl)
	service := NewPermissions(mockClient, log)

	id := uuid.MustParse("10000000-1000-1000-3000-000000000001")

	tests := []struct {
		name     string
		before   func()
		expected bool
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeletePermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(&emptypb.Empty{}, nil)
			},
			expected: true,
			error:    nil,
		},
		{
			name: "InvalidArgument status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeletePermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: false,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "NotFound status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeletePermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expected: false,
			error:    errors.ErrRecordNotFound,
		},
		{
			name: "Internal status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeletePermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: false,
			error:    errors.ErrFailedToDeleteRecord,
		},
		{
			name: "Error",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeletePermissionRequest{
					Id: "10000000-1000-1000-3000-000000000001",
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
