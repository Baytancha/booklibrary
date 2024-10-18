package router

import (
	"net/http"

	"github.com/go-chi/chi"

	_ "github.com/mattn/go-sqlite3"

	"test/internal/infrastructure/components"
	"test/internal/modules"
	swagger "test/static"

	"test/internal/infrastructure/helpers"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func Routes(ctrl *modules.Controllers, comp *components.Components) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Use(jwtauth.Verifier(helpers.TokenAuth))
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token, _, err := jwtauth.FromContext(r.Context())

				if err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				if token == nil {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				} else if jwt.Validate(token) != nil {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}

				next.ServeHTTP(w, r)
			})
		})

	})

	r.Post("/books/book", ctrl.BookHandler.CreateBook)
	r.Post("/books/author", ctrl.BookHandler.CreateAuthor)

	r.Patch("/books/rent/{bookID}", ctrl.BookHandler.RentBook)     //	creates new record in junction rable
	r.Patch("/books/return/{bookID}", ctrl.BookHandler.ReturnBook) // erases record in junction rable

	r.Get("/books/listUsers", ctrl.BookHandler.ListUsers)     //list users with rented books
	r.Get("/books/listBooks", ctrl.BookHandler.ListBooks)     //list books with authors
	r.Get("/books/listAuthors", ctrl.BookHandler.ListAuthors) //list authors with books

	r.Get("/books/rate", ctrl.BookHandler.ListTopRatedAuthors)

	r.Post("/user", ctrl.UserHandler.CreateUser)
	r.Get("/user/login", ctrl.UserHandler.Login)
	r.Get("/user/logout", ctrl.UserHandler.Logout)
	r.Get("/user/{username}", ctrl.UserHandler.GetUserByName)
	r.Put("/user/{username}", ctrl.UserHandler.UpdateUser)
	r.Delete("/user/{username}", ctrl.UserHandler.DeleteUser)
	r.Post("/user/CreateWithList", ctrl.UserHandler.CreateWithList)
	r.Post("/user/CreateWithArray", ctrl.UserHandler.CreateWithArray)

	fileServer := http.FileServerFS(swagger.Swaggerfile)
	r.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/swagger", fileServer)
		fs.ServeHTTP(w, r)
	})

	return r
}
