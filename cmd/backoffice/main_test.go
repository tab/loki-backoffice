package main

import (
	"flag"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"loki-backoffice/internal/config"
	"loki-backoffice/pkg/spec"
)

func Test_Main(t *testing.T) {
	certDir := spec.GenerateCertificates(t)

	os.Setenv("CERT_PATH", certDir)

	if _, err := os.Stat("certs"); os.IsNotExist(err) {
		err = os.Mkdir("certs", 0755)
		require.NoError(t, err)
		t.Cleanup(func() {
			os.RemoveAll("certs")
		})
	}

	certFiles := []string{
		spec.CaFile,
		spec.ClientCertFile,
		spec.ClientKeyFile,
		spec.ServerCertFile,
		spec.ServerKeyFile,
	}

	for _, file := range certFiles {
		srcPath := filepath.Join(certDir, file)
		destPath := filepath.Join("certs", file)

		data, err := os.ReadFile(srcPath)
		require.NoError(t, err)

		err = os.WriteFile(destPath, data, 0644)
		require.NoError(t, err)
	}

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8081",
		GrpcAddr: "localhost:50051",
		CertPath: certDir,
	}
	const BaseURL = "http://localhost:8081"

	tests := []struct {
		name   string
		signal os.Signal
	}{
		{
			name:   "SIGTERM",
			signal: syscall.SIGTERM,
		},
		{
			name:   "SIGINT",
			signal: syscall.SIGINT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = []string{oldArgs[0], cfg.AppAddr}

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			done := make(chan struct{})

			go func() {
				main()
				close(done)
			}()

			spec.WaitForServerStart(t, BaseURL+"/live")

			p, err := os.FindProcess(os.Getpid())
			require.NoError(t, err)
			require.NotNil(t, p)

			require.NoError(t, p.Signal(tt.signal))

			select {
			case <-done:
				// main() exited successfully
			case <-time.After(spec.ServerStopTimeout):
				t.Fatal("timeout: main() did not exit after signal")
			}
		})
	}
}
