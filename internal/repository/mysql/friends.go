package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/olgoncharov/otbook/internal/entity"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

func (r *Repository) AddFriend(ctx context.Context, user, newFriend string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO friends (user, friend)
		 VALUES (?, ?)`,
		user, newFriend,
	)

	var mysqlErr *mysql.MySQLError

	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == duplicateEntryErrorCode {
			return fmt.Errorf("AddFriend: %w", repoErrors.ErrUniqueConstraintViolated)
		}

		return fmt.Errorf("AddFriend: %w", err)
	}

	return nil
}

func (r *Repository) DeleteFriend(ctx context.Context, user, friend string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM friends WHERE user = ? AND friend = ?`,
		user, friend,
	)

	if err != nil {
		return fmt.Errorf("DeleteFriend: %w", err)
	}

	return nil
}

func (r *Repository) GetFriendsOfUser(ctx context.Context, username string, limit, offset uint) ([]entity.Profile, error) {
	friends := make([]entity.Profile, 0, limit)

	query := queryWithLimitOffset(
		`SELECT
			p.user,
			p.first_name,
			p.last_name,
			p.birthdate,
			p.city,
			p.sex,
			p.hobby
		FROM
			profiles AS p
			JOIN friends AS f ON p.user = f.friend AND f.user = ?
		ORDER BY p.user`, limit, offset,
	)

	rows, err := r.db.QueryContext(ctx, query, username)
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
		 WHERE user = ?`,
		username,
	).Scan(&totalCount)

	if err != nil {
		return 0, fmt.Errorf("GetCountOfFriends: %w", err)
	}

	return totalCount, nil
}

// GetCelebrityFriends returns usernames of friends which are celebrity persons.
func (r *Repository) GetCelebrityFriends(ctx context.Context, username string) ([]string, error) {
	query :=
		`SELECT
			f.friend
		 FROM
		 	friends AS f
			JOIN profiles AS p ON f.friend = p.user
		WHERE f.user = ? AND p.is_celebrity`

	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("GetCelebrityFriends: %w", err)
	}
	defer rows.Close()

	friends := make([]string, 0)
	for rows.Next() {
		var friend string
		if err = rows.Scan(&friend); err != nil {
			return nil, fmt.Errorf("GetCelebrityFriends: %w", err)
		}

		friends = append(friends, friend)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCelebrityFriends: %w", err)
	}

	return friends, nil
}

func (r *Repository) GetFollowersOfUser(ctx context.Context, username string) ([]string, error) {
	followers := make([]string, 0)

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT user FROM friends WHERE friend = ?`,
		username,
	)
	if err != nil {
		return nil, fmt.Errorf("GetFollowersOfUser: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var follower string
		if err = rows.Scan(&follower); err != nil {
			return nil, fmt.Errorf("GetFollowersOfUser: %w", err)
		}

		followers = append(followers, follower)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFollowersOfUser: %w", err)
	}

	return followers, nil
}
