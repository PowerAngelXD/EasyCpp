package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv             string        // 当前运行环境（development/staging/production）。
	HTTPAddr           string        // HTTP 服务监听地址。
	PostgresDSN        string        // PostgreSQL 连接字符串。
	RedisAddr          string        // Redis 连接地址。
	RedisPassword      string        // Redis 访问密码。
	RedisDB            int           // Redis 数据库编号。
	JWTSecret          string        // JWT 签名密钥。
	AccessTokenTTL     time.Duration // 访问令牌有效期。
	SessionTTL         time.Duration // 会话在 Redis 中的存活时长。
	PasswordBcryptCost int           // bcrypt 哈希强度参数。
}

func Load() Config {
	cfg := Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		HTTPAddr:           getEnv("HTTP_ADDR", ":8080"),
		PostgresDSN:        getEnv("POSTGRES_DSN", "postgres://postgres:postgres@127.0.0.1:5432/easycpp?sslmode=disable"),
		RedisAddr:          getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
		JWTSecret:          getEnv("JWT_SECRET", "change-me-in-production"),
		AccessTokenTTL:     getEnvAsDuration("JWT_ACCESS_TTL", 2*time.Hour),
		SessionTTL:         getEnvAsDuration("SESSION_TTL", 24*time.Hour),
		PasswordBcryptCost: getEnvAsInt("BCRYPT_COST", 12),
	}
	if cfg.PasswordBcryptCost < 10 {
		cfg.PasswordBcryptCost = 10
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
