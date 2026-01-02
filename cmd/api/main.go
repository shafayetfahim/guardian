package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/shafayetfahim/guardian/internal/crawler"
	"github.com/shafayetfahim/guardian/internal/extractor"
	"github.com/shafayetfahim/guardian/internal/store"
	"github.com/shafayetfahim/guardian/internal/validator"
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

	db.SetMaxOpenConns(10)

	initDatabase(db)

	files, err := crawler.Search(ingestPath, []string{".jpg", ".jpeg", ".png", ".ARW", ".heic"})
	if err != nil {
		log.Fatal(err)
	}

	tasks := make(chan string, len(files))
	var wg sync.WaitGroup
	const workerCount = 5

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for file := range tasks {
				if ok, err := validator.IsValidMedia(file); !ok {
					fmt.Printf("[Worker %d] ⚠️ Skipping %s: %v\n", id, file, err)
					continue
				}

				meta, err := extractor.ExtractPhotoData(file)
				if err != nil {
					meta = map[string]interface{}{"path": file, "error": "metadata_extraction_failed"}
				}

				if err := store.SaveAsset(db, file, "Photography", meta); err != nil {
					fmt.Printf("[Worker %d] ❌ Error saving %s: %v\n", id, file, err)
				}
			}
		}(i)
	}

	for _, file := range files {
		tasks <- file
	}
	close(tasks)

	wg.Wait()

	fmt.Printf("✅ Success! Indexed %d files in %v using %d workers.\n", len(files), time.Since(start), workerCount)
}

func initDatabase(db *sql.DB) {
	tableQuery := `CREATE TABLE IF NOT EXISTS daily_logs (
		id SERIAL PRIMARY KEY,
		category TEXT NOT NULL,
		content JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(tableQuery); err != nil {
		log.Fatalf("Table migration failed: %v", err)
	}

	viewQuery := `CREATE OR REPLACE VIEW photo_analytics AS
	SELECT 
		id,
		content->>'path' as path,
		content->>'camera' as camera,
		content->>'lens' as lens,
		content->>'aperture' as aperture,
		content->>'iso' as iso,
		created_at
	FROM daily_logs
	WHERE category = 'Photography';`
	if _, err := db.Exec(viewQuery); err != nil {
		log.Fatalf("View creation failed: %v", err)
	}
}
