# Users API

# Location User API

The **Location User API** is a microservice responsible for managing user location updates and providing nearby user search functionality. It is part of a broader system that includes a second microservice for historical location tracking and distance calculation.

## Features

- **Update current user location**: Accepts user coordinates and stores them in a PostgreSQL database.
- **Search for nearby users**: Find users within a given radius of a location using geospatial queries.
- **Forward updates to gRPC Location Tracker API**: Integrates with a separate tracking service via gRPC.

## API Endpoints

### POST `/location`

Updates the current location of a user.

#### Request Body

```json
{
  "username": "kirill",
  "latitude": 55.751244,
  "longitude": 37.618423
}
```

- `username`: 4–16 alphanumeric characters.
- `latitude`, `longitude`: Coordinates with up to 8 decimal places.

### GET `/nearby`

Search for users near a specific location.

#### Query Parameters

- `latitude`: Latitude of the search center.
- `longitude`: Longitude of the search center.
- `radius`: Search radius in kilometers.
- `page` (optional): Page number (default: 1).
- `limit` (optional): Number of results per page (default: 10).

#### Example

```
GET /nearby?latitude=55.751244&longitude=37.618423&radius=5
```

## Architecture

- Written in Go using the Gin web framework.
- PostgreSQL for persistent storage of user location.
- gRPC integration with `LocationTracker` service for historical tracking.

## Development

### Running Locally

Make sure you have a PostgreSQL instance running.

Set the environment variable:

```bash
export DATABASE_URL=postgres://user:password@localhost:5432/location?sslmode=disable
```

Start the server:

```bash
go run main.go
```

### Docker

Build and run using Docker Compose:

```bash
docker-compose up --build
```

## License

MIT