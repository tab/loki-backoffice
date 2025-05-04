package config

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki-backoffice/pkg/spec"
)

func TestMain(m *testing.M) {
	if err := spec.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

func Test_LoadConfig(t *testing.T) {
	certDir := spec.GenerateCertificates(t)

	tests := []struct {
		name     string
		args     []string
		env      map[string]string
		expected *Config
	}{
		{
			name: "Success",
			args: []string{},
			env:  map[string]string{},
			expected: &Config{
				AppEnv:      "test",
				AppAddr:     "0.0.0.0:8081",
				GrpcAddr:    "0.0.0.0:50051",
				ClientURL:   "http://localhost:3001",
				CertPath:    certDir,
				DatabaseDSN: "postgres://postgres:postgres@localhost:5432/loki-backoffice-test?sslmode=disable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
			}
			os.Setenv("CERT_PATH", certDir)

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			result := LoadConfig()

			assert.Equal(t, tt.expected.AppEnv, result.AppEnv)
			assert.Equal(t, tt.expected.AppAddr, result.AppAddr)
			assert.Equal(t, tt.expected.GrpcAddr, result.GrpcAddr)
			assert.Equal(t, tt.expected.ClientURL, result.ClientURL)
			assert.Equal(t, tt.expected.CertPath, result.CertPath)
			assert.Equal(t, tt.expected.DatabaseDSN, result.DatabaseDSN)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
				os.Unsetenv("CERT_PATH")
			})
		})
	}
}
