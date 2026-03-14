package main

import "net/http"

func (app *application) routes() http.Handler {
	app.logger.Info("creating new multiplexer")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/ehf", app.listEhfHandler)
	mux.HandleFunc("GET /api/v1/ehf/{id}", app.getEhfHandler)
	mux.HandleFunc("POST /api/v1/ehf", app.createEhfHandler)
	mux.HandleFunc("DELETE /api/v1/ehf/{id}", app.deleteEhfHandler)

	return app.enableCORS(mux)
}
