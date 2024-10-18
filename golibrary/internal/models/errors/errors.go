package errors

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrBookNotFound   = errors.New("book not found")
	ErrAuthorNotFound = errors.New("author not found")
	ErrRentInvalid    = errors.New("this book cannot be rented")
	ErrDuplicateEmail = errors.New("duplicate email")
)
