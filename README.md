# URL Shortener using Go

A simple URL shortener service built with Go.

## Overview
- Shorten any URL through a concise endpoint.
- Fast HTTP redirection to the original URL.
- Tracks click counts for basic usage analytics.

## Features
- Create short URLs by sending a POST request to `/encurtar`.
- Redirect using the generated `/{shortCode}`.
- View statistics (click count, etc.) at `/stats/{shortCode}`.
- PostgreSQL integration (optional) for data persistence.
- Docker support (optional) for easy containerization.

## Installation
1. Install Go (version 1.18+ recommended).
2. Clone this repository:
   ```
   git clone https://github.com/Arth19/url-shortener-go.git
   cd url-shortener-go
   ```
3. Configure environment variables in a `.env` file (for example, `DATABASE_URL` if you are using PostgreSQL).

## Usage
1. Run the server:
   ```
   go run ./cmd/shortener/main.go
   ```
2. By default, the application listens on port 8080:
   - POST `/encurtar` – Create a new shortened URL by sending JSON with the field "url".
   - GET `/{shortCode}` – Redirects to the original URL.
   - GET `/stats/{shortCode}` – Retrieves click-count statistics in JSON format.
3. Test the endpoints using curl, Postman, or any REST client:
   ```
   # Example: Shorten a URL
   curl -X POST -H "Content-Type: application/json" \
        -d '{"url":"https://example.com"}' \
        http://localhost:8080/encurtar

   # Example: Access stats
   curl http://localhost:8080/stats/abc123
   ```

## Contributing
Contributions are welcome! Feel free to open issues or submit pull requests.

## License
This project is released under the MIT License.
