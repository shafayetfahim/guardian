package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/shafayetfahim/guardian/internal/crawler"
	"github.com/shafayetfahim/guardian/internal/extractor"
	"github.com/shafayetfahim/guardian/internal/store"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbUrl := os.Getenv("DB_URL")
	ingestPath := os.Getenv("INGEST_PATH")

	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDatabase(db)

	fmt.Printf("🛡️  Guardian: Scanning %s...\n", ingestPath)

	files, err := crawler.Search(ingestPath, []string{".jpg", ".jpeg", ".png", ".ARW", ".heic"})
	if err != nil {
		log.Fatal(err)
	}

	indexedCount := 0
	for _, file := range files {
		meta, err := extractor.ExtractPhotoData(file)
		if err != nil {
			meta = map[string]interface{}{"path": file}
		}

		err = store.SaveAsset(db, file, "Photography", meta)
		if err != nil {
			fmt.Printf("❌ Failed to save %s: %v\n", file, err)
			continue
		}
		indexedCount++
	}

	fmt.Printf("✅ Success! Truly indexed %d files into the Vault.\n", indexedCount)
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
