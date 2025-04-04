package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki-backoffice/internal/app/errors"
)

func Test_JWT_Decode(t *testing.T) {
	service := NewJWT()

	tests := []struct {
		name     string
		token    string
		expected *Payload
		error    error
	}{
		{
			name:  "Success",
			token: generateToken("PNOEE-30303039914", []string{"read:all"}, []string{"admin"}, []string{"scope:all"}),
			expected: &Payload{
				ID:          "PNOEE-30303039914",
				Roles:       []string{"admin"},
				Permissions: []string{"read:all"},
				Scope:       []string{"scope:all"},
			},
		},
		{
			name:  "Success with empty payload",
			token: generateToken("PNOEE-30303039914", nil, nil, nil),
			expected: &Payload{
				ID: "PNOEE-30303039914",
			},
		},
		{
			name:     "Invalid JWT format",
			token:    "invalid.jwt.token",
			expected: nil,
			error:    errors.ErrInvalidToken,
		},
		{
			name:     "Empty token",
			token:    "",
			expected: nil,
			error:    errors.ErrInvalidToken,
		},
		{
			name:     "Malformed token",
			token:    "only.one.part",
			expected: nil,
			error:    errors.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Decode(tt.token)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_JWT_Decode_Mocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockJwt(ctrl)

	tests := []struct {
		name     string
		token    string
		setup    func()
		expected *Payload
		error    error
	}{
		{
			name:  "Success",
			token: "valid.token.signature",
			setup: func() {
				mockService.EXPECT().Decode("valid.token.signature").Return(&Payload{
					ID:          "user-123",
					Roles:       []string{"admin"},
					Permissions: []string{"read:all"},
					Scope:       []string{"api"},
				}, nil)
			},
			expected: &Payload{
				ID:          "user-123",
				Roles:       []string{"admin"},
				Permissions: []string{"read:all"},
				Scope:       []string{"api"},
			},
		},
		{
			name:  "Error",
			token: "invalid.token",
			setup: func() {
				mockService.EXPECT().Decode("invalid.token").Return(nil, errors.ErrInvalidToken)
			},
			expected: nil,
			error:    errors.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := mockService.Decode(tt.token)

			if tt.error != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func generateToken(id string, permissions, roles, scope []string) string {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Roles:       roles,
		Permissions: permissions,
		Scope:       scope,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("test-secret-key"))
	return signedToken
}
