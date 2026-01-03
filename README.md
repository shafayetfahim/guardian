Guardian is a concurrent file indexing system built in Go for managing photography metadata. It scans a directory for image files, extracts technical metadata (Lens, Aperture, ISO), and stores it in a PostgreSQL database. It is designed to handle local image libraries using a modular backend architecture.

Here's some features I implemented:

- Concurrent Processing
The system uses a worker pool to process files in parallel. By limiting the number of active goroutines, it prevents disk I/O bottlenecks and manages CPU usage during large scans.

- Data Validation
A validation layer checks file health before extraction. It verifies file existence, read permissions, and size limits to ensure the metadata parser does not encounter corrupted or excessively large files.

- Flexible Storage
The database uses PostgreSQL JSONB columns. This allows the system to store varied metadata formats from different camera manufacturers without requiring schema migrations.

- Database Abstraction
A SQL view maps the JSONB data into relational formatting. This allows for SQL queries on specific photo data, like the camera model or lens type.

Here's the internal structure:
* cmd/api: Application entry point and database initialization.
* internal/crawler: Filesystem traversal logic.
* internal/extractor: Binary EXIF data parsing.
* internal/validator: File integrity and permission checks.
* internal/store: Database persistence layer.

To run this...
1. Configure the .env file with DB_URL and INGEST_PATH.
2. Start the database: docker-compose up -d.
3. Run the indexer: go run cmd/api/main.go.
