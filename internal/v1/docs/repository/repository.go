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

		if err := rows.Scan(&doc.Id, &doc.AuthorId, &doc.UploaderId,
			&doc.Title, &doc.Size, &doc.Path, &doc.Hash, &doc.CreatedAt, &doc.ChangedAt); err != nil {
			continue
		}

		res = append(res, doc)
	}

	return res, err
}
