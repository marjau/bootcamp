package controller

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/marcos-wz/capstone-go-bootcamp/internal/repository"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/service"

	"github.com/go-chi/render"
)

const (
	repoCsvErrType     errType = "RepositoryCSVError"
	repoDataApiErrType errType = "RepositoryDataAPIError"
	repoWPErrType      errType = "RepositoryWorkerPoolError"
	svcFilterErrType   errType = "ServiceFilterError"
	svcArgsErrType     errType = "ServiceArgumentsError"
)

var _ fmt.Stringer = errType("")

// errHTTP represents the message in the http error responses.
type errHTTP struct {
	Code      int     `json:"code"`
	ErrorType errType `json:"status"`
	Message   string  `json:"message"`
}

// errType represents the error classification or status.
type errType string

func (e errType) String() string {
	return string(e)
}

// errJSON returns an error JSON response.
func errJSON(w http.ResponseWriter, r *http.Request, err error) {
	errHttp := newErrHTTP(err)
	render.Status(r, errHttp.Code)
	render.JSON(w, r, errHttp)
}

// newErrHTTP returns a new errHTTP instance.
// It sets the Code, ErrorType, and Message values based on the error type.
func newErrHTTP(err error) errHTTP {
	var (
		repoCsvErr     *repository.CsvErr
		repoDataApiErr *repository.DataApiErr
		svcFilterErr   *service.FilterErr
		svcArgsErr     *service.ArgsErr
	)

	switch {

	// ###########  REPOSITORY ERRORS ###########

	case errors.As(err, &repoCsvErr):
		return errHTTP{
			Code:      http.StatusInternalServerError,
			ErrorType: repoCsvErrType,
			Message:   err.Error(),
		}
	case errors.As(err, &repoDataApiErr):
		return errHTTP{
			Code:      http.StatusBadGateway,
			ErrorType: repoDataApiErrType,
			Message:   err.Error(),
		}
	case errors.Is(err, repository.ErrWPInvalidArgs):
		return errHTTP{
			Code:      http.StatusInternalServerError,
			ErrorType: repoWPErrType,
			Message:   err.Error(),
		}

	// ########### SERVICE ERRORS ###########

	case errors.As(err, &svcFilterErr):
		return errHTTP{
			Code:      http.StatusUnprocessableEntity,
			ErrorType: svcFilterErrType,
			Message:   err.Error(),
		}
	case errors.As(err, &svcArgsErr):
		return errHTTP{
			Code:      http.StatusUnprocessableEntity,
			ErrorType: svcArgsErrType,
			Message:   err.Error(),
		}

	// ########### CONTROLLER ERRORS ###########

	// ########### DEFAULT ERRORS ###########

	default:
		return errHTTP{
			Code:      http.StatusBadRequest,
			ErrorType: errType(reflect.TypeOf(err).String()),
			Message:   err.Error(),
		}
	}
}
