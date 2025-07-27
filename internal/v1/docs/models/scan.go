package models

import "database/sql"

func ScanDocument(rows *sql.Rows) (Document, error) {
	doc := Document{}
	if err := rows.Scan(&doc.Id, &doc.AuthorId, &doc.UploaderId,
		&doc.Title, &doc.Size, &doc.Path, &doc.Hash, &doc.CreatedAt, &doc.ChangedAt); err != nil {
		return Document{}, err
	}
	return doc, nil
}
