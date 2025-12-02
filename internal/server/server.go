package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"gks.com/gohl-test/internal/config"
	"gks.com/gohl-test/internal/handler"
		"gks.com/gohl-test/internal/repo"
	"go.uber.org/zap"
)

type Server struct {
	// Logger *zap.Logger
	http.Server
}

func NewServer(cfg *config.Config) (*Server, error) {
	
	dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.DBName,
        cfg.Database.SSLMode,
    )

    conn, err := pgx.Connect(context.Background(), dsn)
    if err != nil {
        return nil, fmt.Errorf("db connect failed: %w", err)
    }

    if err := conn.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("db ping failed: %w", err)
    }

	logger, _ := zap.NewProduction()

	r := mux.NewRouter()
	userRepository := repo.NewUserRepository(conn, logger)
	txnRepository := repo.NewTransactionsRepository(conn, logger)
	userHandler := handler.NewUserHandler(userRepository, logger)
	r.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.GetUserBalance).Methods("GET")
	txnHandler := handler.NewTransactionsHandler(txnRepository, userRepository, logger)
	r.HandleFunc("/transactions", txnHandler.ListTransactions).Methods("GET")
	r.HandleFunc("/transactions", txnHandler.CreateTransactions).Methods("POST")



	return &Server{
		Server: http.Server{
			Addr: fmt.Sprintf(":%s", cfg.Server.Port),
			Handler: r,
		},
	}, nil
}