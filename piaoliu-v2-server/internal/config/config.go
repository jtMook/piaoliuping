package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret string

	COSBucketURL string
	COSSecretID  string
	COSSecretKey string

	WxAppID     string
	WxAppSecret string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "drift_bottle"),

		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,

		JWTSecret: getEnv("JWT_SECRET", "drift-bottle-secret-key"),

		COSBucketURL: getEnv("COS_BUCKET_URL", ""),
		COSSecretID:  getEnv("COS_SECRET_ID", ""),
		COSSecretKey: getEnv("COS_SECRET_KEY", ""),

		WxAppID:     getEnv("WX_APP_ID", ""),
		WxAppSecret: getEnv("WX_APP_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
