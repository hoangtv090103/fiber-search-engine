# Fiber Search Engine

This is a search engine project built with Go and the Fiber framework.

## Project Structure

- `db/`: Contains database related files like `index.go`, `search_index.go`, `search_settings.go`, `url.go`, `user.go`.
- `go.mod`, `go.sum`: Go module and checksum files.
- `main.go`: The entry point of the application.
- `routes/`: Contains routing related files like `admin.go`, `routes.go`, `search.go`.
- `search/`: Contains search engine related files like `crawler.go`, `crawler_test.go`, `engine.go`, `indexer.go`, `tokenizer.go`.
- `tmp/`: Contains temporary files and build logs.
- `utills/`: Contains utility files like `cron.go`, `jwt.go`.
- `views/`: Contains view templates and related Go files.

## Implementation

The project is implemented in Go using the Fiber framework for handling HTTP requests and responses. It uses JWT for authentication. The search engine functionality is implemented in the `search/` directory, with a crawler for fetching data, an indexer for indexing the data, and a tokenizer for tokenizing the search queries.

The database related operations are handled in the `db/` directory. The `routes/` directory contains the routing logic for the application.

## Setup

1. Clone the repository.
2. Install the dependencies with `go mod download`.
3. Copy `.env.example` to `.env` and fill in your environment variables.
4. Run the application with `go run main.go`.

## Usage

- The application provides a search engine functionality.
- It uses the Fiber framework for handling requests and responses.
- It uses JWT for authentication.
- It uses a cron job, defined in [`utills/cron.go`](utills/cron.go), to run the search engine every hour.
