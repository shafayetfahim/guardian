package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time" // For measuring performance

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/shafayetfahim/guardian/internal/crawler"
	"github.com/shafayetfahim/guardian/internal/extractor"
	"github.com/shafayetfahim/guardian/internal/store"
)

func main() {
	start := time.Now()
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

	// Set connection pool limits - another High-Signal move
	db.SetMaxOpenConns(10)

	initDatabase(db)

	files, err := crawler.Search(ingestPath, []string{".jpg", ".jpeg", ".png", ".ARW"})
	if err != nil {
		log.Fatal(err)
	}

	// --- CONCURRENCY BLOCK ---
	tasks := make(chan string, len(files))
	var wg sync.WaitGroup
	const workerCount = 5 // Opinionated choice: stay within MacBook I/O limits

	// Start Workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for file := range tasks {
				meta, err := extractor.ExtractPhotoData(file)
				if err != nil {
					meta = map[string]interface{}{"path": file, "error": "metadata_extraction_failed"}
				}

				if err := store.SaveAsset(db, file, "Photography", meta); err != nil {
					fmt.Printf("[Worker %d] ❌ Error saving %s: %v\n", workerID, file, err)
				}
			}
		}(i)
	}

	// Feed Tasks
	for _, file := range files {
		tasks <- file
	}
	close(tasks) // Signal workers to stop when queue is empty

	wg.Wait()
	// -------------------------

	fmt.Printf("✅ Success! Indexed %d files in %v using %d workers.\n", len(files), time.Since(start), workerCount)
}

func initDatabase(db *sql.DB) {
	tableQuery := `CREATE TABLE IF NOT EXISTS daily_logs (
		id SERIAL PRIMARY KEY,
		category TEXT NOT NULL,
		content JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	db.Exec(tableQuery)

	viewQuery := `CREATE OR REPLACE VIEW photo_analytics AS
	SELECT id, content->>'path' as path, content->>'camera' as camera, 
	content->>'lens' as lens, content->>'aperture' as aperture, created_at
	FROM daily_logs WHERE category = 'Photography';`
	db.Exec(viewQuery)
}
