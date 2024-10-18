package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	//"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	filter "test/internal/infrastructure/filters"
	"test/internal/models"
)

type UserModel struct {
	DB *sqlx.DB
}

func NewUserModel(db *sqlx.DB) IUserStorage {
	return &UserModel{DB: db}
}

var (
	ErrRecordNotFound = errors.New("user not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

func (m UserModel) Get(id int64) (*models.User, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT  id,  name, email, password_hash, deleted, version
        FROM users
        WHERE id = $1 AND deleted = false`

	user := &models.User{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(

		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Deleted,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("no rows", err)
			return nil, ErrRecordNotFound
		default:
			fmt.Println("some error", err)
			return nil, err
		}
	}

	return user, nil
}

func (m UserModel) Insert(user *models.User) error {
	query := `
        INSERT INTO users (name, email, password_hash, deleted) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, version`

	args := []any{user.Name, user.Email, user.Password.Hash, user.Deleted}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			return ErrDuplicateEmail
		default:
			fmt.Println("some error", err)
			return err
		}
	}

	return nil
}

func (m UserModel) GetByName(username string) (*models.User, error) {
	query := `
        SELECT id, name, email, password_hash, deleted, version
        FROM users
        WHERE name = $1 AND deleted = false`

	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Deleted,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("no rows", err)
			return nil, ErrRecordNotFound
		default:
			fmt.Println("some error", err)
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *models.User) error {
	query := `
        UPDATE users 
        SET name = $1, email = $2, password_hash = $3, version = version + 1
        WHERE id = $4 AND version = $5 AND deleted = false
        RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			fmt.Println("duplicate", err)
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("no rows", err)
			return ErrEditConflict
		default:
			fmt.Println("some error", err)
			return err
		}
	}

	return nil
}

func (m UserModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        UPDATE users 
        SET deleted = true, version = version + 1
		WHERE id = $1 AND deleted = false
		RETURNING version
       `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &models.User{}

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			fmt.Println("duplicate", err)
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("no rows", err)
			return ErrEditConflict
		default:
			fmt.Println("some error", err)
			return err
		}
	}

	return nil

}

func (u UserModel) GetAll(filters filter.Filters) ([]*models.User, filter.Metadata, error) {

	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id,  name, email, password_hash, deleted, version
        FROM users  
		WHERE deleted = false
        ORDER BY %s %s, id ASC
        LIMIT $1 OFFSET $2`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{filters.Limit(), filters.Offset()}

	rows, err := u.DB.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Println("some query error", err)
		return nil, filter.Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0

	users := []*models.User{}

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password.Hash,
			&user.Deleted,
			&user.Version,
		)
		if err != nil {
			fmt.Println("some row error", err)
			return nil, filter.Metadata{}, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("some iter error", err)
		return nil, filter.Metadata{}, err
	}

	metadata := filter.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return users, metadata, nil

}
