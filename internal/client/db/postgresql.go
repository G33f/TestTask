package client

import (
	"TestTask/internal/client"
	"TestTask/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)
import "TestTask/pkg/repository/PostgreSQL"

type repository struct {
	client PostgreSQL.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, client *client.Client) error {
	q := `
		INSERT INTO ClientValet 
		    (firstName, lastName, balance) 
		VALUES 
		       ($1, $2, $3) 
		RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, client.FirstName, client.LastName, client.Balance).Scan(&client.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) (c []client.Client, err error) {
	q := `
		SELECT id, firstName, lastName, balance FROM public.clientvalet;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	clients := make([]client.Client, 0)
	for rows.Next() {
		var cli client.Client
		err = rows.Scan(&cli.Id, &cli.FirstName, &cli.LastName, &cli.Balance)
		if err != nil {
			return nil, err
		}
		clients = append(clients, cli)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (client.Client, error) {
	q := `
		SELECT id, name FROM public.author WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	var c client.Client
	err := r.client.QueryRow(ctx, q, id).Scan(&c.Id, &c.FirstName, &c.LastName, &c.Balance)
	if err != nil {
		return client.Client{}, err
	}
	return c, nil
}

func (r *repository) Update(ctx context.Context, c client.Client) error {
	q := `
		UPDATE clientvalet
		SET balance = $2
		WHERE id = $1;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	err := r.client.QueryRow(ctx, q, c.Id, c.Balance).Scan()
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(client PostgreSQL.Client, logger *logging.Logger) client.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
