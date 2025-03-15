package main

import (
    "flag"
    "os"
    "syscall"
    "testing"
    "time"

    "github.com/stretchr/testify/require"

    "loki-backoffice/internal/config"
    "loki-backoffice/pkg/spec"
)

func Test_Main(t *testing.T) {
    cfg := &config.Config{
        AppAddr: "localhost:8081",
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
