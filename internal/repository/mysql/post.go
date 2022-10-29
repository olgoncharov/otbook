package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/olgoncharov/otbook/internal/entity"
	"github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

func (r *Repository) GetAllPosts(ctx context.Context, limit, offset uint) ([]dto.PostShortInfo, error) {
	posts := make([]dto.PostShortInfo, 0, limit)
	query := queryWithLimitOffset(
		`SELECT
			id,
			author,	
			title,
			created_at
		FROM posts
		ORDER BY created_at DESC`, limit, offset)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetAllPosts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post dto.PostShortInfo
		err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Title,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAllPosts: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllPosts: %w", err)
	}

	return posts, nil
}

func (r *Repository) GetPostsCount(ctx context.Context) (uint, error) {
	var totalCount uint

	err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM posts").Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("GetPostsCount: %w", err)
	}

	return totalCount, nil
}

func (r *Repository) GetPostByID(ctx context.Context, postID uint64) (*entity.Post, error) {
	var post entity.Post

	err := r.db.QueryRowContext(
		ctx,
		`SELECT
			id,
			author,
			title,
			text,
			created_at
		FROM posts
		WHERE id = ?`,
		postID,
	).Scan(
		&post.ID,
		&post.Author,
		&post.Title,
		&post.Text,
		&post.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repoErrors.ErrNoRowsFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetPostByID: %w", err)
	}

	return &post, nil
}

func (r *Repository) CreatePost(ctx context.Context, post entity.Post) (uint64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO posts(author,title,text,created_at)
		VALUES (?,?,?,?)`,
		post.Author, post.Title, post.Text, post.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("CreatePost: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreatePost: %w", err)
	}

	return uint64(id), nil
}
