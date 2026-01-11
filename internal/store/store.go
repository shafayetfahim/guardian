package store

import (
	"database/sql"

	"github.com/shafayetfahim/guardian/internal/extractor"
)

func SaveAsset(db *sql.DB, category string, meta *extractor.PhotoMetadata) error {
	query := `INSERT INTO daily_logs (category, content) 
              VALUES ($1, $2)`
	_, err := db.Exec(query, category, meta)
	return err
}
