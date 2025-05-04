package rpcs

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"loki-backoffice/internal/app/rpcs/interceptors"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/logger"
)

const (
	CaFile   = "ca.pem"
	CertFile = "client.pem"
	KeyFile  = "client.key"

	KeepaliveTime    = 5 * time.Second
	KeepaliveTimeout = 1 * time.Second
)

type Client interface {
	Connection() *grpc.ClientConn
	Close() error
}

type client struct {
	connection *grpc.ClientConn
	log        *logger.Logger
}

func NewClient(
	cfg *config.Config,
	authInterceptor interceptors.AuthenticationInterceptor,
	traceInterceptor interceptors.TraceInterceptor,
	logInterceptor interceptors.LoggerInterceptor,
	log *logger.Logger,
) (Client, error) {
	tlsConfig, err := setupTLS(cfg, log)
	if err != nil {
		return nil, err
	}

	options := keepalive.ClientParameters{
		Time:                KeepaliveTime,
		Timeout:             KeepaliveTimeout,
		PermitWithoutStream: true,
	}

	connection, err := grpc.NewClient(
		cfg.GrpcAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithKeepaliveParams(options),
		grpc.WithChainUnaryInterceptor(
			authInterceptor.Authenticate(),
			traceInterceptor.Trace(),
			logInterceptor.Log(),
		),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to dial gRPC server")
		return nil, err
	}

	log.Info().Msgf("Connected to gRPC server at %s", cfg.GrpcAddr)
	return &client{
		connection: connection,
		log:        log,
	}, nil
}

func (c *client) Connection() *grpc.ClientConn {
	return c.connection
}

func (c *client) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}

	return nil
}

func setupTLS(cfg *config.Config, log *logger.Logger) (*tls.Config, error) {
	caCert, err := os.ReadFile(filepath.Join(cfg.CertPath, CaFile))
	if err != nil {
		log.Error().Err(err).Msg("Failed to load CA certificate")
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(
		filepath.Join(cfg.CertPath, CertFile),
		filepath.Join(cfg.CertPath, KeyFile),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load client certificate and private key")
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}
