package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/tmc/langchaingo/memory/sqlite3"

	_ "github.com/mattn/go-sqlite3"
)

var db *pgxpool.Pool

// Init the db
func Init() error {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		return errors.New("DATABASE URL environment variable not set")
	}

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return err
	}

	db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	log.Info().Msg("Database connected successfully")

	// Close the database connection when the application is interrupted
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		Close()
		log.Info().Msg("Application stopped")
		os.Exit(0)
	}()

	return nil
}

// Close closes the database connection pool.
func Close() {
	if db != nil {
		db.Close()
		log.Info().Msg("Closed the database connection")
	}
}

// GetDB returns the database connection pool.
func GetDB() *pgxpool.Pool {
	return db
}

func GetSqlite(dbName string, sessionID string) (*sql.DB, *sqlite3.SqliteChatMessageHistory, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s.db", dbName))
	if err != nil {
		return nil, nil, err
	}

	chatHistory := sqlite3.NewSqliteChatMessageHistory(
		sqlite3.WithSession(sessionID),
		sqlite3.WithDB(db),
	)

	return db, chatHistory, nil
}