# Fiber Search Engine

This search engine is built with Go and the Fiber framework. It uses an in-memory inverted index for quick search results and a web crawler for fetching data from the web.

## Project Structure

- `db/`: Contains database related files like `index.go`, `search_index.go`, `search_settings.go`, `url.go`, `user.go`. These files handle the database operations.
- `main.go`: The entry point of the application.
- `routes/`: Contains routing related files like `admin.go`, `routes.go`, `search.go`. These files handle the routing logic for the application.
- `search/`: Contains search engine related files like `crawler.go`, `crawler_test.go`, `engine.go`, `indexer.go`, `tokenizer.go`. These files implement the search engine functionality.
- `utils/`: Contains utility files like `cron.go`, `jwt.go`. These files provide utility functions like JWT authentication and scheduling cron jobs.
- `views/`: Contains view templates and related Go files. These files handle the rendering of the user interface.

## How It Works

1. Crawling: The search engine starts by crawling the web. This is done by the `crawler.go` file. It fetches data from the web and extracts useful information such as the page title, description, headings, and external links.

2. Indexing: The extracted data is then indexed by the `indexer.go` file. It creates an in-memory inverted index, which is a data structure that maps tokens (words) to the URLs where they were found. This allows for quick search results.

3. Searching: When a search query is received, it is tokenized by the `tokenizer.go` file. The tokens are then used to search the index and return the matching URLs.

4. Updating: The search engine is updated every hour by a cron job defined in `cron.go`. This ensures that the search results are always up-to-date.

5. User Interface: The user interface is rendered by the files in the `views/` directory. It provides a form for users to enter their search queries and displays the search results.

## User Settings

Users can customize their search settings through the user interface. They can set the number of URLs to be crawled per hour and choose whether to add new URLs to the database. These settings are handled by the `index.templ` file.
