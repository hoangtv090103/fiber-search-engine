# Fiber Search Engine

This is a search engine project built with Go and the Fiber framework.

## Project Structure

- `.air.toml`, `.env`, `.gitignore`: Configuration files.
- `.vscode/`: Contains settings for Visual Studio Code.
- `.zed/`: Contains task configurations.
- `db/`: Contains database related files like `index.go`, `search_settings.go`, `user.go`.
- `go.mod`, `go.sum`: Go module and checksum files.
- `main.go`: The entry point of the application.
- `routes/`: Contains routing related files like `admin.go`, `routes.go`.
- `tmp/`: Contains temporary files and build logs.
- `utills/`: Contains utility files like `jwt.go`.
- `views/`: Contains view templates and related Go files.

## Setup

1. Clone the repository.
2. Install the dependencies with `go mod download`.
3. Copy `.env.example` to `.env` and fill in your environment variables.
4. Run the application with `go run main.go`.

## Usage

- The application provides a search engine functionality.
- It uses the Fiber framework for handling requests and responses.
- It uses JWT for authentication.