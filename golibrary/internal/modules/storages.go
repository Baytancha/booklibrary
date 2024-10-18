package modules

import (
	book_storage "test/internal/modules/books/repository"
	user_storage "test/internal/modules/user/repository"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Storages struct {
	UserStorage user_storage.IUserStorage
	BookStorage book_storage.IBookStorage
}

func NewStorages(sql *sqlx.DB, logger *zap.Logger) *Storages {
	return &Storages{
		UserStorage: user_storage.NewUserModel(sql),
		BookStorage: book_storage.NewBookStorage(sql, logger),
	}
}
