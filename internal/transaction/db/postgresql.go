package transaction

import (
	"TestTask/internal/transaction"
	"TestTask/pkg/logging"
	"TestTask/pkg/repository/PostgreSQL"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

type repository struct {
	client PostgreSQL.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, trans *transaction.Transaction) error {
	q := `
		INSERT INTO transaction 
		    (sender, receiver, amount, status) 
		VALUES 
		       ($1, $2, $3, $4) 
		RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, trans.Sender, trans.Receiver, trans.Amount, trans.Status).Scan(&trans.Id); err != nil {
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

func (r *repository) FindAllNotDone(ctx context.Context) (t []transaction.Transaction, err error) {
	q := `
		SELECT id, sender, receiver, amount, status
		FROM public.transaction 
		WHERE status != 'closest';
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	transs := make([]transaction.Transaction, 0)
	for rows.Next() {
		var tr transaction.Transaction
		err = rows.Scan(&tr.Id, &tr.Sender, &tr.Receiver, &tr.Amount, &tr.Status)
		if err != nil {
			return nil, err
		}
		transs = append(transs, tr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transs, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (transaction.Transaction, error) {
	q := `
		SELECT id, name FROM public.transaction WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	var tr transaction.Transaction
	err := r.client.QueryRow(ctx, q, id).Scan(&tr.Id, &tr.Sender, &tr.Receiver, &tr.Amount, &tr.Status)
	if err != nil {
		return transaction.Transaction{}, err
	}
	return tr, nil
}

func (r *repository) Update(ctx context.Context, t transaction.Transaction) error {
	q := `
		UPDATE transaction
		SET status = $1
		WHERE id = $2;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	_, err := r.client.Exec(ctx, q, t.Status, t.Id)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(client PostgreSQL.Client, logger *logging.Logger) transaction.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
