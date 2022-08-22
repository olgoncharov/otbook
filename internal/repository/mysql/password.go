package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

func (r *Repository) GetPasswordHashByUsername(ctx context.Context, username string) (string, error) {
	var password string

	err := r.db.QueryRowContext(
		ctx,
		"SELECT password FROM users WHERE username = ?",
		username,
	).Scan(&password)

	if errors.Is(err, sql.ErrNoRows) {
		return "", repoErrors.ErrNoRowsFound
	}
	if err != nil {
		return "", fmt.Errorf("GetPasswordHashByUsername: %w", err)
	}

	return password, nil
}
