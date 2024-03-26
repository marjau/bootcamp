package service

import (
	"errors"
	"fmt"
)

var (
	ErrFltrTypeEmpty  = errors.New("filter type empty")
	ErrFltrValueEmpty = errors.New("filter value empty")
	ErrFltrInvalid    = errors.New("invalid filter")

	ErrInvalidNumType   = errors.New("invalid number type")
	ErrZeroValue        = errors.New("zero value is not allowed")
	ErrJobsWorkerHigher = errors.New("jobs per worker higher than maximum jobs")
)

// FilterErr covers all errors related to Filters and wraps the error that caused it.
type FilterErr struct {
	Err error
}

func (e FilterErr) Error() string {
	return fmt.Sprintf("service filter: %s", e.Err)
}

func (e FilterErr) Unwrap() error {
	return e.Err
}

// ArgsErr covers all errors related to the given arguments and wraps the error that caused it.
type ArgsErr struct {
	Err error
}

func (e ArgsErr) Error() string {
	return fmt.Sprintf("service arguments: %s", e.Err)
}

func (e ArgsErr) Unwrap() error {
	return e.Err
}
