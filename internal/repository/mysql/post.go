package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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

// GetPostFeed returns newest posts of friends.
func (r *Repository) GetPostFeed(ctx context.Context, username string, limit uint) ([]dto.PostShortInfo, error) {
	posts := make([]dto.PostShortInfo, 0, limit)
	query := queryWithLimitOffset(
		`SELECT
			p.id,
			p.author,
			p.title,
			p.created_at
		FROM
			posts AS p
			JOIN friends AS f ON p.author = f.friend AND f.user = ?
		ORDER BY p.created_at DESC`, limit, 0,
	)

	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("GetPostFeed: %w", err)
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
			return nil, fmt.Errorf("GetPostFeed: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPostFeed: %w", err)
	}

	return posts, nil
}

func (r *Repository) GetPostFeedWithoutCelebrities(ctx context.Context, username string, limit uint) ([]dto.PostShortInfo, error) {
	posts := make([]dto.PostShortInfo, 0, limit)
	query := queryWithLimitOffset(
		`SELECT
			posts.id,
			posts.author,
			posts.title,
			posts.created_at
		FROM
			posts
			JOIN friends ON posts.author = friends.friend AND friends.user = ?
			JOIN profiles ON posts.author = profiles.user
		WHERE profiles.is_celebrity = false
		ORDER BY posts.created_at DESC`, limit, 0,
	)

	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("GetPostFeedWithoutCelebrities: %w", err)
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
			return nil, fmt.Errorf("GetPostFeedWithoutCelebrities: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPostFeedWithoutCelebrities: %w", err)
	}

	return posts, nil
}

func (r *Repository) GetPostsByFilters(ctx context.Context, filters dto.PostFilters, limit uint, offset uint) ([]dto.PostShortInfo, error) {
	queryBuilder := sq.Select("id", "author", "title", "created_at").From("posts").OrderBy("created_at DESC")

	if len(filters.Authors) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"author": filters.Authors})
	}

	if filters.DateFrom != nil {
		queryBuilder = queryBuilder.Where(sq.LtOrEq{"created_at": filters.DateFrom})
	}

	if limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(limit))
	}

	if offset > 0 {
		queryBuilder = queryBuilder.Offset(uint64(offset))
	}

	rows, err := queryBuilder.RunWith(r.db).Query()
	if err != nil {
		return nil, fmt.Errorf("GetPostsByFilters: %w", err)
	}
	defer rows.Close()

	posts := make([]dto.PostShortInfo, 0, limit)

	for rows.Next() {
		var post dto.PostShortInfo
		err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Title,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetPostsByFilter: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPostsByFilter: %w", err)
	}

	return posts, nil
}
