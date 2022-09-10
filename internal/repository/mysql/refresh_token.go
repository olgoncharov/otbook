package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dto "github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

func (r *Repository) ReplaceRefreshToken(ctx context.Context, username string, newToken dto.RefreshToken) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ReplaceRefreshToken: can't begin transaction: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE user = ?", username)
	if err != nil {
		return fmt.Errorf("ReplaceRefreshToken: can't delete old refresh tokens: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO refresh_tokens (user, token, created_at, ttl)
		 VALUES (?, ?, ?, ?)`,
		username, newToken.Value, newToken.CreatedAt, newToken.TTL,
	)
	if err != nil {
		return fmt.Errorf("ReplaceRefreshToken: can't save new refresh token: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ReplaceRefreshToken: can't commit tx: %w", err)
	}

	return nil
}

func (r *Repository) GetRefreshTokenForUser(ctx context.Context, username string) (*dto.RefreshToken, error) {
	token := dto.RefreshToken{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT token, created_at, ttl
		 FROM refresh_tokens
		 WHERE user = ?`,
		username,
	).Scan(
		&token.Value,
		&token.CreatedAt,
		&token.TTL,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repoErrors.ErrNoRowsFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetRefreshTokenForUser: %w", err)
	}

	return &token, nil
}

func (r *Repository) DeleteRefreshTokenForUser(ctx context.Context, username string) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM refresh_tokens WHERE user = ?",
		username,
	)

	if err != nil {
		return fmt.Errorf("DeleteRefreshTokenForUser: %w", err)
	}

	return nil
}
