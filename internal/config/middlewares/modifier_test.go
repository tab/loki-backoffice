package middlewares

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki-backoffice/pkg/jwt"
)

func Test_NewContextModifier(t *testing.T) {
	ctx := context.Background()
	ctxModifier := NewContextModifier(ctx)

	assert.NotNil(t, ctxModifier)
	assert.Equal(t, ctx, ctxModifier.Context())
}

func Test_Modifier_WithClaim(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		claims *jwt.Payload
	}{
		{
			name: "Success",
			claims: &jwt.Payload{
				ID:          "test-user",
				Permissions: []string{"read:self", "write:self"},
			},
		},
		{
			name: "Empty permissions",
			claims: &jwt.Payload{
				ID:          "test-user",
				Permissions: []string{},
			},
		},
		{
			name:   "Nil claims",
			claims: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxModifier := NewContextModifier(ctx).WithClaim(tt.claims)

			claim, ok := ctxModifier.Context().Value(Claim{}).(*jwt.Payload)

			if tt.claims == nil {
				assert.True(t, ok)
				assert.Nil(t, claim)
			} else {
				assert.True(t, ok)
				assert.Equal(t, tt.claims.ID, claim.ID)
				assert.Equal(t, tt.claims.Permissions, claim.Permissions)
			}
		})
	}
}

func Test_Modifier_WithToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "Valid token",
			token: "valid-token",
		},
		{
			name:  "Empty token",
			token: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxModifier := NewContextModifier(ctx).WithToken(tt.token)

			token, ok := ctxModifier.Context().Value(Token{}).(string)
			assert.True(t, ok)
			assert.Equal(t, tt.token, token)
		})
	}
}

func Test_Modifier_WithTraceId(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		traceId string
	}{
		{
			name:    "Valid trace ID",
			traceId: "valid-trace-id",
		},
		{
			name:    "Empty trace ID",
			traceId: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxModifier := NewContextModifier(ctx).WithTraceId(tt.traceId)

			traceId, ok := ctxModifier.Context().Value(TraceId{}).(string)
			assert.True(t, ok)
			assert.Equal(t, tt.traceId, traceId)
		})
	}
}
