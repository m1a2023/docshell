package repository

import (
	"context"
	"database/sql"
	"docshell/internal/v1/docs/models"
)

func GetAllDocuments(ctx context.Context, con *sql.DB) ([]models.Document, error) {
	// Get all documents
	rows, err := con.QueryContext(ctx, get_all_documents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build response
	res := []models.Document{}
	for rows.Next() {
		doc := models.Document{}
		// Scan row
		if err := rows.Scan(&doc.Id, &doc.AuthorId, &doc.UploaderId,
			&doc.Title, &doc.Size, &doc.Path, &doc.Hash, &doc.CreatedAt, &doc.ChangedAt); err != nil {
			return nil, err
		}
		res = append(res, doc)
	}

	return res, err
}

func GetDocumentById(ctx context.Context, con *sql.DB, id int) (models.Document, error) {
	// Get document
	rows, err := con.QueryContext(ctx, get_document_by_id, id)
	if err != nil {
		return models.Document{}, err
	}

	// Build response
	// In fact, there is one itereation because of query
	doc := models.Document{}
	for rows.Next() {
		// Scan row
		if err := rows.Scan(&doc.Id, &doc.AuthorId, &doc.UploaderId,
			&doc.Title, &doc.Size, &doc.Path, &doc.Hash, &doc.CreatedAt, &doc.ChangedAt); err != nil {
			return models.Document{}, err
		}
	}
	return doc, nil
}

// TODO
func CreateDocument(ctx context.Context, con *sql.DB, dc models.DocumentCreation) (models.Document, error) {
	return models.Document{}, nil
}
