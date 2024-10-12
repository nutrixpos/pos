package customerrors

import "errors"

var (
	ErrModuleNotRegistered     = errors.New("can't register module")
	ErrModuleNameAlreadyExists = errors.New("module name already exists")
	ErrSeedingModule           = errors.New("error seeding module")
	ErrTypeAssersionFailed     = errors.New("type assertion failed")
	ErrInvalidObjectId         = errors.New("invalid object id")
)
