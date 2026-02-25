// internal/app/app.go
// Основной файл приложения для инициализации и запуска сервера.
package app

import (
    "chitchat/config"
    "chitchat/internal/app/server"
    "github.com/jackc/pgx/v5/pgxpool"
    "log"
    "net/http"
)

// App представляет приложение
type App struct {
    db     *pgxpool.Pool
    server *http.Server
}

// NewApp создаёт новое приложение
func NewApp(db *pgxpool.Pool) *App {
    r := server.NewServer(db)
    serverHttp := &http.Server{
        Addr:    ":" + config.AppConfig.AppPort,
        Handler: r,
    }
    return &App{
        db:     db,
        server: serverHttp,
    }
}

// Run запускает сервер
func (a *App) Run() error {
    log.Printf("HTTP-сервер запущен на :%s", a.server.Addr)
    if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Printf("Ошибка запуска сервера: %v", err)
        return err
    }
    return nil
}