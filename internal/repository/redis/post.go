package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	jsoniter "github.com/json-iterator/go"
	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

const (
	feedPostsKeyPattern = "feed:%s"
)

// GetPostFeed returns newest posts of friends. It does not contain posts from celebrity friends.
func (r *Repository) GetPostFeed(ctx context.Context, username string, limit uint) ([]dto.PostShortInfo, error) {
	key := fmt.Sprintf(feedPostsKeyPattern, username)
	exist, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("Exists failed: %w", err)
	}

	if exist == 0 {
		return nil, repoErrors.ErrCacheMiss
	}

	values, err := r.client.LRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("LRange failed: %w", err)
	}

	posts := make([]dto.PostShortInfo, 0, limit)
	for _, v := range values {
		if v == "" {
			continue
		}

		var post dto.PostShortInfo
		if err = jsoniter.UnmarshalFromString(v, &post); err != nil {
			return nil, fmt.Errorf("can't unmarshal value from cache: %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// SetPostFeed replaces current list of posts by new value.
func (r *Repository) SetPostFeed(ctx context.Context, username string, posts []dto.PostShortInfo) error {
	feedKey := fmt.Sprintf(feedPostsKeyPattern, username)

	return r.execWithTransaction(ctx, func(tx *redis.Tx) error {
		tx.Del(ctx, feedKey)
		for _, post := range posts {
			encodedPost, err := jsoniter.MarshalToString(post)
			if err != nil {
				return fmt.Errorf("can't marshal post: %w", err)
			}
			tx.RPush(ctx, feedKey, encodedPost)
		}
		if len(posts) == 0 {
			tx.RPush(ctx, feedKey, "")
		}
		return nil
	}, feedKey)
}

// PusthPostToFeed adds new post to feed of given users.
func (r *Repository) PushPostToFeed(ctx context.Context, post entity.Post, usernames []string, feedLimit uint) error {
	pipe := r.client.Pipeline()

	encodedPost, _ := jsoniter.MarshalToString(dto.PostShortInfo{
		ID:        post.ID,
		Author:    post.Author,
		Title:     post.Title,
		CreatedAt: post.CreatedAt,
	})

	for _, username := range usernames {
		key := fmt.Sprintf(feedPostsKeyPattern, username)
		pipe.LPush(ctx, key, encodedPost)
		pipe.LTrim(ctx, key, 0, int64(feedLimit))
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("can't execute pipeline: %w", err)
	}

	return nil
}

// DeletePostsFromFeed removes given posts from feed.
func (r *Repository) DeletePostsFromFeed(ctx context.Context, username string, posts []dto.PostShortInfo) error {
	if len(posts) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()
	key := fmt.Sprintf(feedPostsKeyPattern, username)

	for _, post := range posts {
		encodedPost, _ := jsoniter.MarshalToString(dto.PostShortInfo{
			ID:        post.ID,
			Author:    post.Author,
			Title:     post.Title,
			CreatedAt: post.CreatedAt,
		})
		_ = pipe.LRem(ctx, key, 0, encodedPost)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("can't execute pipeline: %w", err)
	}

	return nil
}
