package httptransport_test

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/application"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
)

const (
	startTimeout = 2 * time.Minute

	defaultUserEmail    = "admin@topsigner.test"
	defaultUserPassword = "password1234"
)

type IntegrationSuite struct {
	suite.Suite

	application *fxtest.App
	server      *httptest.Server
	client      *httpserver.ClientWithResponses
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(s.T().Context(), startTimeout)
	defer cancel()

	config := s.newConfig()

	s.startPostgres(ctx, config)
	s.startMinio(ctx, config)
	s.startRedis(ctx, config)

	// todo: описать тестовый сервер vk
	vkServer := httptest.NewServer(http.NewServeMux())
	s.T().Cleanup(vkServer.Close)

	// todo: описать тестовый сервер vkid
	vkidServer := httptest.NewServer(http.NewServeMux())
	s.T().Cleanup(vkidServer.Close)

	config.VK.BaseURL = vkServer.URL
	config.VK.OAuthBaseURL = vkServer.URL
	config.VKID.BaseURL = vkidServer.URL

	var handler http.Handler

	s.application = fxtest.New(
		s.T(),
		application.New(),
		fx.Replace(
			config,
			zaptest.NewLogger(s.T(), zaptest.Level(zap.WarnLevel)),
		),
		fx.Populate(&handler),
	)
	s.application.RequireStart()

	s.server = httptest.NewServer(handler)
	s.T().Cleanup(s.server.Close)

	client, err := httpserver.NewClientWithResponses(s.server.URL)
	s.Require().NoError(err)

	s.client = client
}

func (s *IntegrationSuite) newConfig() *config.Config {
	return &config.Config{
		Auth: config.AuthConfig{
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 30 * 24 * time.Hour,
			SigningKey:      "test-signing-key",
			EncryptionKey:   base64.StdEncoding.EncodeToString(make([]byte, 32)),
		},
		User: config.UserConfig{
			Default: config.DefaultUserConfig{
				Email:    defaultUserEmail,
				Password: defaultUserPassword,
				Role:     model.RoleAdmin,
			},
		},
		Cors: config.CorsConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "DELETE"},
			AllowedHeaders: []string{"*"},
		},
		Postgres: config.PostgresConfig{
			User:            "test",
			Password:        "test",
			Database:        "topsigner_test",
			Schema:          "public",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnTimeout:     5 * time.Second,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: 15 * time.Minute,
		},
		S3: config.S3Config{
			Region:      "us-east-1",
			AccessKey:   "minioadmin",
			SecretKey:   "minioadmin",
			ImageBucket: "images-test",
			FontBucket:  "fonts-test",
			Secure:      false,
		},
		Redis: config.RedisConfig{
			DialTimeout: 5 * time.Second,
		},
	}
}

func (s *IntegrationSuite) startPostgres(ctx context.Context, config *config.Config) {
	container, err := postgres.Run(ctx,
		"postgres:18.4-alpine",
		postgres.WithDatabase(config.Postgres.Database),
		postgres.WithUsername(config.Postgres.User),
		postgres.WithPassword(config.Postgres.Password),
		postgres.BasicWaitStrategies(),
	)
	s.Require().NoError(err, "start postgres container")

	testcontainers.CleanupContainer(s.T(), container)

	host, err := container.Host(ctx)
	s.Require().NoError(err)

	port, err := container.MappedPort(ctx, "5432/tcp")
	s.Require().NoError(err)

	config.Postgres.Host = host
	config.Postgres.Port = int(port.Num())
}

func (s *IntegrationSuite) startMinio(ctx context.Context, config *config.Config) {
	container, err := minio.Run(ctx,
		"pgsty/silo:RELEASE.2026-09-16T00-00-00Z",
		minio.WithUsername(config.S3.AccessKey),
		minio.WithPassword(config.S3.SecretKey),
	)
	s.Require().NoError(err, "start minio container")

	testcontainers.CleanupContainer(s.T(), container)

	endpoint, err := container.ConnectionString(ctx)
	s.Require().NoError(err)

	config.S3.Endpoint = endpoint
}

func (s *IntegrationSuite) startRedis(ctx context.Context, config *config.Config) {
	container, err := redis.Run(ctx, "redis:8.10-alpine")
	s.Require().NoError(err, "start redis container")

	testcontainers.CleanupContainer(s.T(), container)

	host, err := container.Host(ctx)
	s.Require().NoError(err)

	port, err := container.MappedPort(ctx, "6379/tcp")
	s.Require().NoError(err)

	config.Redis.Addr = net.JoinHostPort(host, port.Port())
}
