package storage

import (
	"database/sql"
)

func ScanSingle[T any](rows *sql.Rows, scan func(*sql.Rows) (T, error)) (T, error) {
	var result T
	if rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return result, err
		}
		result = item
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func ScanMany[T any](rows *sql.Rows, scan func(*sql.Rows) (T, error)) ([]T, error) {
	var result []T
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
