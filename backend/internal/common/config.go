package common

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `validate:"required"`
	Database DatabaseConfig `validate:"required"`
	Redis    RedisConfig    `validate:"required"`
	Auth     AuthConfig     `validate:"required"`
	Email    EmailConfig    `validate:"required"`
}

type ServerConfig struct {
	Port         int           `validate:"required,min=1,max=65535"`
	ReadTimeout  time.Duration `validate:"required"`
	WriteTimeout time.Duration `validate:"required"`
}

type DatabaseConfig struct {
	Host     string `validate:"required,hostname"`
	Port     int    `validate:"required,min=1,max=65535"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Database string `validate:"required"`
	PoolSize int    `validate:"min=1"`
}

type RedisConfig struct {
	Host     string `validate:"required,hostname"`
	Port     int    `validate:"required,min=1,max=65535"`
	Password string
	DB       int `validate:"min=0,max=15"`
}

type AuthConfig struct {
	SessionTTL     time.Duration `validate:"required,min=1m"`
	PasswordMinLen int           `validate:"required,min=4"`
}

type EmailConfig struct {
	SMTPHost     string `validate:"required,hostname"`
	SMTPPort     int    `validate:"required,min=1,max=65535"`
	SMTPUser     string
	SMTPPassword string
	FromEmail    string `validate:"required,email"`
	FromName     string `validate:"required"`
}

func Validate(cfg *Config) error {
	validate := validator.New()
	return validate.Struct(cfg)
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.projek")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	cfg := &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "projek",
			Password: "projek",
			Database: "projek",
			PoolSize: 20,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		Auth: AuthConfig{
			SessionTTL:     7 * 24 * time.Hour,
			PasswordMinLen: 8,
		},
		Email: EmailConfig{
			SMTPHost:  "localhost",
			SMTPPort:  587,
			FromEmail: "noreply@projek.local",
			FromName:  "Projek",
		},
	}

	viper.RegisterAlias("server.port", "PORT")
	viper.RegisterAlias("database.host", "DB_HOST")
	viper.RegisterAlias("database.port", "DB_PORT")
	viper.RegisterAlias("database.user", "DB_USER")
	viper.RegisterAlias("database.password", "DB_PASSWORD")
	viper.RegisterAlias("database.database", "DB_NAME")
	viper.RegisterAlias("redis.host", "REDIS_HOST")
	viper.RegisterAlias("redis.port", "REDIS_PORT")
	viper.RegisterAlias("redis.password", "REDIS_PASSWORD")

	viper.Unmarshal(cfg)

	return cfg
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		c.User, c.Password, c.Host, c.Port, c.Database, c.PoolSize,
	)
}

func NewDatabase(cfg *Config, logger zerolog.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 5; i++ {
		pool, err = pgxpool.New(ctx, cfg.Database.DSN())
		if err == nil {
			if err := pool.Ping(ctx); err == nil {
				logger.Info().Str("host", cfg.Database.Host).Int("port", cfg.Database.Port).Msg("database connected")
				return pool, nil
			}
			pool.Close()
		}
		logger.Warn().Int("attempt", i+1).Err(err).Msg("database connection failed, retrying...")
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
}

func DatabaseHealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	return pool.Ping(ctx)
}

func NewRedis(cfg *Config, logger zerolog.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info().Str("host", cfg.Redis.Host).Int("port", cfg.Redis.Port).Msg("redis connected")
	return client, nil
}
