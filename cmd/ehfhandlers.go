package main

import (
	"errors"
	"net/http"
	"testFagprove/internal/data"
	"testFagprove/internal/loggingutils"
	"testFagprove/internal/rest"
	"time"

	"github.com/google/uuid"
)

type EhfResponse struct {
	EhfID          uuid.UUID  `json:"ehf_id"`
	FileName       string     `json:"file_name"`
	CustomerID     int        `json:"customer_id"`
	SupplierID     int        `json:"supplier_id"`
	InvoiceNo      string     `json:"invoice_number"`
	BuyerReference string     `json:"buyer_reference"`
	IssueDate      *time.Time `json:"issue_date"`
	DueDate        *time.Time `json:"due_date"`
	Currency       string     `json:"currency"`
	Amount         float64    `json:"amount"`
}

type EhfListresponse struct {
	Ehfs []*data.Ehf `json:"ehfs"`
}

func (app *application) listEhfHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := loggingutils.LoggerFromContext(ctx)

	ehfs, err := app.models.Ehf.List(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "unable to get ehfs", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.InfoContext(ctx, "returning ehfs")

	rest.WriteJSON(
		w,
		http.StatusOK,
		EhfListresponse{Ehfs: ehfs},
		nil,
	)
}

func (app *application) getEhfHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		rest.BadRequestResponse(w, r, "unable to parse id from path")
		return
	}

	ehf, err := app.models.Ehf.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			rest.NotFoundResponse(w, r)
		default:
			logger.ErrorContext(ctx, "unable to get ehf", "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.InfoContext(ctx, "returning ehf")

	rest.WriteJSON(
		w,
		http.StatusOK,
		EhfResponse{
			EhfID:          ehf.EhfID,
			FileName:       ehf.FileName,
			CustomerID:     ehf.CustomerID,
			SupplierID:     ehf.SupplierID,
			InvoiceNo:      ehf.InvoiceNo,
			BuyerReference: ehf.BuyerReference,
			IssueDate:      ehf.IssueDate,
			DueDate:        ehf.DueDate,
			Currency:       ehf.Currency,
			Amount:         ehf.Amount,
		},
		nil,
	)
}

func (app *application) createEhfHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	var req data.Ehf
	err := rest.ReadJSON(r, &req)
	if err != nil {
		rest.BadRequestResponse(w, r, "unable to decode data from request")
		return
	}

	ehf := &data.Ehf{
		EhfID:          req.EhfID,
		FileName:       req.FileName,
		CustomerID:     req.CustomerID,
		SupplierID:     req.SupplierID,
		InvoiceNo:      req.InvoiceNo,
		BuyerReference: req.BuyerReference,
		IssueDate:      req.IssueDate,
		DueDate:        req.DueDate,
		Currency:       req.Currency,
		Amount:         req.Amount,
	}

	result, err := app.models.Ehf.Insert(ctx, ehf)
	if err != nil {
		logger.ErrorContext(ctx, "unable to insert ehf", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.InfoContext(ctx, "ehf created")

	rest.WriteJSON(
		w,
		http.StatusCreated,
		EhfResponse{
			FileName:       result.FileName,
			CustomerID:     result.CustomerID,
			SupplierID:     result.SupplierID,
			InvoiceNo:      result.InvoiceNo,
			BuyerReference: result.BuyerReference,
			IssueDate:      result.IssueDate,
			DueDate:        result.DueDate,
			Currency:       result.Currency,
			Amount:         result.Amount,
		},
		nil,
	)
}

func (app *application) deleteEhfHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		rest.BadRequestResponse(w, r, "unable to parse id from path")
		return
	}

	err = app.models.Ehf.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			rest.NotFoundResponse(w, r)
		default:
			logger.ErrorContext(ctx, "unable to delete ehf", "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.InfoContext(ctx, "ehf deleted")

	rest.WriteJSON(w, http.StatusNoContent, nil, nil)
}
