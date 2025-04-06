package services

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/app/models"
	proto "loki-backoffice/internal/app/rpcs/proto/sso/v1"
	"loki-backoffice/internal/config/logger"
)

type Roles interface {
	List(ctx context.Context, pagination *Pagination) ([]models.Role, uint64, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Role, error)
	Create(ctx context.Context, params *models.Role) (*models.Role, error)
	Update(ctx context.Context, params *models.Role) (*models.Role, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type roles struct {
	client proto.RoleServiceClient
	log    *logger.Logger
}

func NewRoles(client proto.RoleServiceClient, log *logger.Logger) Roles {
	return &roles{
		client: client,
		log:    log,
	}
}

//nolint:dupl
func (p *roles) List(ctx context.Context, pagination *Pagination) ([]models.Role, uint64, error) {
	response, err := p.client.List(ctx, &proto.PaginatedListRequest{
		Limit:  pagination.Page,
		Offset: pagination.PerPage,
	})
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch roles")

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, 0, errors.ErrInvalidArguments
			case codes.Unavailable:
				return nil, 0, errors.ErrFailedToFetchResults
			case codes.Internal:
				return nil, 0, errors.ErrFailedToFetchResults
			}
		}

		return nil, 0, err
	}

	collection := make([]models.Role, 0, len(response.Data))
	for _, item := range response.Data {
		collection = append(collection, models.Role{
			ID:          uuid.MustParse(item.Id),
			Name:        item.Name,
			Description: item.Description,
		})
	}

	return collection, response.Meta.Total, nil
}

//nolint:dupl
func (p *roles) FindById(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	response, err := p.client.Get(ctx, &proto.GetRoleRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to get role")

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, errors.ErrInvalidArguments
			case codes.NotFound:
				return nil, errors.ErrRecordNotFound
			case codes.Internal:
				return nil, errors.ErrFailedToFetchResults
			}
		}

		return nil, err
	}

	permissionIds := make([]uuid.UUID, 0, len(response.Data.PermissionIds))
	if response.Data.PermissionIds != nil {
		for _, permissionId := range response.Data.PermissionIds {
			itemId, err := uuid.Parse(permissionId)
			if err != nil {
				p.log.Error().Err(err).Str("permission_id", permissionId).Msg("Invalid permission ID format")
				return nil, errors.ErrInvalidArguments
			}

			permissionIds = append(permissionIds, itemId)
		}
	}

	return &models.Role{
		ID:            uuid.MustParse(response.Data.Id),
		Name:          response.Data.Name,
		Description:   response.Data.Description,
		PermissionIDs: permissionIds,
	}, nil
}

//nolint:dupl
func (p *roles) Create(ctx context.Context, params *models.Role) (*models.Role, error) {
	paramsPermissionIds := []string{}
	if params.PermissionIDs != nil {
		paramsPermissionIds = make([]string, 0, len(params.PermissionIDs))
		for _, permissionId := range params.PermissionIDs {
			paramsPermissionIds = append(paramsPermissionIds, permissionId.String())
		}
	}

	response, err := p.client.Create(ctx, &proto.CreateRoleRequest{
		Name:          params.Name,
		Description:   params.Description,
		PermissionIds: paramsPermissionIds,
	})
	if err != nil {
		p.log.Error().Err(err).Str("name", params.Name).Msg("Failed to create role")

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, errors.ErrInvalidArguments
			case codes.Internal:
				return nil, errors.ErrFailedToCreateRecord
			}
		}

		return nil, err
	}

	permissionIds := make([]uuid.UUID, 0, len(response.Data.PermissionIds))
	if response.Data.PermissionIds != nil {
		for _, permissionId := range response.Data.PermissionIds {
			itemId, err := uuid.Parse(permissionId)
			if err != nil {
				p.log.Error().Err(err).Str("permission_id", permissionId).Msg("Invalid permission ID format")
				return nil, errors.ErrInvalidArguments
			}

			permissionIds = append(permissionIds, itemId)
		}
	}

	return &models.Role{
		ID:            uuid.MustParse(response.Data.Id),
		Name:          response.Data.Name,
		Description:   response.Data.Description,
		PermissionIDs: permissionIds,
	}, nil
}

//nolint:dupl
func (p *roles) Update(ctx context.Context, params *models.Role) (*models.Role, error) {
	paramsPermissionIds := []string{}
	if params.PermissionIDs != nil {
		paramsPermissionIds = make([]string, 0, len(params.PermissionIDs))
		for _, permissionId := range params.PermissionIDs {
			paramsPermissionIds = append(paramsPermissionIds, permissionId.String())
		}
	}

	response, err := p.client.Update(ctx, &proto.UpdateRoleRequest{
		Id:            params.ID.String(),
		Name:          params.Name,
		Description:   params.Description,
		PermissionIds: paramsPermissionIds,
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", params.ID.String()).Msg("Failed to update role")

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, errors.ErrInvalidArguments
			case codes.NotFound:
				return nil, errors.ErrRecordNotFound
			case codes.Internal:
				return nil, errors.ErrFailedToUpdateRecord
			}
		}

		return nil, err
	}

	permissionIds := make([]uuid.UUID, 0, len(response.Data.PermissionIds))
	if response.Data.PermissionIds != nil {
		for _, permissionId := range response.Data.PermissionIds {
			itemId, err := uuid.Parse(permissionId)
			if err != nil {
				p.log.Error().Err(err).Str("permission_id", permissionId).Msg("Invalid permission ID format")
				return nil, errors.ErrInvalidArguments
			}

			permissionIds = append(permissionIds, itemId)
		}
	}

	return &models.Role{
		ID:            uuid.MustParse(response.Data.Id),
		Name:          response.Data.Name,
		Description:   response.Data.Description,
		PermissionIDs: permissionIds,
	}, nil
}

//nolint:dupl
func (p *roles) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := p.client.Delete(ctx, &proto.DeleteRoleRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to delete role")

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return false, errors.ErrInvalidArguments
			case codes.NotFound:
				return false, errors.ErrRecordNotFound
			case codes.Internal:
				return false, errors.ErrFailedToDeleteRecord
			}
		}

		return false, err
	}

	return true, nil
}
