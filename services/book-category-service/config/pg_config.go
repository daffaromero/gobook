package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"

	logs "github.com/daffaromero/gobook/services/common/helper/logger"
	"github.com/daffaromero/gobook/services/common/utils"
)

type DBConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	MinConns        int32
	MaxConns        int32
	TimeOutDuration time.Duration
}

func loadDBConfig() (*DBConfig, error) {
	minConns, err := strconv.Atoi(utils.GetEnv("DB_MIN_CONNS"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MIN_CONNS: %w", err)
	}

	maxConns, err := strconv.Atoi(utils.GetEnv("DB_MAX_CONNS"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONNS: %w", err)
	}

	timeoutDuration, err := strconv.Atoi(utils.GetEnv("DB_CONNECTION_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONNECTION_TIMEOUT: %w", err)
	}

	return &DBConfig{
		Host:            utils.GetEnv("DB_HOST"),
		Port:            utils.GetEnv("DB_PORT"),
		Username:        utils.GetEnv("DB_USERNAME"),
		Password:        utils.GetEnv("DB_PASSWORD"),
		DBName:          utils.GetEnv("DB_NAME"),
		MinConns:        int32(minConns),
		MaxConns:        int32(maxConns),
		TimeOutDuration: time.Duration(timeoutDuration) * time.Second,
	}, nil
}

func NewPostgresDatabase() (*pgxpool.Pool, error) {
	logger := logs.New("database_connection")

	dbConfig, err := loadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database configuration: %w", err)
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("Failed to parse configuration dsn " + dsn)
	}

	poolConfig.MinConns = dbConfig.MinConns
	poolConfig.MaxConns = dbConfig.MaxConns
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeDescribeExec

	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.TimeOutDuration)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected dsn " + dsn)

	return pool, nil
}
