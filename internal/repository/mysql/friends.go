package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/olgoncharov/otbook/internal/entity"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

func (r *Repository) CreateFriends(ctx context.Context, firstUsername, secondUsername string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO friends (user1, user2)
		 VALUES (?, ?)`,
		firstUsername, secondUsername,
	)

	var mysqlErr *mysql.MySQLError

	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == duplicateEntryErrorCode {
			return fmt.Errorf("CreateFriends: %w", repoErrors.ErrUniqueConstraintViolated)
		}

		return fmt.Errorf("CreateFriends: %w", err)
	}

	return nil
}

func (r *Repository) GetFriendsOfUser(ctx context.Context, username string, limit, offset uint) ([]entity.Profile, error) {
	friends := make([]entity.Profile, 0, limit)

	query := queryWithLimitOffset(
		`SELECT
			user,
			first_name,
			last_name,
			birthdate,
			city,
			sex,
			hobby
		FROM profiles
		WHERE user IN (
			SELECT
				CASE WHEN user1 = ? THEN user2 ELSE user1 END AS user
			FROM friends
			WHERE user1 = ? OR user2 = ?
		)
		ORDER BY user`, limit, offset,
	)

	rows, err := r.db.QueryContext(ctx, query, username, username, username)
	if err != nil {
		return nil, fmt.Errorf("GetFriendsOfUser: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var profile entity.Profile
		err = rows.Scan(
			&profile.Username,
			&profile.FirstName,
			&profile.LastName,
			&profile.Birthdate,
			&profile.City,
			&profile.Sex,
			&profile.Hobby,
		)
		if err != nil {
			return nil, fmt.Errorf("GetFriendsOfUser: %w", err)
		}

		friends = append(friends, profile)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFriendsOfUser: %w", err)
	}

	return friends, nil
}

func (r *Repository) GetCountOfFriends(ctx context.Context, username string) (uint, error) {
	var totalCount uint

	err := r.db.QueryRowContext(
		ctx,
		`SELECT count(*)
		 FROM friends
		 WHERE user1 = ? OR user2 = ?`,
		username, username,
	).Scan(&totalCount)

	if err != nil {
		return 0, fmt.Errorf("GetCountOfFriends: %w", err)
	}

	return totalCount, nil
}
