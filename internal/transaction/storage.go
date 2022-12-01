package transaction

import "context"

type Repository interface {
	Create(ctx context.Context, c *Transaction) error
	FindAllNotDone(ctx context.Context) (cl []Transaction, err error)
	FindOne(ctx context.Context, id string) (Transaction, error)
	Update(ctx context.Context, c Transaction) error
}
