package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	filter "test/internal/infrastructure/filters"
	"test/internal/models"
	book_errors "test/internal/models/errors"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type BookStorage struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func NewBookStorage(db *sqlx.DB, logger *zap.Logger) IBookStorage {
	return &BookStorage{
		logger: logger,
		db:     db,
	}
}

func (bs *BookStorage) CreateBook(book *models.Book) error {
	query := `
        INSERT INTO books (year, title, author_id) 
        VALUES ($1, $2, $3)
        RETURNING id, year, title, available`

	args := []any{book.Year, book.Title, book.Author.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := bs.db.QueryRowContext(ctx, query, args...).Scan(
		&book.ID,
		&book.Year,
		&book.Title,
		&book.Available)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "foreign key"):
			bs.logger.Error("foreign key constraint error", zap.Error(err))
			return book_errors.ErrAuthorNotFound
		default:
			bs.logger.Error("errors on inserting into books", zap.Error(err))
			return err
		}
	}

	query = `
        SELECT  books.id,  books.year, books.title, books.available, authors.id, authors.name
        FROM books
		INNER JOIN authors ON books.author_id = authors.id
        WHERE books.id = $1
	`
	args = []any{book.ID}

	err = bs.db.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.Year, &book.Title, &book.Available, &book.Author.ID, &book.Author.Name)
	if err != nil {
		bs.logger.Error("error on getting a book", zap.Error(err))
		return err
	}

	return nil
}

func (bs *BookStorage) CreateAuthor(author *models.Author) error {
	query := `
        INSERT INTO authors (name) 
        VALUES ($1)
        RETURNING id`

	args := []any{author.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := bs.db.QueryRowContext(ctx, query, args...).Scan(&author.ID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "foreign key"):
			bs.logger.Error("foreign ley constraint error", zap.Error(err))
			return err
		default:
			bs.logger.Error("error on inserting into authors", zap.Error(err))
			return err
		}
	}

	query = `
        SELECT  books.id, books.year, books.title, books.available, authors.id, authors.name
        FROM books
		INNER JOIN authors ON books.author_id = authors.id 
        WHERE authors.id = $1
	`
	args = []any{author.ID}

	rows, err := bs.db.QueryContext(ctx, query, args...)
	if err != nil {
		bs.logger.Error("some error", zap.Error(err))
		return err
	}

	defer rows.Close()
	for rows.Next() {
		book := models.Book{
			Author: &models.Author{},
		}
		err := rows.Scan(&book.ID, &book.Year, &book.Title, &book.Available, &book.Author.ID, &book.Author.Name)
		if err != nil {
			bs.logger.Error("error on creating book array", zap.Error(err))
			return err
		}
		author.Books = append(author.Books, book)
	}
	if err = rows.Err(); err != nil {
		bs.logger.Error("errors on iterating", zap.Error(err))
		return err
	}

	return nil
}

func (bs *BookStorage) ListUsers(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), users.*
        FROM users
        ORDER BY %s %s, id ASC
        LIMIT $1 OFFSET $2`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{filters.Limit(), filters.Offset()}

	rows, err := bs.db.QueryContext(ctx, query, args...)
	if err != nil {
		bs.logger.Error("error on getting users", zap.Error(err))
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
			bs.logger.Error(" error encountered during iteration", zap.Error(err))
			return nil, filter.Metadata{}, err
		}

		query := `
        SELECT  books.id, books.year,books.title,  books.available, authors.id, authors.name
        FROM books
		INNER JOIN authors ON books.author_id = authors.id 
		INNER JOIN rented ON books.id = rented.book_id
        WHERE rented.user_id = $1
	`
		args := []any{user.ID}

		rows, err := bs.db.QueryContext(ctx, query, args...)
		if err != nil {
			bs.logger.Error("some error", zap.Error(err))
			return nil, filter.Metadata{}, err
		}

		defer rows.Close()
		for rows.Next() {
			book := models.Book{
				Author: &models.Author{},
			}
			err := rows.Scan(&book.ID, &book.Year, &book.Title, &book.Available, &book.Author.ID, &book.Author.Name)
			if err != nil {
				bs.logger.Error("some error during iteration", zap.Error(err))
				return nil, filter.Metadata{}, err
			}
			user.RentedBooks = append(user.RentedBooks, book)
		}
		if err = rows.Err(); err != nil {
			bs.logger.Error("some error on closing rows", zap.Error(err))
			return nil, filter.Metadata{}, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		bs.logger.Error("some error on closing rows", zap.Error(err))
		return nil, filter.Metadata{}, err
	}

	metadata := filter.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return users, metadata, nil

}

func (bs *BookStorage) ListBooks(filters filter.Filters) ([]*models.Book, filter.Metadata, error) {

	query := fmt.Sprintf(`
        SELECT count(*) OVER(), books.*, authors.name
        FROM books
		INNER JOIN authors on books.author_id = authors.id
        ORDER BY %s %s, id ASC
        LIMIT $1 OFFSET $2`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{filters.Limit(), filters.Offset()}

	rows, err := bs.db.QueryContext(ctx, query, args...)
	if err != nil {
		bs.logger.Error("some error", zap.Error(err))
		return nil, filter.Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0

	books := []*models.Book{}

	for rows.Next() {
		book := models.Book{
			Author: &models.Author{},
		}

		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.Year,
			&book.Title,
			&book.Available,
			&book.Author.ID,
			&book.Author.Name,
		)
		if err != nil {
			bs.logger.Error("some error during iteration", zap.Error(err))
			return nil, filter.Metadata{}, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		bs.logger.Error("some error on closing the row", zap.Error(err))
		return nil, filter.Metadata{}, err
	}

	metadata := filter.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return books, metadata, nil
}

func (bs *BookStorage) ListAuthors(filters filter.Filters) ([]*models.Author, filter.Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), authors.id, authors.name, authors.times_ordered
        FROM authors
        ORDER BY %s %s, id ASC
        LIMIT $1 OFFSET $2`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{filters.Limit(), filters.Offset()}

	rows, err := bs.db.QueryContext(ctx, query, args...)
	if err != nil {
		bs.logger.Error("some error on getting authors", zap.Error(err))
		return nil, filter.Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0

	authors := []*models.Author{}

	for rows.Next() {
		var author models.Author

		err := rows.Scan(
			&totalRecords,
			&author.ID,
			&author.Name,
			&author.Times_ordered,
		)
		if err != nil {
			bs.logger.Error("some error", zap.Error(err))
			return nil, filter.Metadata{}, err
		}

		query := `
        SELECT  books.id, books.year, books.title, books.available, authors.id, authors.name
        FROM books
		INNER JOIN authors ON books.author_id = authors.id 
        WHERE authors.id = $1
	`
		args := []any{author.ID}

		rows, err := bs.db.QueryContext(ctx, query, args...)
		if err != nil {
			bs.logger.Error("some error", zap.Error(err))
			return nil, filter.Metadata{}, err
		}

		defer rows.Close()
		for rows.Next() {
			book := models.Book{
				Author: &models.Author{},
			}
			err := rows.Scan(&book.ID, &book.Year, &book.Title, &book.Available, &book.Author.ID, &book.Author.Name)
			if err != nil {
				bs.logger.Error("some error", zap.Error(err))
				return nil, filter.Metadata{}, err
			}
			author.Books = append(author.Books, book)
		}
		if err = rows.Err(); err != nil {
			bs.logger.Error("some error", zap.Error(err))
			return nil, filter.Metadata{}, err
		}

		authors = append(authors, &author)
	}

	if err = rows.Err(); err != nil {
		bs.logger.Error("some error", zap.Error(err))
		return nil, filter.Metadata{}, err
	}

	metadata := filter.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return authors, metadata, nil
}

func (bs *BookStorage) RentBook(bookID, userID int64) error {

	query := `
	INSERT INTO rented  (book_id, user_id) 
	VALUES ($1, $2)
	RETURNING id`

	args := []any{bookID, userID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	err := bs.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "foreign key"):
			bs.logger.Error("attempting to rent a book that is not available", zap.Error(err))
			return book_errors.ErrRentInvalid
		default:
			bs.logger.Error("some error", zap.Error(err))
			return err
		}
	}

	query = `
        UPDATE books
        SET available = false
        WHERE id = $1 AND available = true
		RETURNING author_id
       `

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	book := models.Book{
		Author: &models.Author{},
	}
	err = bs.db.QueryRowContext(ctx, query, bookID).Scan(&book.Author.ID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			bs.logger.Error("foreign key violation error", zap.Error(err))
			return err
		case errors.Is(err, sql.ErrNoRows):
			bs.logger.Error("no rows", zap.Error(err))
			return book_errors.ErrRentInvalid
		default:
			bs.logger.Error("some error", zap.Error(err))
			return err
		}
	}
	query = `
	UPDATE authors
	SET times_ordered = times_ordered + 1
	WHERE id = $1
	`
	err = bs.db.QueryRowContext(ctx, query, book.Author.ID).Err()
	if err != nil {
		bs.logger.Error("some error on updating author table", zap.Error(err))
		return err
	}

	return nil
}

func (bs *BookStorage) ReturnBook(bookID, userID int64) error {

	query := `
	DELETE FROM rented
        WHERE book_id = $1 AND user_id = $2`

	args := []any{bookID, userID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := bs.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return book_errors.ErrBookNotFound
	}

	query = `
        UPDATE books 
        SET available = true
        WHERE id = $1 AND available = false
		RETURNING id, year, title, available, author_id
       `

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	book := models.Book{
		Author: &models.Author{},
	}
	err = bs.db.QueryRowContext(ctx, query, bookID).Scan(&book.ID, &book.Year, &book.Title, &book.Available, &book.Author.ID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			bs.logger.Error("some error", zap.Error(err))
			return err
		case errors.Is(err, sql.ErrNoRows):
			bs.logger.Error("some error", zap.Error(err))
			return book_errors.ErrBookNotFound
		default:
			bs.logger.Error("some error", zap.Error(err))
			return err
		}
	}

	return nil
}
