package handler

import (
	"encoding/json"
	"net/http"
)

var (
	internalServerErrorResponse = newJsonResponse(
		http.StatusInternalServerError,
		newHandlerError(http.StatusInternalServerError, ErrInternalServer.Error()),
	)
	unauthorizedResponse = newJsonResponse(
		http.StatusUnauthorized,
		newHandlerError(http.StatusUnauthorized, ErrUnauthorized.Error()),
	)
	badRequestResponse = newJsonResponse(
		http.StatusBadRequest,
		newHandlerError(http.StatusBadRequest, ErrBadRequest.Error()),
	)
	unsupportedMediaTypeResponse = newJsonResponse(
		http.StatusUnsupportedMediaType,
		newHandlerError(http.StatusUnsupportedMediaType, ErrUnsupportedMediaType.Error()),
	)
	invalidCredentialsResponse = newJsonResponse(
		http.StatusBadRequest,
		newHandlerError(http.StatusBadRequest, ErrInvalidCredentials.Error()),
	)
	usernameAlreadyTakenResponse = newJsonResponse(
		http.StatusConflict,
		newHandlerError(http.StatusConflict, ErrUsernameAlreadyTaken.Error()),
	)
)

type response interface {
	write(w http.ResponseWriter)
}

type jsonResponse struct {
	statusCode int
	body       any
}

func newJsonResponse(statusCode int, body any) *jsonResponse {
	return &jsonResponse{statusCode: statusCode, body: body}
}

func (r *jsonResponse) write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.statusCode)
	err := json.NewEncoder(w).Encode(r.body)
	if err != nil {
		// It's not possible to change the status code header
		// as it can only be sent once.
		// TODO: Somehow send the correct status code header
		_ = json.NewEncoder(w).Encode(newHandlerError(
			http.StatusInternalServerError,
			ErrInternalServer.Error(),
		))
	}
}

type noContentResponse struct{}

func (r noContentResponse) write(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

type createdResponse struct{}

func (r createdResponse) write(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}
