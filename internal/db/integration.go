package integration

import (
	"context"
	"fmt"
	"net/url"
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
	dbHost string
	dbPort string
}

func getEnv() (envVar, error) {
	err := godotenv.Load("../../docker/.env")
	if err != nil {
		return envVar{}, fmt.Errorf("loading .env file: %w", err)
	}

	result := envVar{
		dbPSWD: url.QueryEscape(os.Getenv("DB_PASSWORD")),
		dbUser: url.QueryEscape(os.Getenv("DB_USER")),
		dbName: url.QueryEscape(os.Getenv("DB_NAME")),
		dbHost: os.Getenv("DB_HOST"),
		dbPort: os.Getenv("DB_PORT"),
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
		return envVar{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return result, nil
}

func Connect(ctx context.Context) (*pgxpool.Pool, error) {

	connVars, err := getEnv()
	if err != nil {
		return nil, fmt.Errorf("getting env: %w", err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", connVars.dbUser, connVars.dbPSWD, connVars.dbHost, connVars.dbPort, connVars.dbName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping db: %w", err)
	}

	return pool, nil
}
