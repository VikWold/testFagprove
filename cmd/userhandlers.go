package main

import (
	"errors"
	"net/http"
	"testFagprove/internal/auth"
	"testFagprove/internal/data"
	"testFagprove/internal/loggingutils"
	"testFagprove/internal/rest"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	CreatedAt   *time.Time `json:"created_at"`
	LastUpdated *time.Time `json:"last_updated"`
}

type UserListresponse struct {
	Users []*data.User `json:"users"`
}

func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := loggingutils.LoggerFromContext(ctx)

	users, err := app.models.User.List(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "unable to get users", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.InfoContext(ctx, "returning users")

	rest.WriteJSON(
		w,
		http.StatusOK,
		UserListresponse{Users: users},
		nil,
	)
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		rest.BadRequestResponse(w, r, "unable to parse id from path")
		return
	}

	user, err := app.models.User.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			rest.NotFoundResponse(w, r)
		default:
			logger.ErrorContext(ctx, "unable to get user", "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.InfoContext(ctx, "returning user")

	rest.WriteJSON(
		w,
		http.StatusOK,
		UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			CreatedAt:   user.CreatedAt,
			LastUpdated: user.LastUpdated,
		},
		nil,
	)
}

func (app *application) createuserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	var req data.User
	err := rest.ReadJSON(r, &req)
	if err != nil {
		rest.BadRequestResponse(w, r, "unable to decode data from request")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.ErrorContext(ctx, "unable to hash password", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	user := &data.User{
		ID:          req.ID,
		Username:    req.Username,
		Password:    hashedPassword,
		CreatedAt:   req.CreatedAt,
		LastUpdated: req.LastUpdated,
	}

	result, err := app.models.User.Insert(ctx, user)
	if err != nil {
		logger.ErrorContext(ctx, "unable to insert user", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.InfoContext(ctx, "user created")

	rest.WriteJSON(
		w,
		http.StatusCreated,
		UserResponse{
			ID:       result.ID,
			Username: result.Username,
		},
		nil,
	)
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		rest.BadRequestResponse(w, r, "unable to parse id from path")
		return
	}

	err = app.models.User.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			rest.NotFoundResponse(w, r)
		default:
			logger.ErrorContext(ctx, "unable to delete user", "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.InfoContext(ctx, "user deleted")
	w.WriteHeader(http.StatusNoContent)
}
