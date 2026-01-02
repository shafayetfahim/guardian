package store

import (
	"database/sql"
)

// SaveAsset now accepts a metadata map to fill our JSONB column
func SaveAsset(db *sql.DB, path string, category string, meta map[string]interface{}) error {
	// Add the path into the metadata map
	meta["path"] = path

	query := `INSERT INTO daily_logs (category, content) 
              VALUES ($1, $2)`

	// We pass the entire map; the pgx driver handles the JSON conversion
	_, err := db.Exec(query, category, meta)
	return err
}
