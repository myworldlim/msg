// config/config.go
// Загрузка и инициализация конфигурации из .env файла.
// Предоставляет глобальную переменную Config.
package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	AppPort        string
	WSPort         string
	FrontendOrigin string
	BackendOrigin  string
	WSOrigin       string
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
	DBName         string
	DBSSLMode      string
	CORSAccepted   []string
	DBMaxConns     int
	DBMinConns     int
	APIPrefix      string
	JWTSecret      string
}

var AppConfig Config

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if fallback == "" {
		log.Fatalf("Переменная окружения %s обязательна", key)
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	log.Printf("Предупреждение: некорректное значение для %s, используем значение по умолчанию %d", key, fallback)
	return fallback
}

func InitConfig() {
	appEnv := getEnv("APP_ENV", "development")

	if err := godotenv.Load(".env." + appEnv); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Файл .env не найден, используем переменные окружения")
		}
	}

	AppConfig = Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "8181"),
		WSPort:         getEnv("WS_PORT", "8182"),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		BackendOrigin:  getEnv("BACKEND_ORIGIN", "http://localhost:8181"),
		WSOrigin:       getEnv("WS_ORIGIN", "http://localhost:8182"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBName:         getEnv("DB_NAME", "chitchat"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		CORSAccepted:   splitEnv("FRONTEND_ORIGIN", ","),
		DBMaxConns:     getEnvInt("DB_MAX_CONNS", 50),
		DBMinConns:     getEnvInt("DB_MIN_CONNS", 5),
		APIPrefix:      getEnv("API_PREFIX", "/api/"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
	}

	if AppConfig.DBPassword == "" {
		log.Fatal("DB_PASSWORD обязателен")
	}
	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET обязателен")
	}
	if AppConfig.DBUser == "" || AppConfig.DBHost == "" ||
		AppConfig.DBPort == "" || AppConfig.DBName == "" {
		log.Fatal("Необходимо указать все параметры базы данных")
	}

	if _, err := strconv.Atoi(AppConfig.AppPort); err != nil {
		log.Fatalf("APP_PORT должен быть числом: %v", err)
	}
	if _, err := strconv.Atoi(AppConfig.WSPort); err != nil {
		log.Fatalf("WS_PORT должен быть числом: %v", err)
	}
}

func splitEnv(key, sep string) []string {
	// Используем os.Getenv чтобы не фаталить при отсутствии переменной
	value := os.Getenv(key)
	if value == "" {
		return []string{}
	}
	values := strings.Split(value, sep)
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}
	return values
}