package store

import "database/sql"

// SaveAsset stores file paths into the JSONB column of our DB
func SaveAsset(db *sql.DB, path string, category string) error {
	query := `INSERT INTO daily_logs (category, content) 
              VALUES ($1, jsonb_build_object('path', $2))`
	_, err := db.Exec(query, category, path)
	return err
}
