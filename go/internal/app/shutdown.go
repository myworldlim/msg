// internal/app/shutdown.go
// Обработка graceful shutdown для приложения.
package app

import (
    "context"
    "log"
    "os/signal"
    "syscall"
    "time"
)

// HandleShutdown запускает сервер и обрабатывает graceful shutdown
func HandleShutdown(app *App) error {
    // Запуск сервера в отдельной горутине
    go func() {
        if err := app.Run(); err != nil {
            log.Fatalf("Ошибка запуска приложения: %v", err)
        }
    }()

    // Обработка сигналов для graceful shutdown
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    <-ctx.Done()
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    log.Println("Останавливаем сервер...")
    if err := app.server.Shutdown(shutdownCtx); err != nil {
        log.Printf("Ошибка остановки сервера: %v", err)
        return err
    }

    log.Println("Закрываем пул соединений...")
    stats := app.db.Stat()
    log.Printf("Активных соединений перед закрытием: %d", stats.TotalConns())
    app.db.Close()
    log.Println("Пул соединений закрыт")
    return nil
}