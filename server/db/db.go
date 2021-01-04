package db

import (
	"context"
	"fmt"

	"github.com/UN0wen/virtualfs/server/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// Db contains the database connection pool and its config
type Db struct {
	Pool *pgxpool.Pool
	cfg  Config
}

// Config is a database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	URL      string
}

// Setup setups the database
func Setup(cfg Config) (Db, error) {
	var db Db
	if cfg.URL == "" {
		if cfg.Host == "" || cfg.Port == "" || cfg.User == "" ||
			cfg.Password == "" || cfg.Database == "" {
			err := errors.New("Provide all fields for config")
			return db, err
		}

		db.cfg = cfg
		cfgDNS := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)

		config, err := pgxpool.ParseConfig(cfgDNS)
		if err != nil {
			err = errors.Wrapf(err, "Cannot parse config string")
			return db, err
		}

		// PGXPool configs
		config.MaxConns = 10

		pool, err := pgxpool.ConnectConfig(context.Background(), config)
		if err != nil {
			err = errors.Wrapf(err, "Unable to connect to database")
			return db, err
		}

		db.Pool = pool
		utils.Sugar.Infof("Database created with config string: %s", config.ConnString())

	} else {
		pool, err := pgxpool.Connect(context.Background(), cfg.URL)
		if err != nil {
			err = errors.Wrapf(err, "Unable to connect to database")
			return db, err
		}
		db.Pool = pool
		utils.Sugar.Infof("Database created with URL: %s", cfg.URL)

	}
	return db, nil
}

// Close closes the connection pool
func (db *Db) Close() {
	if db.Pool == nil {
		return
	}

	db.Pool.Close()
	return
}

// CreateTable executes the query given
func (db *Db) CreateTable(query string) (err error) {
	utils.Sugar.Infof("SQL Query: %s", query)

	if _, err = db.Pool.Exec(context.Background(), query); err != nil {
		err = errors.Wrapf(err, "Table creation query failed")
	}

	return
}
