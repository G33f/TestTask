package PostgreSQL

import (
	"context"
	"example.com/m/author/repository/PostgreSQL/utils"
	"example.com/m/author/repository/PostgreSQL/utils/config"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
}

func NewClient(ctx context.Context, sc config.StorageConfig) (err error, pool *pgxpool.Pool) {
	Username, Password, Host, Port, Database := sc.GetStorageConfig()
	dsn := fmt.Sprintf("postgressql://%s:%s@%s:%s/%s", Username, Password, Host, Port, Database)
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, sc.GetMaxAttempt(), 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}
	return
}
