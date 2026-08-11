package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cook_book/internal/app"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered: %v\n", r)
			os.Exit(1)
		}
	}()

	ctx := context.Background()

	deps, err := app.NewDependencies(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	deps.StartCleanup(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	serverError := make(chan error, 1)

	go func() {
		deps.Logger.Info("Server is starting", "port", deps.Config.Server.Port)
		if err := deps.Server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			deps.Logger.Error("Failed to run server", "error", err)
			serverError <- err
		}
	}()

	select {
	case <-quit:
		deps.Logger.Info("Shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := deps.Server.Shutdown(shutdownCtx); err != nil {
			deps.Logger.Error("Server forced to shutdown", "error", err)
		}

		deps.Logger.Info("Disconnecting from database")

		if err := deps.CloseDB(shutdownCtx); err != nil {
			deps.Logger.Error("Failed to close database connection", "error", err)
		} else {
			deps.Logger.Info("Database connection closed successfully")
		}

		deps.StopCleanup()
		deps.StopCache()

	case err := <-serverError:
		deps.Logger.Error("Server error", "error", err)
		os.Exit(1)
	}

	deps.Logger.Info("Server exited")
}
