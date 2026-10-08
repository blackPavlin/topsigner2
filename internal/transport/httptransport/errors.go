package httptransport

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

var (
	ErrInvalidMultipartForm    = errors.New("invalid multipart form")
	ErrMultipartFileIsRequired = errors.New("multipart file is required")
	ErrInvalidOAuthProvider    = errors.New("invalid oauth provider")
)

func NewInternalError() httpserver.InternalErrorJSONResponse {
	return httpserver.InternalErrorJSONResponse{Message: "internal server error"}
}

func NewUnauthorizedError() httpserver.UnauthorizedJSONResponse {
	return httpserver.UnauthorizedJSONResponse{Message: "unauthorized"}
}

func NewBadRequestError(err error) httpserver.BadRequestJSONResponse {
	message := err.Error()

	if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
		message = fmt.Sprintf(
			"field validation for '%s' failed on the '%s' tag",
			validationErrors[0].Field(),
			validationErrors[0].Tag(),
		)
	}

	return httpserver.BadRequestJSONResponse{Message: message}
}

func NewNotFoundError(err error) httpserver.NotFoundJSONResponse {
	return httpserver.NotFoundJSONResponse{Message: err.Error()}
}

func NewConflictError(err error) httpserver.ConflictJSONResponse {
	return httpserver.ConflictJSONResponse{Message: err.Error()}
}

func requestErrorHandlerFunc(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.Respond(w, r, NewBadRequestError(err))
}

func responseErrorHandlerFunc(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)
	render.Respond(w, r, NewInternalError())
}
