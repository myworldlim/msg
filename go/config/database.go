// config/database.go
// Подключение к базе данных (PostgreSQL). Возвращает пул соединений.
package config

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/jackc/pgx/v5/pgxpool"
)

// InitDB инициализирует подключение к базе данных
func InitDB() *pgxpool.Pool {
    connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
        AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBName, AppConfig.DBSSLMode)

    poolConfig, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        log.Fatalf("Невозможно проанализировать конфигурацию базы данных: %v", err)
    }

    // Настройка пула соединений
    poolConfig.MaxConns = int32(AppConfig.DBMaxConns) // Максимум соединений
    poolConfig.MinConns = int32(AppConfig.DBMinConns) // Минимум соединений
    poolConfig.MaxConnLifetime = 30 * time.Minute     // Время жизни соединения
    poolConfig.MaxConnIdleTime = 5 * time.Minute      // Время простоя соединения
    poolConfig.HealthCheckPeriod = 1 * time.Minute    // Период проверки здоровья

    // Повторные попытки подключения
    maxRetries := 5
    retryDelay := 2 * time.Second
    var dbPool *pgxpool.Pool

    for attempt := 1; attempt <= maxRetries; attempt++ {
        dbPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
        if err == nil {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            err = dbPool.Ping(ctx)
            cancel()
            if err == nil {
                log.Printf("Успешное подключение к базе данных %s", AppConfig.DBName)
                return dbPool
            }
        }

        log.Printf("Попытка %d: не удалось подключиться к базе данных: %v", attempt, err)
        if attempt < maxRetries {
            time.Sleep(retryDelay)
            retryDelay *= 2 // Экспоненциальная задержка
        }
    }

    log.Fatalf("Не удалось подключиться к базе данных после %d попыток", maxRetries)
    return nil
}