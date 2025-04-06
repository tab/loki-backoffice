package middlewares

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki-backoffice/pkg/jwt"
)

func Test_CurrentClaimFromContext(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		claim  *jwt.Payload
		exists bool
	}{
		{
			name: "Success",
			ctx: context.WithValue(context.Background(), Claim{}, &jwt.Payload{
				ID:          "test-user",
				Permissions: []string{"read:users"},
			}),
			claim: &jwt.Payload{
				ID:          "test-user",
				Permissions: []string{"read:users"},
			},
			exists: true,
		},
		{
			name:   "Claim does not exist",
			ctx:    context.Background(),
			claim:  nil,
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claim, exists := CurrentClaimFromContext(tt.ctx)
			assert.Equal(t, tt.exists, exists)

			if tt.exists {
				assert.Equal(t, tt.claim.ID, claim.ID)
				assert.Equal(t, tt.claim.Permissions, claim.Permissions)
			} else {
				assert.Nil(t, claim)
			}
		})
	}
}

func Test_CurrentTokenFromContext(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		token  string
		exists bool
	}{
		{
			name:   "Success",
			ctx:    context.WithValue(context.Background(), Token{}, "test-token"),
			token:  "test-token",
			exists: true,
		},
		{
			name:   "Token does not exist",
			ctx:    context.Background(),
			token:  "",
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, exists := CurrentTokenFromContext(tt.ctx)
			assert.Equal(t, tt.exists, exists)

			if tt.exists {
				assert.Equal(t, tt.token, token)
			} else {
				assert.Empty(t, token)
			}
		})
	}
}

func Test_CurrentTraceIdFromContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		traceId string
		exists  bool
	}{
		{
			name:    "Success",
			ctx:     context.WithValue(context.Background(), TraceId{}, "9809b3e0-484b-438c-80b2-73cb9af51cd4"),
			traceId: "9809b3e0-484b-438c-80b2-73cb9af51cd4",
			exists:  true,
		},
		{
			name:    "TraceId does not exist",
			ctx:     context.Background(),
			traceId: "",
			exists:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traceId, exists := CurrentTraceIdFromContext(tt.ctx)
			assert.Equal(t, tt.exists, exists)

			if tt.exists {
				assert.Equal(t, tt.traceId, traceId)
			} else {
				assert.Empty(t, traceId)
			}
		})
	}
}
