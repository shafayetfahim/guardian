package store
package store

import (
	"database/sql"
)

// SaveAsset stores the file path and metadata in PostgreSQL
func SaveAsset(db *sql.DB, path string, category string) error {
	query := `INSERT INTO daily_logs (category, content) 
              VALUES ($1, jsonb_build_object('path', $2))`
	_, err := db.Exec(query, category, path)
	return err
}