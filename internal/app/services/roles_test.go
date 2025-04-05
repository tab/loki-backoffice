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
	"loki-backoffice/pkg/logger"
)

func Test_Roles_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockClient := rpcs.NewMockRoleServiceClient(ctrl)
	log := logger.NewLogger()
	service := NewRoles(mockClient, log)

	tests := []struct {
		name     string
		before   func()
		expected []models.Role
		total    uint64
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(&proto.ListRolesResponse{
					Data: []*proto.Role{
						{
							Id:          "10000000-1000-1000-1000-000000000001",
							Name:        "admin",
							Description: "Admin role",
						},
						{
							Id:          "10000000-1000-1000-1000-000000000002",
							Name:        "manager",
							Description: "Manager role",
						},
						{
							Id:          "10000000-1000-1000-1000-000000000003",
							Name:        "user",
							Description: "User role",
						},
					},
					Meta: &proto.PaginationMeta{
						Page:  1,
						Per:   10,
						Total: 3,
					},
				}, nil)
			},
			expected: []models.Role{
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
			total: 3,
			error: nil,
		},
		{
			name: "Empty",
			before: func() {
				mockClient.EXPECT().List(ctx, &proto.PaginatedListRequest{
					Limit:  uint64(1),
					Offset: uint64(10),
				}).Return(&proto.ListRolesResponse{
					Data: []*proto.Role{},
					Meta: &proto.PaginationMeta{
						Page:  1,
						Per:   10,
						Total: 0,
					},
				}, nil)
			},
			expected: []models.Role{},
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

func Test_Roles_FindById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockClient := rpcs.NewMockRoleServiceClient(ctrl)
	log := logger.NewLogger()
	service := NewRoles(mockClient, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")
	permissionIds := []uuid.UUID{
		uuid.MustParse("10000000-1000-1000-3000-000000000001"),
		uuid.MustParse("10000000-1000-1000-3000-000000000002"),
	}
	permissionIdsStr := []string{
		"10000000-1000-1000-3000-000000000001",
		"10000000-1000-1000-3000-000000000002",
	}

	tests := []struct {
		name     string
		before   func()
		expected *models.Role
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
				}).Return(&proto.GetRoleResponse{
					Data: &proto.Role{
						Id:            "10000000-1000-1000-1000-000000000001",
						Name:          "admin",
						Description:   "Admin role",
						PermissionIds: permissionIdsStr,
					},
				}, nil)
			},
			expected: &models.Role{
				ID:            id,
				Name:          "admin",
				Description:   "Admin role",
				PermissionIDs: permissionIds,
			},
			error: nil,
		},
		{
			name: "InvalidArgument status code",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: nil,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "NotFound status code",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
				}).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expected: nil,
			error:    errors.ErrRecordNotFound,
		},
		{
			name: "Internal status code",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: nil,
			error:    errors.ErrFailedToFetchResults,
		},
		{
			name: "Error",
			before: func() {
				mockClient.EXPECT().Get(ctx, &proto.GetRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
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

func Test_Roles_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockClient := rpcs.NewMockRoleServiceClient(ctrl)
	log := logger.NewLogger()
	service := NewRoles(mockClient, log)

	permissionIds := []uuid.UUID{
		uuid.MustParse("10000000-1000-1000-3000-000000000001"),
		uuid.MustParse("10000000-1000-1000-3000-000000000002"),
	}

	tests := []struct {
		name     string
		params   *models.Role
		before   func()
		expected *models.Role
		error    error
	}{
		{
			name: "Success",
			params: &models.Role{
				Name:          models.AdminRoleType,
				Description:   "Admin role",
				PermissionIDs: permissionIds,
			},
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreateRoleRequest{
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{"10000000-1000-1000-3000-000000000001", "10000000-1000-1000-3000-000000000002"},
				}).Return(&proto.CreateRoleResponse{
					Data: &proto.Role{
						Id:            "10000000-1000-1000-1000-000000000001",
						Name:          "admin",
						Description:   "Admin role",
						PermissionIds: []string{"10000000-1000-1000-3000-000000000001", "10000000-1000-1000-3000-000000000002"},
					},
				}, nil)
			},
			expected: &models.Role{
				ID:          uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:        "admin",
				Description: "Admin role",
				PermissionIDs: []uuid.UUID{
					uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					uuid.MustParse("10000000-1000-1000-3000-000000000002"),
				},
			},
			error: nil,
		},
		{
			name: "No permission IDs",
			params: &models.Role{
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreateRoleRequest{
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(&proto.CreateRoleResponse{
					Data: &proto.Role{
						Id:            "10000000-1000-1000-1000-000000000001",
						Name:          "admin",
						Description:   "Admin role",
						PermissionIds: []string{},
					},
				}, nil)
			},
			expected: &models.Role{
				ID:            uuid.MustParse("10000000-1000-1000-1000-000000000001"),
				Name:          "admin",
				Description:   "Admin role",
				PermissionIDs: []uuid.UUID{},
			},
			error: nil,
		},
		{
			name: "InvalidArgument status code",
			params: &models.Role{
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreateRoleRequest{
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: nil,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "Internal status code",
			params: &models.Role{
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreateRoleRequest{
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: nil,
			error:    errors.ErrFailedToCreateRecord,
		},
		{
			name: "Error",
			params: &models.Role{
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Create(ctx, &proto.CreateRoleRequest{
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Create(ctx, tt.params)

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

func Test_Roles_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockClient := rpcs.NewMockRoleServiceClient(ctrl)
	log := logger.NewLogger()
	service := NewRoles(mockClient, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")
	permissionIds := []uuid.UUID{
		uuid.MustParse("10000000-1000-1000-3000-000000000001"),
		uuid.MustParse("10000000-1000-1000-3000-000000000002"),
	}

	tests := []struct {
		name     string
		params   *models.Role
		before   func()
		expected *models.Role
		error    error
	}{
		{
			name: "Success",
			params: &models.Role{
				ID:            id,
				Name:          models.AdminRoleType,
				Description:   "Admin role",
				PermissionIDs: permissionIds,
			},
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdateRoleRequest{
					Id:            "10000000-1000-1000-1000-000000000001",
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{"10000000-1000-1000-3000-000000000001", "10000000-1000-1000-3000-000000000002"},
				}).Return(&proto.UpdateRoleResponse{
					Data: &proto.Role{
						Id:            "10000000-1000-1000-1000-000000000001",
						Name:          "admin",
						Description:   "Admin role",
						PermissionIds: []string{"10000000-1000-1000-3000-000000000001", "10000000-1000-1000-3000-000000000002"},
					},
				}, nil)
			},
			expected: &models.Role{
				ID:          id,
				Name:        "admin",
				Description: "Admin role",
				PermissionIDs: []uuid.UUID{
					uuid.MustParse("10000000-1000-1000-3000-000000000001"),
					uuid.MustParse("10000000-1000-1000-3000-000000000002"),
				},
			},
			error: nil,
		},
		{
			name: "No permission IDs",
			params: &models.Role{
				ID:          id,
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdateRoleRequest{
					Id:            "10000000-1000-1000-1000-000000000001",
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(&proto.UpdateRoleResponse{
					Data: &proto.Role{
						Id:            "10000000-1000-1000-1000-000000000001",
						Name:          "admin",
						Description:   "Admin role",
						PermissionIds: []string{},
					},
				}, nil)
			},
			expected: &models.Role{
				ID:            id,
				Name:          "admin",
				Description:   "Admin role",
				PermissionIDs: []uuid.UUID{},
			},
			error: nil,
		},
		{
			name: "InvalidArgument status code",
			params: &models.Role{
				ID:          id,
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdateRoleRequest{
					Id:            "10000000-1000-1000-1000-000000000001",
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: nil,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "NotFound status code",
			params: &models.Role{
				ID:          id,
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdateRoleRequest{
					Id:            "10000000-1000-1000-1000-000000000001",
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expected: nil,
			error:    errors.ErrRecordNotFound,
		},
		{
			name: "Internal status code",
			params: &models.Role{
				ID:          id,
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdateRoleRequest{
					Id:            "10000000-1000-1000-1000-000000000001",
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: nil,
			error:    errors.ErrFailedToUpdateRecord,
		},
		{
			name: "Error",
			params: &models.Role{
				ID:          id,
				Name:        models.AdminRoleType,
				Description: "Admin role",
			},
			before: func() {
				mockClient.EXPECT().Update(ctx, &proto.UpdateRoleRequest{
					Id:            "10000000-1000-1000-1000-000000000001",
					Name:          "admin",
					Description:   "Admin role",
					PermissionIds: []string{},
				}).Return(nil, assert.AnError)
			},
			expected: nil,
			error:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := service.Update(ctx, tt.params)

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

func Test_Roles_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockClient := rpcs.NewMockRoleServiceClient(ctrl)
	log := logger.NewLogger()
	service := NewRoles(mockClient, log)

	id := uuid.MustParse("10000000-1000-1000-1000-000000000001")

	tests := []struct {
		name     string
		before   func()
		expected bool
		error    error
	}{
		{
			name: "Success",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
				}).Return(&emptypb.Empty{}, nil)
			},
			expected: true,
			error:    nil,
		},
		{
			name: "InvalidArgument status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
				}).Return(nil, status.Error(codes.InvalidArgument, "invalid arguments"))
			},
			expected: false,
			error:    errors.ErrInvalidArguments,
		},
		{
			name: "NotFound status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
				}).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expected: false,
			error:    errors.ErrRecordNotFound,
		},
		{
			name: "Internal status code",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
				}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expected: false,
			error:    errors.ErrFailedToDeleteRecord,
		},
		{
			name: "Error",
			before: func() {
				mockClient.EXPECT().Delete(ctx, &proto.DeleteRoleRequest{
					Id: "10000000-1000-1000-1000-000000000001",
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
