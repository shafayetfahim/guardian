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
	// 1. Load configuration
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system defaults")
	}

	dbUrl := os.Getenv("DB_URL")
	ingestPath := os.Getenv("INGEST_PATH")

	// 2. Database Connection
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. Initialize Schema
	initDatabase(db)

	fmt.Printf("🛡️  Guardian: Scanning %s...\n", ingestPath)

	// 4. Run the Crawler
	files, err := crawler.Search(ingestPath, []string{".jpg", ".png", ".ARW"})
	if err != nil {
		log.Fatal(err)
	}

	// 5. Save to Store
	for _, file := range files {
		store.SaveAsset(db, file, "Photography")
	}

	fmt.Printf("✅ Success! Indexed %d files into the Vault.\n", len(files))
}

func initDatabase(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS daily_logs (
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
