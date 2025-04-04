package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"loki-backoffice/internal/app/errors"
)

type Payload struct {
	ID          string   `json:"id"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scope       []string `json:"scope,omitempty"`
}

type Jwt interface {
	Decode(token string) (*Payload, error)
}

type jwtService struct {
	parser *jwt.Parser
}

type Claims struct {
	jwt.RegisteredClaims
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scope       []string `json:"scope,omitempty"`
}

func NewJWT() Jwt {
	return &jwtService{
		parser: jwt.NewParser(),
	}
}

func (j *jwtService) Decode(token string) (*Payload, error) {
	result, _, err := j.parser.ParseUnverified(token, &Claims{})
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	if claims, ok := result.Claims.(*Claims); ok {
		return &Payload{
			ID:          claims.ID,
			Roles:       claims.Roles,
			Permissions: claims.Permissions,
			Scope:       claims.Scope,
		}, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
