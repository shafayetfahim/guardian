package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/shafayetfahim/guardian/internal/crawler"
	"github.com/shafayetfahim/guardian/internal/store"
)

func main() {
	// 1. Load configuration from .env file
	// This keeps your personal file paths out of GitHub.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system defaults")
	}

	dbUrl := os.Getenv("DB_URL")
	ingestPath := os.Getenv("INGEST_PATH")

	// 2. Establish Database Connection
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer db.Close()

	// 3. Initialize Schema (Surgical Catch: Ensure tables exist before scanning)
	// This makes your system "Self-Healing".
	initDatabase(db)

	fmt.Printf("🛡️  Guardian: Scanning %s for assets...\n", ingestPath)

	// 4. Crawl the target folder
	// We are looking for common image formats.
	files, err := crawler.Search(ingestPath, []string{".jpg", ".png", ".ARW"})
	if err != nil {
		log.Fatalf("Crawl failed: %v", err)
	}

	// 5. Persist discovered assets to PostgreSQL
	count := 0
	for _, file := range files {
		err := store.SaveAsset(db, file, "Photography")
		if err != nil {
			fmt.Printf("⚠️  Error saving %s: %v\n", file, err)
			continue
		}
		count++
	}

	fmt.Printf("✅ Success! Indexed %d files into the Vault.\n", count)
}

// initDatabase ensures the necessary tables exist in the database.
func initDatabase(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS daily_logs (
		id SERIAL PRIMARY KEY,
		category TEXT NOT NULL,
		content JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
