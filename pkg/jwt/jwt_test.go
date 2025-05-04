package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/config"
)

func Test_NewJWT(t *testing.T) {
	tempDir := generateTestKeys(t)

	cfg := &config.Config{
		CertPath: tempDir,
	}
	service, err := NewJWT(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func Test_JWT_Decode(t *testing.T) {
	tempDir := generateTestKeys(t)

	cfg := &config.Config{
		CertPath: tempDir,
	}
	service, err := NewJWT(cfg)
	assert.NoError(t, err)

	privateKey, err := loadPrivateKey(cfg)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		token    string
		expected *Payload
		error    error
	}{
		{
			name:  "Success",
			token: generateToken("PNOEE-30303039914", []string{"read:all"}, []string{"admin"}, []string{"scope:all"}, privateKey),
			expected: &Payload{
				ID:          "PNOEE-30303039914",
				Roles:       []string{"admin"},
				Permissions: []string{"read:all"},
				Scope:       []string{"scope:all"},
			},
		},
		{
			name:  "Success with empty payload",
			token: generateToken("PNOEE-30303039914", nil, nil, nil, privateKey),
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

func generateToken(id string, permissions, roles, scope []string, privateKey *rsa.PrivateKey) string {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Roles:       roles,
		Permissions: permissions,
		Scope:       scope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, _ := token.SignedString(privateKey)
	return signedToken
}

func generateTestKeys(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "jwt-test-*")
	require.NoError(t, err)

	t.Cleanup(func() { os.RemoveAll(tempDir) })

	jwtDir := filepath.Join(tempDir, Dir)
	err = os.MkdirAll(jwtDir, 0755)
	require.NoError(t, err)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err = os.WriteFile(filepath.Join(jwtDir, PrivateKeyFile), privateKeyPEM, 0600)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(jwtDir, PublicKeyFile), publicKeyPEM, 0644)
	require.NoError(t, err)

	return tempDir
}
