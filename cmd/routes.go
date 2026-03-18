package main

import "net/http"

func (app *application) routes() http.Handler {
	app.logger.Info("creating new multiplexer")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/ehf", app.listEhfHandler)
	mux.HandleFunc("GET /api/v1/ehf/{id}", app.getEhfHandler)
	mux.HandleFunc("POST /api/v1/ehf", app.createEhfHandler)
	mux.HandleFunc("DELETE /api/v1/ehf/{id}", app.deleteEhfHandler)

	mux.HandleFunc("GET /api/v1/user", app.listUsersHandler)
	mux.HandleFunc("GET /api/v1/user/{id}", app.getUserHandler)
	mux.HandleFunc("POST /api/v1/user", app.createUserHandler)
	mux.HandleFunc("DELETE /api/v1/user/{id}", app.deleteUserHandler)

	return app.enableCORS(mux)
}
