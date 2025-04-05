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

type Scopes interface {
	List(ctx context.Context, pagination *Pagination) ([]models.Scope, uint64, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Scope, error)
	Create(ctx context.Context, params *models.Scope) (*models.Scope, error)
	Update(ctx context.Context, params *models.Scope) (*models.Scope, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type scopes struct {
	client proto.ScopeServiceClient
	log    *logger.Logger
}

func NewScopes(client proto.ScopeServiceClient, log *logger.Logger) Scopes {
	return &scopes{
		client: client,
		log:    log,
	}
}

//nolint:dupl
func (p *scopes) List(ctx context.Context, pagination *Pagination) ([]models.Scope, uint64, error) {
	response, err := p.client.List(ctx, &proto.PaginatedListRequest{
		Limit:  pagination.Page,
		Offset: pagination.PerPage,
	})
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch scopes")

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

	collection := make([]models.Scope, 0, len(response.Data))
	for _, item := range response.Data {
		collection = append(collection, models.Scope{
			ID:          uuid.MustParse(item.Id),
			Name:        item.Name,
			Description: item.Description,
		})
	}

	return collection, response.Meta.Total, nil
}

//nolint:dupl
func (p *scopes) FindById(ctx context.Context, id uuid.UUID) (*models.Scope, error) {
	response, err := p.client.Get(ctx, &proto.GetScopeRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to get scope")

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

	return &models.Scope{
		ID:          uuid.MustParse(response.Data.Id),
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}, nil
}

//nolint:dupl
func (p *scopes) Create(ctx context.Context, params *models.Scope) (*models.Scope, error) {
	response, err := p.client.Create(ctx, &proto.CreateScopeRequest{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		p.log.Error().Err(err).Str("name", params.Name).Msg("Failed to create scope")

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

	return &models.Scope{
		ID:          uuid.MustParse(response.Data.Id),
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}, nil
}

//nolint:dupl
func (p *scopes) Update(ctx context.Context, params *models.Scope) (*models.Scope, error) {
	response, err := p.client.Update(ctx, &proto.UpdateScopeRequest{
		Id:          params.ID.String(),
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", params.ID.String()).Msg("Failed to update scope")

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

	return &models.Scope{
		ID:          uuid.MustParse(response.Data.Id),
		Name:        response.Data.Name,
		Description: response.Data.Description,
	}, nil
}

//nolint:dupl
func (p *scopes) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := p.client.Delete(ctx, &proto.DeleteScopeRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to delete scope")

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
