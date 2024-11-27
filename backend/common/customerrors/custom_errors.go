// Package customerrors contains custom error types and values.
//
// Error types are used by functions to indicate that an error has occurred.
// The error type is a string that describes the error.
//
// Error values are used by functions to indicate that an error has occurred.
// The error value is an instance of the error type.
//
// The variables declared in this file are error values.
package customerrors

import "errors"

// ErrModuleNotRegistered is an error returned when a module can't be registered.
var ErrModuleNotRegistered = errors.New("can't register module")

// ErrModuleNameAlreadyExists is an error returned when a module name already exists.
var ErrModuleNameAlreadyExists = errors.New("module name already exists")

// ErrSeedingModule is an error returned when an error occurs while seeding a module.
var ErrSeedingModule = errors.New("error seeding module")

// ErrTypeAssersionFailed is an error returned when a type assertion fails.
var ErrTypeAssersionFailed = errors.New("type assertion failed")

// ErrInvalidObjectId is an error returned when an object id is invalid.
var ErrInvalidObjectId = errors.New("invalid object id")

// ErrInsufficientReady is an error returned when there is not enough ready.
var ErrInsufficientReady = errors.New("insufficient ready")
