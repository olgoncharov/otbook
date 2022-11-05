package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v9"
)

const (
	transactionMaxRetries = 5
)

type Repository struct {
	client redis.Client
}

func NewRepository(client redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) execWithTransaction(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	for i := 0; i < transactionMaxRetries; i++ {
		err := r.client.Watch(ctx, fn, keys...)
		if errors.Is(err, redis.TxFailedErr) {
			continue
		}

		return err
	}

	return errors.New("increment reached maximum number of retries")
}
