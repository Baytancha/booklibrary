package modules

import (
	"test/internal/infrastructure/components"
	book_controller "test/internal/modules/books/controller"
	user_controller "test/internal/modules/user/controller"
)

type Controllers struct {
	UserHandler user_controller.IUserHandler
	BookHandler book_controller.IBookController
}

func NewControllers(services *Services, components *components.Components) *Controllers {
	return &Controllers{
		UserHandler: user_controller.NewUserHandler(components.Responder, services.UserService),
		BookHandler: book_controller.NewBookController(components.Responder, services.BookService),
	}
}
