package main

import (
	"context"
	"expense_tracker/internal/config"
	"expense_tracker/internal/handler"
	"expense_tracker/internal/middleware"
	"expense_tracker/internal/repository"
	"expense_tracker/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"os/signal"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
)

func main() {

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := repository.NewDB(ctx, cfg)
	if err != nil {
		log.Fatalf("cmd: failed to connect to database: %v", err)
	}
	defer db.Pool.Close()

	if err := applyMigrations(cfg); err != nil {
		log.Fatalf("cmd: migrations failed: %v", err)
	}

	userRep := repository.NewUserRepository(db)
	expenseRep := repository.NewExpenseRepository(db)

	authService := service.NewAuthService(
		userRep,
		cfg.JWTSecret,
		24*time.Hour,
	)
	userService := service.NewUserServise(userRep)
	expeneseService := service.NewExpenseServise(expenseRep)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	expenseHandler := handler.NewExpenseHandler(expeneseService)

	router := http.NewServeMux()
	authMiddleware := middleware.AuthMiddleware(authService)

	router.HandleFunc("POST /auth/register", authHandler.Register)
	router.HandleFunc("POST /auth/login", authHandler.Login)

	router.Handle("PUT /user/username", authMiddleware(http.HandlerFunc(userHandler.UpdateUsername)))
	router.Handle("DELETE /user", authMiddleware(http.HandlerFunc(userHandler.DeleteUser)))
	router.Handle("GET /user", authMiddleware(http.HandlerFunc(userHandler.GetProfile)))

	router.Handle("POST /expenses", authMiddleware(http.HandlerFunc(expenseHandler.CreateExpense)))
	router.Handle("GET /expenses/{id}", authMiddleware(http.HandlerFunc(expenseHandler.GetExpense)))
	router.Handle("PUT /expenses/{id}", authMiddleware(http.HandlerFunc(expenseHandler.UpdateExpense)))
	router.Handle("DELETE /expenses/{id}", authMiddleware(http.HandlerFunc(expenseHandler.DeleteExpense)))
	router.Handle("GET /expenses", authMiddleware(http.HandlerFunc(expenseHandler.GetExpensesList)))
	router.Handle("GET /expenses/period", authMiddleware(http.HandlerFunc(expenseHandler.GetExpensesByPeriod)))
	router.Handle("GET /expenses/category", authMiddleware(http.HandlerFunc(expenseHandler.GetExpensesByCategory)))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("cmd: server shutdown error: %v", err)
		}
	}()

	log.Printf("cmd: server starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cmd: server failed: %v", err)
	}
}

func applyMigrations(cfg *config.Config) error {
	dsn := cfg.DBURL

	m, err := migrate.New(
		"file:///app/migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("cmd: migration init failed: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("cmd: migration up failed: %w", err)
	}

	return nil
}
