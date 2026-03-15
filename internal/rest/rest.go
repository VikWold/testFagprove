package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"testFagprove/internal/loggingutils"
)

type ErrorMessage struct {
	Message any `json:"message"`
}

func ReadJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header,
) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		return err
	}

	return nil
}

func LogError(r *http.Request, err error) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	logger.Error(
		"an error occurred",
		"request_method", r.Method,
		"request_url", r.URL.String(),
		"error", err,
	)
}

func ParameterNotFoundResponse(
	w http.ResponseWriter,
	r *http.Request,
	param string,
) {
	logger := loggingutils.LoggerFromContext(r.Context())

	logger.Info(
		"parameter for requested resource could not be found", slog.String("parameter", param),
	)
	ErrorResponse(w, r, http.StatusNotFound, fmt.Sprintf("%s not found in path", param))
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	logger := loggingutils.LoggerFromContext(r.Context())

	const notFoundMsg string = "the requested resource could not be found"

	logger.Info(notFoundMsg)
	ErrorResponse(w, r, http.StatusNotFound, notFoundMsg)
}

func ErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message any,
) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	logger.Info("writing error response", slog.Int("status", status), slog.Any("message", message))
	err := WriteJSON(w, status, ErrorMessage{Message: message}, nil)
	if err != nil {
		logger.Error("error writing response", "error", err)
		LogError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger := loggingutils.LoggerFromContext(r.Context())

	LogError(r, err)
	const serverErrorMsg string = "the server encountered a problem and could not process your request"

	logger.Info(serverErrorMsg)
	ErrorResponse(w, r, http.StatusInternalServerError, serverErrorMsg)
}
func BadRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	logger := loggingutils.LoggerFromContext(r.Context())

	logger.Info("bad request response", slog.String("message", message))
	ErrorResponse(w, r, http.StatusBadRequest, message)
}

func RespondWithJSON(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data any,
	headers http.Header,
) {
	logger := loggingutils.LoggerFromContext(r.Context())

	logger.Info("marshalling data")
	js, err := json.Marshal(data)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}

	js = append(js, '\n')

	logger.Info("adding headers")
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	logger.Info("writing response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		ServerErrorResponse(w, r, err)
	}
}
