package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/olgoncharov/otbook/internal/entity"
	dto "github.com/olgoncharov/otbook/internal/repository/dto"
	repoErrors "github.com/olgoncharov/otbook/internal/repository/errors"
)

func (r *Repository) CreateProfile(ctx context.Context, regInfo dto.RegistrationInfo) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("CreateProfile: can't begin transaction: %w", err)
	}

	defer tx.Rollback()

	var mysqlErr *mysql.MySQLError

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO users (username, password)
		 VALUES (?, ?)`,
		regInfo.Username, regInfo.Password,
	)

	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == duplicateEntryErrorCode {
			return fmt.Errorf("CreateProfile: %w", repoErrors.ErrUniqueConstraintViolated)
		}

		return fmt.Errorf("CreateProfile: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO profiles (user, first_name, last_name, birthdate, city, sex, hobby)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		regInfo.Username, regInfo.FirstName, regInfo.LastName, regInfo.Birthdate, regInfo.City, regInfo.Sex, regInfo.Hobby,
	)

	if err != nil {
		return fmt.Errorf("CreateProfile: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CreateProfile: can't commit tx: %w", err)
	}

	return nil
}

func (r *Repository) UpdateProfile(ctx context.Context, profile entity.Profile) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE profiles
		 SET
		 	first_name = ?,
			last_name = ?,
			birthdate = ?,
			city = ?,
			sex = ?,
			hobby = ?
		 WHERE user = ?`,

		profile.FirstName,
		profile.LastName,
		profile.Birthdate,
		profile.City,
		profile.Sex,
		profile.Hobby,
		profile.Username,
	)

	if err != nil {
		return fmt.Errorf("UpdateProfile: can't execute query: %w", err)
	}

	return nil
}

func (r *Repository) GetProfileByUsername(ctx context.Context, username string) (*entity.Profile, error) {
	profile := entity.Profile{
		Username: username,
	}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT
			first_name,
			last_name,
			birthdate,
			city,
			sex,
			hobby
		FROM profiles
		WHERE user = ?`,
		username,
	).Scan(
		&profile.FirstName,
		&profile.LastName,
		&profile.Birthdate,
		&profile.City,
		&profile.Sex,
		&profile.Hobby,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repoErrors.ErrNoRowsFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetProfileByUsername: %w", err)
	}

	return &profile, nil
}

func (r *Repository) GetAllProfiles(ctx context.Context, limit, offset uint) ([]entity.Profile, error) {
	profiles := make([]entity.Profile, 0, limit)

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
		ORDER BY user`, limit, offset)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetAllProfiles: %w", err)
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
			return nil, fmt.Errorf("GetAllProfiles: %w", err)
		}

		profiles = append(profiles, profile)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllProfiles: %w", err)
	}

	return profiles, nil
}

func (r *Repository) GetProfilesCount(ctx context.Context) (uint, error) {
	var totalCount uint

	err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM profiles").Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("GetProfilesCount: %w", err)
	}

	return totalCount, nil
}

func (r *Repository) CheckUsersExistence(ctx context.Context, usernames ...string) (map[string]bool, error) {
	if len(usernames) == 0 {
		return nil, nil
	}

	result := make(map[string]bool)

	query := `SELECT username FROM users WHERE username IN (?` + strings.Repeat(", ?", len(usernames)-1) + `)`

	args := make([]any, len(usernames))
	for i, u := range usernames {
		args[i] = u
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("CheckUsersExistence: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var current string
		err = rows.Scan(&current)
		if err != nil {
			return nil, fmt.Errorf("CheckUsersExistence: %w", err)
		}

		result[current] = true
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("CheckUsersExistence: %w", err)
	}

	return result, nil
}

func (r *Repository) SearchProfiles(ctx context.Context, firstNamePrefix, lastNamePrefix string) ([]dto.ProfileShortInfo, error) {
	profiles := make([]dto.ProfileShortInfo, 0)

	rows, err := r.db.QueryContext(ctx,
		`SELECT
			user,
			first_name,
			last_name
		FROM profiles
		WHERE
			first_name LIKE ? AND
			last_name LIKE ?`,
		firstNamePrefix+"%",
		lastNamePrefix+"%",
	)

	if err != nil {
		return nil, fmt.Errorf("SearchProfiles: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var profile dto.ProfileShortInfo
		err = rows.Scan(
			&profile.Username,
			&profile.FirstName,
			&profile.LastName,
		)
		if err != nil {
			return nil, fmt.Errorf("SearchProfiles: %w", err)
		}

		profiles = append(profiles, profile)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("SearchProfiles: %w", err)
	}

	return profiles, nil
}
