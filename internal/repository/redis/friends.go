package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

const (
	celebrityFriendsKeyPattern = "celebrity_friends:%s"
)

// GetCelebrityFriends returns usernames of friends which are celebrity persons.
func (r *Repository) GetCelebrityFriends(ctx context.Context, username string) ([]string, error) {
	key := fmt.Sprintf(celebrityFriendsKeyPattern, username)

	exist, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("Exists failed: %w", err)
	}

	if exist == 0 {
		return nil, repoErrors.ErrCacheMiss
	}

	values, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("SMemebers failed: %w", err)
	}

	result := make([]string, 0, len(values))
	for _, v := range values {
		if v != "" {
			result = append(result, v)
		}
	}

	return result, nil
}

// SetCelebrityFriends replaces collection of celebrity friends for given user.
func (r *Repository) SetCelebrityFriends(ctx context.Context, username string, friends []string) error {
	key := fmt.Sprintf(celebrityFriendsKeyPattern, username)

	return r.execWithTransaction(ctx, func(tx *redis.Tx) error {
		tx.Del(ctx, key)
		for _, f := range friends {
			tx.SAdd(ctx, key, f)
		}
		if len(friends) == 0 {
			tx.SAdd(ctx, key, "")
		}
		return nil
	}, key)
}

// AddCelebrityFriend pushes new celebrity friend to set.
func (r *Repository) AddCelebrityFriend(ctx context.Context, user string, friend string) error {
	_, err := r.client.SAdd(ctx, fmt.Sprintf(celebrityFriendsKeyPattern, user), friend).Result()

	return err
}

// DeleteCelebrityFriend removes celebrity friend from set.
func (r *Repository) DeleteCelebrityFriend(ctx context.Context, user string, friend string) error {
	_, err := r.client.SRem(ctx, fmt.Sprintf(celebrityFriendsKeyPattern, user), friend).Result()

	return err
}
