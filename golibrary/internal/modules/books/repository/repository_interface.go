package repository

import (
	filter "test/internal/infrastructure/filters"
	"test/internal/models"
)

type IBookStorage interface {
	CreateBook(book *models.Book) error
	CreateAuthor(author *models.Author) error
	ListUsers(filters filter.Filters) ([]*models.User, filter.Metadata, error)
	ListBooks(filters filter.Filters) ([]*models.Book, filter.Metadata, error)
	ListAuthors(filters filter.Filters) ([]*models.Author, filter.Metadata, error)
	RentBook(userID, bookID int64) error
	ReturnBook(userID, bookID int64) error
}
