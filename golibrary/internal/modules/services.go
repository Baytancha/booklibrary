package modules

import (
	"test/internal/infrastructure/components"
	book_service "test/internal/modules/books/service"
	user_service "test/internal/modules/user/service"
)

type Services struct {
	UserService user_service.IUserService
	BookService book_service.IBookService
}

func NewServices(cmp *components.Components, storages *Storages) *Services {
	return &Services{
		UserService: user_service.NewUserService(storages.UserStorage),
		BookService: book_service.NewBookService(storages.BookStorage),
	}
}
