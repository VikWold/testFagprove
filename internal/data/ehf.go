package data

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"testFagprove/internal/loggingutils"
	"time"

	"github.com/google/uuid"
)

type Ehf struct {
	ID             uuid.UUID  `json:"id"`
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

type EhfModel struct {
	Timeout *time.Duration
	DB      *sql.DB
}

func (e EhfModel) Get(ctx context.Context, id uuid.UUID) (*Ehf, error) {
	logger := loggingutils.LoggerFromContext(ctx)

	stmt := `
SELECT CAST(id AS CHAR(36)),
	file_name,
	customer_id,
	supplier_id,
	invoice_number,
	buyer_reference,
	issue_date,
	due_date,
	currency,
	amount
FROM public.ehf
WHERE id = $1
	`

	var ehf Ehf

	ctx, cancel := context.WithTimeout(ctx, *e.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", stmt),
			slog.String("id", id.String()),
		),
	)

	logger.InfoContext(ctx, "performing query")
	err := e.DB.QueryRowContext(ctx, stmt, id).Scan(
		&ehf.ID,
		&ehf.FileName,
		&ehf.CustomerID,
		&ehf.SupplierID,
		&ehf.InvoiceNo,
		&ehf.BuyerReference,
		&ehf.IssueDate,
		&ehf.DueDate,
		&ehf.Currency,
		&ehf.Amount,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.InfoContext(ctx, "no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.InfoContext(ctx, "an error occured while performing query", "error", err)
		}
	}

	return &ehf, nil
}

func (e EhfModel) List(ctx context.Context) ([]*Ehf, error) {
	logger := loggingutils.LoggerFromContext(ctx)

	stmt := `
SELECT id,
	file_name,
	customer_id,
	supplier_id,
	invoice_number,
	buyer_reference,
	issue_date,
	due_date,
	currency,
	amount
FROM public.ehf
WHERE id = $1
	`

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", stmt),
		),
	)

	ctx, cancel := context.WithTimeout(ctx, *e.Timeout)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, stmt)
	if err != nil {
		logger.ErrorContext(ctx, "error executing query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var el []*Ehf

	for rows.Next() {
		var ehf Ehf

		err := rows.Scan(
			&ehf.ID,
			&ehf.FileName,
			&ehf.CustomerID,
			&ehf.SupplierID,
			&ehf.InvoiceNo,
			&ehf.BuyerReference,
			&ehf.IssueDate,
			&ehf.DueDate,
			&ehf.Currency,
			&ehf.Amount,
		)
		if err != nil {
			logger.ErrorContext(ctx, "error scanning row", "error", err)
			return nil, err
		}

		el = append(el, &ehf)
	}

	if err = rows.Err(); err != nil {
		logger.ErrorContext(ctx, "error with row iteration", "error", err)
		return nil, err
	}

	logger.InfoContext(ctx, "successfully fetched all ehfs")

	return el, nil
}

func (e EhfModel) Insert(ctx context.Context, ehf *Ehf) (*Ehf, error) {
	logger := loggingutils.LoggerFromContext(ctx)

	stmt := `
INSERT INTO public.ehf (
	id,
	file_name,
	customer_id,
	supplier_id,
	invoice_number,
	buyer_reference,
	issue_date,
	due_date,
	currency,
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, file_name, customer_id, supplier_id, invoice_number, buyer_reference, issue_date, due_date, currency;
	`

	ctx, cancel := context.WithTimeout(ctx, *e.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", stmt),
			"ehf", ehf,
		),
	)

	var result Ehf

	logger.InfoContext(ctx, "performing query")
	err := e.DB.QueryRowContext(
		ctx,
		stmt,
		ehf.ID,
		ehf.FileName,
		ehf.CustomerID,
		ehf.SupplierID,
		ehf.InvoiceNo,
		ehf.BuyerReference,
		ehf.IssueDate,
		ehf.DueDate,
		ehf.Currency,
	).Scan(
		&result.ID,
		&result.FileName,
		&result.CustomerID,
		&result.SupplierID,
		&result.InvoiceNo,
		&result.BuyerReference,
		&result.IssueDate,
		&result.DueDate,
		&result.Currency,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.InfoContext(ctx, "no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.InfoContext(ctx, "an error occured while performing query", "error", err)
			return nil, err
		}
	}

	return &result, nil
}

func (e *EhfModel) Delete(ctx context.Context, ehfID uuid.UUID) error {
	logger := loggingutils.LoggerFromContext(ctx)
	stmt := `
DELETE FROM public.ehf
WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, *e.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", stmt),
			slog.String("ehf_id", ehfID.String()),
		),
	)

	_, err := e.DB.ExecContext(ctx, stmt, ehfID)

	if err != nil {
		logger.Error(
			"an error ocurred while trying to delete the ehf",
			slog.String("error", err.Error()),
		)
		return err
	}

	logger.Info("ehf deleted successfully")
	return nil
}
