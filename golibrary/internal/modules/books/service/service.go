package service

import (
	"sort"
	filter "test/internal/infrastructure/filters"
	"test/internal/models"
	"test/internal/modules/books/repository"
)

type IBookService interface {
	CreateBook(book *models.Book) error
	CreateAuthor(author *models.Author) error
	ListUsers(filters filter.Filters) ([]*models.User, filter.Metadata, error)
	ListBooks(filters filter.Filters) ([]*models.Book, filter.Metadata, error)
	ListAuthors(filters filter.Filters) ([]*models.Author, filter.Metadata, error)
	ListTopRatedAuthors(filters filter.Filters) ([]*models.Author, filter.Metadata, error)
	RentBook(userID, bookID int64) error
	ReturnBook(userID, bookID int64) error
}

type BookService struct {
	storage repository.IBookStorage
}

func NewBookService(repo repository.IBookStorage) *BookService {
	return &BookService{storage: repo}
}

func (s *BookService) CreateBook(book *models.Book) error {
	return s.storage.CreateBook(book)
}

func (s *BookService) CreateAuthor(author *models.Author) error {
	return s.storage.CreateAuthor(author)
}

func (s *BookService) ListUsers(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
	return s.storage.ListUsers(filters)
}

func (s *BookService) ListBooks(filters filter.Filters) ([]*models.Book, filter.Metadata, error) {
	return s.storage.ListBooks(filters)
}

func (s *BookService) ListAuthors(filters filter.Filters) ([]*models.Author, filter.Metadata, error) {
	return s.storage.ListAuthors(filters)
}

func (s *BookService) ListTopRatedAuthors(filters filter.Filters) ([]*models.Author, filter.Metadata, error) {
	authors, meta, err := s.storage.ListAuthors(filters)

	if err != nil {
		return nil, filter.Metadata{}, err
	}
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Times_ordered > authors[j].Times_ordered
	})
	return authors, meta, nil
}
func (s *BookService) RentBook(bookID, userID int64) error {
	return s.storage.RentBook(bookID, userID)
}

func (s *BookService) ReturnBook(bookID, userID int64) error {
	return s.storage.ReturnBook(bookID, userID)
}
