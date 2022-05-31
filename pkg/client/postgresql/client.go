package postgresql

import (
	"context"
	"critical-path-analysis-api/internal/config"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, attempts int, config *config.Config) (*pgxpool.Pool, error) {
	cs := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)

	for attempts > 0 {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		pool, err := pgxpool.Connect(ctx, cs)
		if err != nil {
			time.Sleep(10 * time.Second)
			attempts--
		} else {
			return pool, nil
		}
	}
	return nil, errors.New("connection to database is failed")
}
