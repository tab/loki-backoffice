package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	AppAddr    = "0.0.0.0:8081"
	ClientURL  = "http://localhost:3001"
	GrpcAddr   = "0.0.0.0:50051"
	DebugLevel = "debug"
)

type Config struct {
	AppEnv       string
	AppName      string
	AppAddr      string
	GrpcAddr     string
	ClientURL    string
	CertPath     string
	DatabaseDSN  string
	TelemetryURI string
	LogLevel     string
}

func LoadConfig() *Config {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	envFiles := []string{
		".env",
		fmt.Sprintf(".env.%s", env),
		fmt.Sprintf(".env.%s.local", env),
	}
	for _, file := range envFiles {
		_ = godotenv.Overload(file)
	}

	flagAppAddr := flag.String("b", AppAddr, "server address")
	flagGrpcAddr := flag.String("g", GrpcAddr, "gRPC server address")
	flagClientURL := flag.String("c", ClientURL, "client address")
	flagCertPath := flag.String("p", "", "certificate path")
	flagDatabaseDSN := flag.String("d", "", "database DSN")
	flagTelemetryURI := flag.String("t", "", "OpenTelemetry collector URI")
	flag.Parse()

	return &Config{
		AppEnv:    env,
		AppName:   getEnvString("APP_NAME"),
		AppAddr:   getFlagOrEnvString(*flagAppAddr, "APP_ADDRESS", AppAddr),
		GrpcAddr:  getFlagOrEnvString(*flagGrpcAddr, "GRPC_ADDRESS", GrpcAddr),
		ClientURL: getFlagOrEnvString(*flagClientURL, "CLIENT_URL", ClientURL),

		CertPath: getFlagOrEnvString(*flagCertPath, "CERT_PATH", ""),

		DatabaseDSN:  getFlagOrEnvString(*flagDatabaseDSN, "DATABASE_DSN", ""),
		TelemetryURI: getFlagOrEnvString(*flagTelemetryURI, "TELEMETRY_URI", ""),

		LogLevel: getEnvString("LOG_LEVEL"),
	}
}

func getFlagOrEnvString(flagValue, envVar, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}

	if envValue, ok := os.LookupEnv(envVar); ok && envValue != "" {
		return envValue
	}

	return defaultValue
}

func getEnvString(envVar string) string {
	if envValue, ok := os.LookupEnv(envVar); ok && envValue != "" {
		return envValue
	}

	return ""
}
