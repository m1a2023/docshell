package repository

import (
	"context"
	"database/sql"
	"docshell/internal/v1/docs/models"
	"docshell/internal/v1/storage"
)

func GetAllDocuments(ctx context.Context, con *sql.DB) ([]models.Document, error) {
	// Get all documents
	rows, err := con.QueryContext(ctx, get_all_documents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build response
	docs, err := storage.ScanMany(rows, models.ScanDocument)
	if err != nil {
		return nil, err
	}

	return docs, err
}

func GetDocumentById(ctx context.Context, con *sql.DB, id int) (models.Document, error) {
	// Get document
	rows, err := con.QueryContext(ctx, get_document_by_id, id)
	if err != nil {
		return models.Document{}, err
	}
	defer rows.Close()

	// Build response
	doc, err := storage.ScanSingle(rows, models.ScanDocument)
	if err != nil {
		return models.Document{}, err
	}
	return doc, nil
}

func CreateDocument(ctx context.Context, con *sql.DB, dc models.DocumentCreation) (models.Document, error) {
	// Insert document and return it
	rows, err := con.QueryContext(ctx, insert_document,
		dc.AuthorId, dc.UploaderId, dc.Title, dc.Size, dc.Path, dc.Hash,
	)
	if err != nil {
		return models.Document{}, nil
	}
	defer rows.Close()

	// Build response
	doc, err := storage.ScanSingle(rows, models.ScanDocument)
	if err != nil {
		return models.Document{}, nil
	}

	return doc, nil
}

func GetManyDocementsByRowWithArg(ctx context.Context, con *sql.DB, row string, arg any) ([]models.Document, error) {
	rows, err := con.QueryContext(ctx, get_documents_by_, row, arg)
	if err != nil {
		return nil, err
	}

	// Build response
	docs, err := storage.ScanMany(rows, models.ScanDocument)
	if err != nil {
		return nil, err
	}

	return docs, nil
}
