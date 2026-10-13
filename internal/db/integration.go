package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type envVar struct {
	dbPSWD string
	dbUser string
	dbName string
}

func getEnv() (envVar, error) {
	err := godotenv.Load("../../docker/.env")
	if err != nil {
		return envVar{}, fmt.Errorf("loading .env file: %w\n", err)
	}

	result := envVar{
		dbPSWD: os.Getenv("DB_PASSWORD"),
		dbUser: os.Getenv("DB_USER"),
		dbName: os.Getenv("DB_NAME"),
	}

	var missing []string
	if result.dbPSWD == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if result.dbUser == "" {
		missing = append(missing, "DB_USER")
	}
	if result.dbName == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return envVar{}, fmt.Errorf("missing required env vars: %s\n", strings.Join(missing, ", "))
	}

	return result, nil
}

func Connect(ctx context.Context) (*pgxpool.Pool, error) {

	connVars, err := getEnv()
	if err != nil {
		return nil, fmt.Errorf("getting env: %v\n", err)
	}

	connString := fmt.Sprintf("postgress://%s:%s@localhost:5432/%s?sslmode=disable", connVars.dbUser, connVars.dbPSWD, connVars.dbName) // TODO: handle URL-encode credentials

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pgxpool: %w\n", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping db: %w\n", err)
	}

	return pool, nil
}
