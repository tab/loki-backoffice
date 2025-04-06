package rpcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki-backoffice/internal/app/rpcs/interceptors"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/logger"
	"loki-backoffice/pkg/spec"
)

func Test_NewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthInterceptor := interceptors.NewMockAuthenticationInterceptor(ctrl)
	mockTraceInterceptor := interceptors.NewMockTraceInterceptor(ctrl)
	mockLogInterceptor := interceptors.NewMockLoggerInterceptor(ctrl)

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
		GrpcAddr: "localhost:9999",
		CertPath: "/non-existent-path",
	}
	log := logger.NewLogger(cfg)

	c, err := NewClient(
		cfg,
		mockAuthInterceptor,
		mockTraceInterceptor,
		mockLogInterceptor,
		log)
	assert.Error(t, err)
	assert.Nil(t, c)
}

func Test_Client_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	c := &client{
		connection: nil,
		log:        log,
	}

	err := c.Close()
	assert.NoError(t, err)
}

func Test_Client_SetupTLS(t *testing.T) {
	certDir := spec.GenerateCertificates(t)

	tests := []struct {
		name     string
		certDir  string
		expected error
	}{
		{
			name:     "Success",
			certDir:  certDir,
			expected: nil,
		},
		{
			name:     "Invalid cert path",
			certDir:  "/non-existent-path",
			expected: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AppEnv:   "test",
				AppAddr:  "localhost:8080",
				GrpcAddr: "localhost:50051",
				LogLevel: "info",
				CertPath: tt.certDir,
			}
			log := logger.NewLogger(cfg)

			tlsConfig, err := setupTLS(cfg, log)

			if tt.expected != nil {
				assert.Error(t, err)
				assert.Nil(t, tlsConfig)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tlsConfig)
			}
		})
	}
}
