# URL Shortener with Tagging

A URL shortener service built with Go, Redis, and Gin, featuring tagging functionality. This service allows users to shorten long URLs, add tags, and manage their links efficiently.

## Features

- Shorten long URLs
- Add tags to shortened URLs
- Custom short URL option
- Set expiry time for shortened URLs
- Rate limiting based on client IP

## Technologies Used

- Go
- Gin (web framework)
- Redis (for storage and rate limiting)
- Docker (for containerization)

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Create a `.env` file based on the provided `.env.example`:
   ```bash
   cp .env.example .env
   ```

3. Modify the `.env` file if necessary (e.g., change database address, API quota).

### Running the Application

1. Start the services using Docker Compose:
   ```bash
   docker-compose up
   ```

2. The API will be available at `http://localhost:8080`.

### API Endpoints

#### Shorten URL

- **POST** `/api/v1/shorten`
- **Request Body**:
  ```json
  {
      "url": "https://www.example.com",
      "custom_short": "optional_custom_short",
      "expiry": 24,
      "tags": ["tag1", "tag2"]
  }
  ```
- **Response**:
  ```json
  {
      "url": "https://www.example.com",
      "custom_short": "http://localhost:8080/abc123",
      "expiry": 24,
      "tags": ["tag1", "tag2"],
      "x_rate_remaining": 10,
      "x_rate_limit_reset": 30
  }
  ```

### Testing

You can test the API using tools like Postman or curl. Here are some example curl commands:

```bash
# Basic URL shortening with tags
curl -X POST http://localhost:8080/api/v1/shorten \
-H "Content-Type: application/json" \
-d '{
    "url": "https://www.google.com",
    "custom_short": "google",
    "expiry": 24,
    "tags": ["search", "google"]
}'
```

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [Redis](https://redis.io/) - In-memory data structure store
