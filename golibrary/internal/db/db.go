package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"test/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func NewSqlDB(cfg *config.Config, logger *zap.Logger) (*sqlx.DB, error) {
	fmt.Println("dbname", cfg.Db.Dsn)
	switch cfg.Db.DBname {

	case "postgres":
		db, err := sql.Open("postgres", cfg.Db.Dsn)
		if err != nil {
			return nil, err
		}

		db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
		db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
		duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
		if err != nil {
			return nil, err
		}
		db.SetConnMaxIdleTime(duration)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			return nil, err
		}

		return sqlx.NewDb(db, "postgres"), nil
	case "mysql":

	case "sqlite3":
		rawdb, _ := sql.Open("sqlite3", ":memory:")
		_, err := rawdb.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, email VARCHAR(100), hashed_password VARCHAR(100))")
		if err != nil {
			logger.Fatal(err.Error())
		}
		db := sqlx.NewDb(rawdb, "sqlite3")
		return db, nil
	case "test":
		return nil, nil
	}

	return nil, errors.New("unknown DB type")
}
