package consts

import "fmt"

const (
	NoRowsFound         = Error("no rows found")
	ForeignKeyViolation = Error("foreign key violation")
	UniqueViolation     = Error("unique violation")
	UndefinedTable      = Error("undefined table")
	NullValueNotAllowed = Error("null value not allowed")
)

// Error to detect response status code 400
type Error string

func (e Error) Error() string {
	return string(e)
}

// ErrNotFound to detect response status code 404
type ErrNotFound string

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s Not Found", string(e))
}

// ErrAlreadyExist to detect response status code 409
type ErrAlreadyExist string

func (e ErrAlreadyExist) Error() string {
	return fmt.Sprintf("%s Already Exist", string(e))
}

// ErrInvalidMetaData for error invalid meta data
func ErrInvalidMetaData(additional string) Error {
	return Error(fmt.Sprintf("Invalid Meta Data %s", additional))
}

const (
	ErrInvalidJSON = Error("Invalid JSON Request")
)
