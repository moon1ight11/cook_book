package main

import (
	"context"
	"cook_book/internal/app"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// обработка паники
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered: %v\n", r)
			os.Exit(1)
		}
	}()

	ctx := context.Background()

	// инициализация всех зависимостей
	deps, err := app.NewDependencies(ctx)
	if err != nil {
		deps.Logger.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	// запуск очистителя токенов
	deps.StartCleanup(ctx)

	// создаем каналы для сигналов завершения и ошибки сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	serverError := make(chan error, 1)

	// запуск сервера
	go func() {
		deps.Logger.Info("Server is starting", "port:", deps.Config.Server.Port)
		if err := deps.Server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			deps.Logger.Error("Failed to run server", "error", err)
			serverError <- err
		}
	}()

	// ждем сигналы
	select {
	case <-quit:
		// если поступил сигнал завершения - делаем шатдаун с таймаутом
		deps.Logger.Info("Shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// завершаем работу сервера
		if err := deps.Server.Shutdown(shutdownCtx); err != nil {
			deps.Logger.Error("Server forced to shutdown", "error", err)
		}

		deps.Logger.Info("Disconnecting from database")

		// завершаем работу с БД
		if err := deps.CloseDB(shutdownCtx); err != nil {
			deps.Logger.Error("Failed to close database connection", "error", err)
		} else {
			deps.Logger.Info("Database connection closed successfully")
		}

		// останавливаем очиститель токенов
		deps.StopCleanup()

		// останавливаем кэш
		deps.StopCache()

	case err := <-serverError:
		// если пришла ошибка от сервера - фаталим
		deps.Logger.Error("Server error", "error", err)
		os.Exit(1)
	}

	deps.Logger.Info("Server exited")
}
