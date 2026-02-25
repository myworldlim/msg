// cmd/main.go
package main

import (
    "chitchat/config"
    "chitchat/internal/app"
    "log"
)

func main() {
    // Инициализация конфигурации
    config.InitConfig()

    // Инициализация базы данных
    dbPool := config.InitDB()

    // Создание приложения
    application := app.NewApp(dbPool)

    // Запуск сервера и обработка graceful shutdown
    if err := app.HandleShutdown(application); err != nil {
        log.Fatalf("Ошибка при выполнении приложения: %v", err)
    }
}