package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV             string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	AccessTokenKey  string
	JwtSecret       string
	JwtTtl          string
	IsAuthKey       string
	HttpSecure      string
	HttpOnly        string
	Domain          string
	RedisHost       string
	AwsMailFrom     string
	AwsMailRegion   string
	AwsAccessKey    string
	AwsSecretKey    string
	CorsDomains     string
	AwsBucket       string
	AwsBucketRegion string
	AwsSqsRegion    string
	AwsSqsUrl       string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using environment variables.")
	}

	return &Config{
		ENV:             getEnv("ENV", "local"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "3306"),
		DBUser:          getEnv("DB_USER", "root"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "gcstatus"),
		AccessTokenKey:  getEnv("ACCESS_TOKEN_KEY", "_gc_9hp1b73cGDCmAPgaVTYOlS6cjPsnDYho"),
		JwtSecret:       getEnv("JWT_SECRET", "5qY51df4G2WkfGhYxsB2bO5yXhc5RG9l"),
		JwtTtl:          getEnv("JWT_TTL", "7"), // in days
		IsAuthKey:       getEnv("IS_AUTH_KEY", "_gc_auth"),
		HttpSecure:      getEnv("HTTP_SECURE", "false"),
		HttpOnly:        getEnv("HTTP_ONLY", "false"),
		Domain:          getEnv("HTTP_DOMAIN", "localhost"),
		RedisHost:       getEnv("REDIS_HOST", "localhost:6379"),
		AwsMailFrom:     getEnv("AWS_MAIL_FROM", "localhost@localhost.com"),
		AwsMailRegion:   getEnv("AWS_MAIL_REGION", "us-west-2"),
		AwsAccessKey:    getEnv("AWS_ACCESS_KEY", ""),
		AwsSecretKey:    getEnv("AWS_SECRET_KEY", ""),
		CorsDomains:     getEnv("CORS_DOMAINS", "http://localhost:5173"),
		AwsBucket:       getEnv("AWS_BUCKET", ""),
		AwsBucketRegion: getEnv("AWS_BUCKET_REGION", ""),
		AwsSqsRegion:    getEnv("AWS_SQS_REGION", ""),
		AwsSqsUrl:       getEnv("AWS_SQS_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetDBConnectionURL(config *Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
}

var JWTSecret []byte

func init() {
	config := LoadConfig()

	JWTSecret = []byte(config.JwtSecret)
}
