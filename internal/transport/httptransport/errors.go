package httptransport

import "github.com/bboykiv/topsigner/gen/httpserver"

func NewInternalError() httpserver.InternalErrorJSONResponse {
	return httpserver.InternalErrorJSONResponse{Message: "internal server error"}
}

func NewUnauthorizedError() httpserver.UnauthorizedJSONResponse {
	return httpserver.UnauthorizedJSONResponse{Message: "unauthorized"}
}

func NewBadRequestError(message string) httpserver.BadRequestJSONResponse {
	return httpserver.BadRequestJSONResponse{Message: message}
}
