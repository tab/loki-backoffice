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

type Tokens interface {
	List(ctx context.Context, pagination *Pagination) ([]models.Token, uint64, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type tokens struct {
	client proto.TokenServiceClient
	log    *logger.Logger
}

func NewTokens(client proto.TokenServiceClient, log *logger.Logger) Tokens {
	return &tokens{
		client: client,
		log:    log,
	}
}

//nolint:dupl
func (p *tokens) List(ctx context.Context, pagination *Pagination) ([]models.Token, uint64, error) {
	response, err := p.client.List(ctx, &proto.PaginatedListRequest{
		Limit:  pagination.Page,
		Offset: pagination.PerPage,
	})
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch tokens")

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

	collection := make([]models.Token, 0, len(response.Data))
	for _, item := range response.Data {
		collection = append(collection, models.Token{
			ID:        uuid.MustParse(item.Id),
			UserId:    uuid.MustParse(item.UserId),
			Type:      item.Type,
			Value:     item.Value,
			ExpiresAt: item.ExpiresAt.AsTime(),
		})
	}

	return collection, response.Meta.Total, nil
}

//nolint:dupl
func (p *tokens) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := p.client.Delete(ctx, &proto.DeleteTokenRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to delete token")

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
