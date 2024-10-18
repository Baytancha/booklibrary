package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"test/internal/infrastructure/filters"
	"test/internal/infrastructure/helpers"
	"test/internal/infrastructure/responder"
	"test/internal/infrastructure/validator"
	"test/internal/models"
	book_error "test/internal/models/errors"
	"test/internal/modules/books/service"

	"github.com/go-chi/chi"
)

type IBookController interface {
	CreateBook(w http.ResponseWriter, r *http.Request)
	CreateAuthor(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	ListBooks(w http.ResponseWriter, r *http.Request)
	ListAuthors(w http.ResponseWriter, r *http.Request)
	ListTopRatedAuthors(w http.ResponseWriter, r *http.Request)
	RentBook(w http.ResponseWriter, r *http.Request)
	ReturnBook(w http.ResponseWriter, r *http.Request)
}

type BookController struct {
	responder responder.Responder
	service   service.IBookService
}

func NewBookController(responder responder.Responder, service service.IBookService) *BookController {
	return &BookController{
		responder: responder,
		service:   service,
	}
}

func (bc *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {

	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		bc.responder.ErrorBadRequest(w, err)
		return
	}
	fmt.Println(book)
	err = bc.service.CreateBook(&book)
	if err != nil {
		switch {
		case errors.Is(err, book_error.ErrAuthorNotFound):
			bc.responder.ErrorInternal(w, errors.New("Author not found"))
		default:
			bc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}
	bc.responder.OutputJSON(w, book)
}

func (bc *BookController) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author *models.Author
	err := json.NewDecoder(r.Body).Decode(&author)

	if err != nil {
		bc.responder.ErrorBadRequest(w, err)
		return
	}

	err = bc.service.CreateAuthor(author)
	if err != nil {
		bc.responder.ErrorInternal(w, err)
		return
	}

	bc.responder.OutputJSON(w, author)
}

func (bc *BookController) ListUsers(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string
		Email string
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "page_size", 20, v)

	input.Filters.Sort = helpers.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "email", "-id", "-name", "-email"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		bc.responder.ErrorInternal(w, errors.New("Internal server error1"))
		return
	}

	users, metadata, err := bc.service.ListUsers(input.Filters)
	if err != nil {
		bc.responder.ErrorInternal(w, errors.New("Internal server error2"))
		return
	}
	bc.responder.OutputJSON(w, map[string]interface{}{"metadata": metadata, "data": users})
}

func (bc *BookController) ListBooks(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string
		Email string
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "page_size", 20, v)

	input.Filters.Sort = helpers.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "email", "-id", "-name", "-email"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		bc.responder.ErrorInternal(w, errors.New("Internal server error1"))
		return
	}

	books, metadata, err := bc.service.ListBooks(input.Filters)
	if err != nil {
		bc.responder.ErrorInternal(w, errors.New("Internal server error2"))
		return
	}
	bc.responder.OutputJSON(w, map[string]interface{}{"metadata": metadata, "data": books})
}

func (bc *BookController) ListAuthors(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string
		Email string
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "page_size", 20, v)

	input.Filters.Sort = helpers.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "email", "-id", "-name", "-email"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		bc.responder.ErrorInternal(w, errors.New("Internal server error1"))
		return
	}

	authors, metadata, err := bc.service.ListAuthors(input.Filters)
	if err != nil {
		bc.responder.ErrorInternal(w, errors.New("Internal server error2"))
		return
	}
	bc.responder.OutputJSON(w, map[string]interface{}{"metadata": metadata, "data": authors})
}

func (bc *BookController) ListTopRatedAuthors(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string
		Email string
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "page_size", 20, v)

	input.Filters.Sort = helpers.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "email", "-id", "-name", "-email"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		bc.responder.ErrorInternal(w, errors.New("Internal server error1"))
		return
	}

	authors, metadata, err := bc.service.ListTopRatedAuthors(input.Filters)
	if err != nil {
		bc.responder.ErrorInternal(w, errors.New("Internal server error2"))
		return
	}
	bc.responder.OutputJSON(w, map[string]interface{}{"metadata": metadata, "data": authors})
}

func (bc *BookController) RentBook(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "bookID")

	bookID, err := strconv.Atoi(rawID)
	if err != nil {
		bc.responder.ErrorBadRequest(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		bc.responder.ErrorBadRequest(w, err)
		return
	}

	rawID = r.FormValue("userID")

	userID, err := strconv.Atoi(rawID)
	if err != nil {
		bc.responder.ErrorBadRequest(w, err)
		return
	}

	err = bc.service.RentBook(int64(bookID), int64(userID))
	if err != nil {
		bc.responder.ErrorInternal(w, err)
		return
	}

	bc.responder.OutputJSON(w, "Book rented successfully")
}

func (bc *BookController) ReturnBook(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "bookID")

	bookID, err := strconv.Atoi(rawID)
	if err != nil {
		bc.responder.ErrorBadRequest(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		bc.responder.ErrorBadRequest(w, err)
		return
	}

	rawID = r.FormValue("userID")

	userID, err := strconv.Atoi(rawID)
	if err != nil {
		bc.responder.ErrorBadRequest(w, err)
		return
	}

	err = bc.service.ReturnBook(int64(bookID), int64(userID))
	if err != nil {
		bc.responder.ErrorInternal(w, err)
		return
	}

	bc.responder.OutputJSON(w, "Book returned successfully")

}
