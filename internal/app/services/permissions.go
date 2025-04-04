package services

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/app/models"
	proto "loki-backoffice/internal/app/rpcs/proto/sso/v1"
	"loki-backoffice/pkg/logger"
)

type Permissions interface {
	List(ctx context.Context, pagination *Pagination) ([]models.Permission, uint64, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	Create(ctx context.Context, params *models.Permission) (*models.Permission, error)
	Update(ctx context.Context, params *models.Permission) (*models.Permission, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type permissions struct {
	client proto.PermissionServiceClient
	log    *logger.Logger
}

func NewPermissions(client proto.PermissionServiceClient, log *logger.Logger) Permissions {
	return &permissions{
		client: client,
		log:    log,
	}
}

func (p *permissions) List(ctx context.Context, pagination *Pagination) ([]models.Permission, uint64, error) {
	response, err := p.client.List(ctx, &proto.PaginatedListRequest{
		Limit:  pagination.Page,
		Offset: pagination.PerPage,
	})
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch permissions")

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

	collection := make([]models.Permission, 0, len(response.Data))
	for _, item := range response.Data {
		collection = append(collection, models.Permission{
			ID:          uuid.MustParse(item.Id),
			Name:        item.Name,
			Description: item.Description,
		})
	}

	return collection, response.Meta.Total, nil
}

func (p *permissions) FindById(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	response, err := p.client.Get(ctx, &proto.GetPermissionRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to get permission")

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

	return &models.Permission{
		ID:          uuid.MustParse(response.Data.Id),
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}, nil
}

func (p *permissions) Create(ctx context.Context, params *models.Permission) (*models.Permission, error) {
	response, err := p.client.Create(ctx, &proto.CreatePermissionRequest{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		p.log.Error().Err(err).Str("name", params.Name).Msg("Failed to create permission")

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

	return &models.Permission{
		ID:          uuid.MustParse(response.Data.Id),
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}, nil
}

func (p *permissions) Update(ctx context.Context, params *models.Permission) (*models.Permission, error) {
	response, err := p.client.Update(ctx, &proto.UpdatePermissionRequest{
		Id:          params.ID.String(),
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", params.ID.String()).Msg("Failed to update permission")

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

	return &models.Permission{
		ID:          uuid.MustParse(response.Data.Id),
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}, nil
}

func (p *permissions) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := p.client.Delete(ctx, &proto.DeletePermissionRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to delete permission")

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
