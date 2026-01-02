package store

import (
	"database/sql"
)

// SaveAsset stores file paths into the JSONB column of our DB
func SaveAsset(db *sql.DB, path string, category string) error {
	// We add ::TEXT so PostgreSQL knows exactly what type the path is
	query := `INSERT INTO daily_logs (category, content) 
              VALUES ($1, jsonb_build_object('path', $2::TEXT))`
	_, err := db.Exec(query, category, path)
	return err
}
