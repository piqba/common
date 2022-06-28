package postgres

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/piqba/common/pkg/databases"
	"time"
)

// NewPostgresDb func for connection to Postgresql database.
func NewPostgresDb(option databases.PgOptions) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", option.DbURI)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set database connection settings.
	db.SetMaxOpenConns(option.MaxConnections)                           // the default is 0 (unlimited)
	db.SetMaxIdleConns(option.MaxIdleConnections)                       // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(time.Duration(option.MaxLifeTimeConnections)) // 0, connections are reused forever

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer func(db *sqlx.DB) {
			_ = db.Close()
		}(db) // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}

// NewPostgresDbPool use a pgxpool is a connection pool for pgx. pgx is entirely decoupled from its default pool implementation. This means that pgx can be used with a different pool or without any pool at all.
func NewPostgresDbPool(ctx context.Context, option databases.PgOptions) (pool *pgxpool.Pool, err error) {
	pool, err = pgxpool.Connect(ctx, option.DbURI)
	if err != nil {
		return
	}

	return
}
