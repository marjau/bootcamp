package repository

import (
	"errors"
	"fmt"
)

var (
	ErrFileNameEmpty   = errors.New("file name empty")
	ErrFilePathIsDir   = errors.New("file path is a directory")
	ErrFileModeInvalid = errors.New("invalid file permissions")
	ErrDirNameEmpty    = errors.New("directory name empty")
	ErrIsNotDir        = errors.New("is not a dataDir")

	ErrURLPathEmpty    = errors.New("url path empty")
	ErrInvalidRespCode = errors.New("invalid response code")

	ErrCSVRecEmpty  = errors.New("csv record empty")
	ErrJsonRecEmpty = errors.New("json record empty")

	ErrCocktailNameEmpty         = errors.New("cocktail name empty")
	ErrCocktailInstructionsEmpty = errors.New("cocktail instructions empty")
	ErrCocktailIngredientsEmpty  = errors.New("cocktail ingredients empty")

	ErrWPInvalidArgs = errors.New("worker pool: invalid arguments")
)

// CsvErr covers all errors related to CSV operations and wraps the error that caused it.
type CsvErr struct {
	Err error
}

func (e CsvErr) Error() string {
	return fmt.Sprintf("csv: %s", e.Err)
}

func (e CsvErr) Unwrap() error {
	return e.Err
}

// DataApiErr covers all errors related to data api operations and wraps the error that caused it.
type DataApiErr struct {
	Err error
}

func (e DataApiErr) Error() string {
	return fmt.Sprintf("data api: %s", e.Err)
}

func (e DataApiErr) Unwrap() error {
	return e.Err
}
