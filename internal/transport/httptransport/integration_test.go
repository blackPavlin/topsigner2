package httptransport_test

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/database"
	"github.com/bboykiv/topsigner/internal/database/repository"
	"github.com/bboykiv/topsigner/internal/s3"
	"github.com/bboykiv/topsigner/internal/s3/storage"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/font"
	"github.com/bboykiv/topsigner/internal/service/image"
	"github.com/bboykiv/topsigner/internal/transport/httptransport"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cfg := &config.Config{
		Auth: config.AuthConfig{
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 30 * 24 * time.Hour,
			SigningKey:      "test-signing-key",
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
	}

	pgContainer, err := postgres.Run(ctx,
		"postgres:18.4-alpine",
		postgres.WithDatabase(cfg.Postgres.Database),
		postgres.WithUsername(cfg.Postgres.User),
		postgres.WithPassword(cfg.Postgres.Password),
		testcontainers.WithAdditionalWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		log.Printf("start postgres container: %v", err)

		return 1
	}
	defer func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			log.Printf("terminate postgres container: %v", err)
		}
	}()

	cfg.Postgres.Host, err = pgContainer.Host(ctx)
	if err != nil {
		log.Printf("get postgres host: %v", err)

		return 1
	}

	port, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Printf("get postgres port: %v", err)

		return 1
	}

	cfg.Postgres.Port = int(port.Num())

	pool, err := database.NewPool(cfg)
	if err != nil {
		log.Printf("create database pool: %v", err)

		return 1
	}
	defer pool.Close()

	if err := database.MakeMigrations(pool, cfg); err != nil {
		log.Printf("run migrations: %v", err)

		return 1
	}

	minioContainer, err := minio.Run(ctx,
		"minio/minio:latest",
		minio.WithUsername(cfg.S3.AccessKey),
		minio.WithPassword(cfg.S3.SecretKey),
	)
	if err != nil {
		log.Printf("start minio container: %v", err)

		return 1
	}
	defer func() {
		if err := minioContainer.Terminate(ctx); err != nil {
			log.Printf("terminate minio container: %v", err)
		}
	}()

	cfg.S3.Endpoint, err = minioContainer.ConnectionString(ctx)
	if err != nil {
		log.Printf("get minio connection string: %v", err)

		return 1
	}

	minioClient, err := s3.NewClient(cfg)
	if err != nil {
		log.Printf("create minio client: %v", err)

		return 1
	}

	if err := s3.CreateBuckets(ctx, minioClient, cfg); err != nil {
		log.Printf("create buckets: %v", err)

		return 1
	}

	var (
		logger            = zap.NewNop()
		userRepository    = repository.NewUserRepository(pool)
		sessionRepository = repository.NewSessionRepository(pool)
		imageRepository   = repository.NewImageRepository(pool)
		fontRepository    = repository.NewFontRepository(pool)
		imageStorage      = storage.NewImageStorage(cfg, minioClient)
		authService       = auth.New(logger, cfg, userRepository, sessionRepository)
		imageService      = image.New(logger, imageRepository, imageStorage)
		fontService       = font.New(logger, fontRepository)
	)

	server = httptest.NewServer(
		httptransport.NewHandler(
			logger,
			cfg,
			authService,
			imageService,
			fontService,
		),
	)
	defer server.Close()

	return m.Run()
}
