package e2e

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/router"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	DB          *pgxpool.Pool
	Redis       *redis.Client
	Server      *httptest.Server
	Client      *TestClient
	CommunityID string

	pgContainer    *postgres.PostgresContainer
	redisContainer *rediscontainer.RedisContainer

	Cleanup func()
}

func SetupTestEnv(t *testing.T) *TestEnv {
	return SetupTestEnvWithSuffix(t, generateRandomSuffix())
}

func SetupTestEnvWithSuffix(t *testing.T, suffix string) *TestEnv {
	ctx := context.Background()

	dbName := fmt.Sprintf("projek_test_%s", suffix)
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get PostgreSQL connection string")

	db, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err, "Failed to connect to PostgreSQL")

	err = db.Ping(ctx)
	require.NoError(t, err, "Failed to ping PostgreSQL")

	err = RunMigrations(ctx, db)
	require.NoError(t, err, "Failed to run migrations")

	t.Logf("PostgreSQL ready: %s", dbName)

	redisContainer, err := rediscontainer.Run(ctx, "redis:7-alpine")
	require.NoError(t, err, "Failed to start Redis container")

	redisHost, err := redisContainer.Host(ctx)
	require.NoError(t, err, "Failed to get Redis host")

	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err, "Failed to get Redis port")

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort.Port()),
	})

	err = redisClient.Ping(ctx).Err()
	require.NoError(t, err, "Failed to ping Redis")

	t.Logf("Redis ready: %s:%s", redisHost, redisPort.Port())

	cfg := &common.Config{
		Server: common.ServerConfig{
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: common.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			Database: dbName,
			PoolSize: 5,
		},
		Redis: common.RedisConfig{
			Host: redisHost,
			Port: 6379,
		},
		Auth: common.AuthConfig{
			SessionTTL:     24 * time.Hour,
			PasswordMinLen: 8,
		},
		Email: common.EmailConfig{
			SMTPHost:  "localhost",
			SMTPPort:  587,
			FromEmail: "test@localhost",
			FromName:  "Test",
		},
	}

	logger := zerolog.Nop()
	r := router.NewRouter(cfg, logger, db, redisClient)

	server := httptest.NewServer(r)
	t.Logf("Test server ready: %s", server.URL)

	cleanup := func() {
		t.Log("Cleaning up test environment...")
		server.Close()
		db.Close()
		redisClient.Close()

		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Warning: failed to terminate PostgreSQL: %v", err)
		}
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Logf("Warning: failed to terminate Redis: %v", err)
		}
		t.Log("Cleanup complete")
	}

	return &TestEnv{
		DB:             db,
		Redis:          redisClient,
		Server:         server,
		Client:         NewTestClient(server.URL),
		pgContainer:    pgContainer,
		redisContainer: redisContainer,
		Cleanup:        cleanup,
	}
}

func generateRandomSuffix() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (e *TestEnv) SetCommunityID(id string) {
	e.CommunityID = id
}

func (e *TestEnv) QueryDB(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	return e.DB.Query(ctx, query, args...)
}

func (e *TestEnv) ExecDB(ctx context.Context, query string, args ...interface{}) error {
	_, err := e.DB.Exec(ctx, query, args...)
	return err
}
