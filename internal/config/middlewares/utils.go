package middlewares

import (
	"context"

	"loki-backoffice/pkg/jwt"
)

const (
	Authorization = "Authorization"
	bearerScheme  = "Bearer "
)

func CurrentClaimFromContext(ctx context.Context) (*jwt.Payload, bool) {
	c, ok := ctx.Value(Claim{}).(*jwt.Payload)
	return c, ok
}

func CurrentTokenFromContext(ctx context.Context) (string, bool) {
	t, ok := ctx.Value(Token{}).(string)
	return t, ok
}
