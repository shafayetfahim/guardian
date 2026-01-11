Guardian, a Metadata Indexer

I'm into photography, but it can be difficult to sift through a lot of RAW/JPEG files trying to find "the right shot". I wanted a way to pre-filter photos, which requires indexing hundreds, if not thousands, of photos by their EXIF metadata.

That's the role Guardian serves. It's a concurrent file indexing system built in Go for managing photography metadata. It scans a directory for images, extracts EXIF metadata (like the lens model, aperture, or ISO value), and stores it in a PostgreSQL database. It is designed to handle local image libraries using a modular backend architecture.

Here's some features I implemented:

- Concurrent Processing: 
The system uses a worker pool to process files in parallel. By limiting the number of active goroutines, it prevents disk I/O bottlenecks and manages CPU usage during large scans.

- Data Validation: 
A validation layer checks file health before extraction. It verifies file existence, read permissions, and size limits to ensure the metadata parser does not encounter corrupted or excessively large files.

- Flexible Storage: 
The database uses PostgreSQL JSONB columns. This allows the system to store varied metadata formats from different camera manufacturers without requiring schema migrations.

- Database Abstraction: 
A SQL view maps the JSONB data into relational formatting. This allows for SQL queries on specific photo data, like the camera model or lens type.

Here's the internal structure: 
* cmd/api: The entry point & database initialization.
* internal/crawler: This defines the logic for the file system crwaler.
* internal/extractor: Parses Binary EXIF data.
* internal/validator: File integrity and permission checks.
* internal/store: Database persistence.

To run this... 
1. Configure the .env file with DB_URL and INGEST_PATH.
2. Start the database: docker-compose up -d.
3. Run the indexer: go run cmd/api/main.go.

In the future, I'd like to continue working on this tool to model data, and maybe create functions that pre-determine which photos should be re-assessed based on an exposure triangle score. 
