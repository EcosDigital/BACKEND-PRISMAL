package core

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Backend_port     string
	Frontend_url     string
	Allow_subdomains bool

	Db_host string
	Db_port string
	Db_user string
	Db_pass string
	Db_name string

	Env         string
	TenancyMode string
	Jwt_secret  string

	//redis
	RedisAddr     string
	RedisPassword string
}

var Cfg Config

func LoadConfig() {

	envFile := ".env"

	// permite ejecutar: go run cmd/main.go --multi
	if len(os.Args) > 1 && os.Args[1] == "--multi" {
		envFile = ".multi.env"
	}

	//cargar archivo.ENV
	err := godotenv.Load(envFile)
	if err != nil {
		log.Println("No se pudo cargar .env, usando variables de entorno del sistema.")
	}

	Cfg = Config{
		Backend_port:     getEnv("BACKEND_PORT", "3000"),
		Frontend_url:     getEnv("FRONTEND_URL", "*"),
		Allow_subdomains: getEnvBool("ALLOW_SUBDOMAINS", false),

		Db_host: getEnv("DB_HOST", ""),
		Db_port: getEnv("DB_PORT", ""),
		Db_user: getEnv("DB_USER", ""),
		Db_pass: getEnv("DB_PASS", ""),
		Db_name: getEnv("DB_NAME", ""),

		Env:         getEnv("APP_ENV", ""),
		TenancyMode: getEnv("TENANCY_MODE", "single"),
		Jwt_secret:  getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(strings.TrimSpace(value)) == "true"
	}
	return fallback
}

func IsProduction() bool {
	return Cfg.Env == "production"
}
