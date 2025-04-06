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

type Users interface {
	List(ctx context.Context, pagination *Pagination) ([]models.User, uint64, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, params *models.User) (*models.User, error)
	Update(ctx context.Context, params *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type users struct {
	client proto.UserServiceClient
	log    *logger.Logger
}

func NewUsers(client proto.UserServiceClient, log *logger.Logger) Users {
	return &users{
		client: client,
		log:    log,
	}
}

//nolint:dupl
func (p *users) List(ctx context.Context, pagination *Pagination) ([]models.User, uint64, error) {
	response, err := p.client.List(ctx, &proto.PaginatedListRequest{
		Limit:  pagination.Page,
		Offset: pagination.PerPage,
	})
	if err != nil {
		p.log.Error().Err(err).Msg("Failed to fetch users")

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

	collection := make([]models.User, 0, len(response.Data))
	for _, item := range response.Data {
		collection = append(collection, models.User{
			ID:             uuid.MustParse(item.Id),
			IdentityNumber: item.IdentityNumber,
			PersonalCode:   item.PersonalCode,
			FirstName:      item.FirstName,
			LastName:       item.LastName,
		})
	}

	return collection, response.Meta.Total, nil
}

//nolint:dupl
func (p *users) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	response, err := p.client.Get(ctx, &proto.GetUserRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to get user")

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

	roleIds := make([]uuid.UUID, 0, len(response.Data.RoleIds))
	if response.Data.RoleIds != nil {
		for _, permissionId := range response.Data.RoleIds {
			itemId, err := uuid.Parse(permissionId)
			if err != nil {
				p.log.Error().Err(err).Str("role_id", permissionId).Msg("Invalid role ID format")
				return nil, errors.ErrInvalidArguments
			}

			roleIds = append(roleIds, itemId)
		}
	}

	scopeIds := make([]uuid.UUID, 0, len(response.Data.ScopeIds))
	if response.Data.ScopeIds != nil {
		for _, scopeId := range response.Data.ScopeIds {
			itemId, err := uuid.Parse(scopeId)
			if err != nil {
				p.log.Error().Err(err).Str("scope_id", scopeId).Msg("Invalid scope ID format")
				return nil, errors.ErrInvalidArguments
			}

			scopeIds = append(scopeIds, itemId)
		}
	}

	return &models.User{
		ID:             uuid.MustParse(response.Data.Id),
		IdentityNumber: response.Data.IdentityNumber,
		PersonalCode:   response.Data.PersonalCode,
		FirstName:      response.Data.FirstName,
		LastName:       response.Data.LastName,
		RoleIDs:        roleIds,
		ScopeIDs:       scopeIds,
	}, nil
}

//nolint:dupl
func (p *users) Create(ctx context.Context, params *models.User) (*models.User, error) {
	paramsRoleIds := []string{}
	if params.RoleIDs != nil {
		paramsRoleIds = make([]string, 0, len(params.RoleIDs))
		for _, permissionId := range params.RoleIDs {
			paramsRoleIds = append(paramsRoleIds, permissionId.String())
		}
	}

	paramsScopeIDs := []string{}
	if params.ScopeIDs != nil {
		paramsScopeIDs = make([]string, 0, len(params.ScopeIDs))
		for _, scopeId := range params.ScopeIDs {
			paramsScopeIDs = append(paramsScopeIDs, scopeId.String())
		}
	}

	response, err := p.client.Create(ctx, &proto.CreateUserRequest{
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		RoleIds:        paramsRoleIds,
		ScopeIds:       paramsScopeIDs,
	})
	if err != nil {
		p.log.Error().Err(err).Str("identity_number", params.IdentityNumber).Msg("Failed to create user")

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

	return &models.User{
		ID:             uuid.MustParse(response.Data.Id),
		IdentityNumber: response.Data.IdentityNumber,
		PersonalCode:   response.Data.PersonalCode,
		FirstName:      response.Data.FirstName,
		LastName:       response.Data.LastName,
	}, nil
}

//nolint:dupl
func (p *users) Update(ctx context.Context, params *models.User) (*models.User, error) {
	paramsRoleIds := []string{}
	if params.RoleIDs != nil {
		paramsRoleIds = make([]string, 0, len(params.RoleIDs))
		for _, permissionId := range params.RoleIDs {
			paramsRoleIds = append(paramsRoleIds, permissionId.String())
		}
	}

	paramsScopeIDs := []string{}
	if params.ScopeIDs != nil {
		paramsScopeIDs = make([]string, 0, len(params.ScopeIDs))
		for _, scopeId := range params.ScopeIDs {
			paramsScopeIDs = append(paramsScopeIDs, scopeId.String())
		}
	}

	response, err := p.client.Update(ctx, &proto.UpdateUserRequest{
		Id:             params.ID.String(),
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		RoleIds:        paramsRoleIds,
		ScopeIds:       paramsScopeIDs,
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", params.ID.String()).Msg("Failed to update user")

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

	return &models.User{
		ID:             uuid.MustParse(response.Data.Id),
		IdentityNumber: response.Data.IdentityNumber,
		PersonalCode:   response.Data.PersonalCode,
		FirstName:      response.Data.FirstName,
		LastName:       response.Data.LastName,
	}, nil
}

//nolint:dupl
func (p *users) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := p.client.Delete(ctx, &proto.DeleteUserRequest{
		Id: id.String(),
	})
	if err != nil {
		p.log.Error().Err(err).Str("id", id.String()).Msg("Failed to delete user")

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
