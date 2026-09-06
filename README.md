# open-movies-db

A high performance REST API built in Go for managing structured movie metadata. Features a layered architecture, relational database transactions, dynamic search, and strict data integrity.


## Architecture & Design

This project implements a clean, layered architecture to ensure separation of concerns, maintainability, and testability:

- **Handler Layer:** Manages HTTP request parsing, routing, and JSON serialization.
- **Service Layer:** Encapsulates business logic, data validation, and transactional orchestration.
- **Repository Layer:** Interfaces directly with the SQLite database, executing raw SQL queries and managing complex `JOIN` operations.

## Core Capabilities

- **Relational Data Management:** Full CRUD support for Movies, Actors, and Genres.
- **Complex Associations:** Handles Many-to-Many relationships (Movies ↔ Genres, Movies ↔ Actors) via junction tables.
- **Dynamic Querying:**
  - Case-insensitive partial text search (`/search?title=matrix`)
- **Strict Data Integrity:** Enforces `ON DELETE RESTRICT` foreign key constraints at the database level to prevent orphaned records.
- **Transactional Operations:** Implements database transactions (`BEGIN`, `COMMIT`, `ROLLBACK`) for multi-table inserts and safe `?force=true` cascading deletions.

## Tech Stack

- **Language:** Go (1.22+)
- **Database:** SQLite3 (`github.com/mattn/go-sqlite3`)
- **Routing:** Go Standard Library (`net/http` )
- **Containerization:** Docker & Docker Compose
- **Testing:** Go testing with isolated in-memory SQLite databases


## Setup and Installation

### Local Development

Ensure Go 1.22+ and a C compiler (required for `go-sqlite3` CGO bindings) are installed.

1. Clone the repository:
   ```bash
   git clone https://gitea.kood.tech/muhammadowaisjaved/movies-api.git
   cd movies-api
   ```
2. Download dependencies:

   ```bash
   go mod tidy
   ```

3. Run the application:

   ```bash
   go run ./cmd/api/main.go
   ```


## Containerization

The application includes a multi-stage Dockerfile and docker-compose.yml for consistent deployment across environments.

    docker-compose up --build

The API will be exposed on http://localhost:8080. The SQLite database file is persisted via a Docker volume in the ./data directory.

## Testing

The project includes a comprehensive automated testing suite. Tests are executed against isolated, in-memory SQLite databases (:memory:) to ensure no impact on the local development environment.

To run the unit tests in the root directory:

    go test ./... -v

## API Reference

### Movies

| Method   | Endpoint             | Description                                                                  |
| :------- | :------------------- | :--------------------------------------------------------------------------- |
| `POST`   | `/api/movies`        | Create a movie and establish actor/genre relationships                       |
| `GET`    | `/api/movies`        | Retrieve movies (Supports `size`, `year`, `genre` and `actor` query params)  |
| `GET`    | `/api/movies/search` | Search movies by title (`?title=`)                                           |
| `GET`    | `/api/movies/{id}`   | Retrieve a specific movie                                                    |
| `PATCH`  | `/api/movies/{id}`   | Partially update movie attributes or relationships                           |
| `DELETE` | `/api/movies/{id}`   | Delete a movie (Fails if relationships exist unless `?force=true` is passed) |

### Genres & Actors

| Method   | Endpoint                                 | Description                               |
| :------- | :--------------------------------------- | :---------------------------------------- |
| `POST`   | `/api/genres` or `/api/actors`           | Create a new entity                       |
| `GET`    | `/api/genres` or `/api/actors`           | Retrieve all entities                     |
| `GET`    | `/api/genres/{id}` or `/api/actors/{id}` | Retrieve a specific entity                |
| `PATCH`  | `/api/genres/{id}` or `/api/actors/{id}` | Partially update an entity                |
| `DELETE` | `/api/genres/{id}` or `/api/actors/{id}` | Delete an entity (Supports `?force=true`) |



## Usage & Payload Examples

Below are examples of how to interact with the API using curl. You can use these JSON payloads in Postman as well.

### 1. Genres

Create a Genre:

```Bash
curl -X POST http://localhost:8080/api/genres \
-H "Content-Type: application/json" \
-d '{
  "name": "Action"
}'
```

Update a Genre (PATCH ):

```Bash
curl -X PATCH http://localhost:8080/api/genres/1 \
-H "Content-Type: application/json" \
-d '{
  "name": "Action & Adventure"
}'
```

### 2. Actors

Create an Actor:

_(Note: birthDate must strictly follow the ISO 8601 YYYY-MM-DD format )_

```Bash
curl -X POST http://localhost:8080/api/actors \
-H "Content-Type: application/json" \
-d '{
  "name": "Leonardo DiCaprio",
  "birthDate": "1974-11-11"
}'
```

Update an Actor (PATCH ):

```Bash
curl -X PATCH http://localhost:8080/api/actors/1 \
-H "Content-Type: application/json" \
-d '{
  "name": "Leo DiCaprio"
}'
```

### 3. Movies

Create a Movie with Relationships:

_(Note: The genres and actors arrays expect objects containing the id of existing records )_

```Bash
curl -X POST http://localhost:8080/api/movies \
-H "Content-Type: application/json" \
-d '{
  "title": "Inception",
  "releaseYear": 2010,
  "duration": 148,
  "genres": [{"id": 1}, {"id": 2}],
  "actors": [{"id": 1}, {"id": 3}]
}'
```

Update a Movie (PATCH ):

_(You can send just one field, or multiple. Providing a new array for relationships will overwrite the old ones)_

```Bash
curl -X PATCH http://localhost:8080/api/movies/1 \
-H "Content-Type: application/json" \
-d '{
  "duration": 150,
  "actors": [{"id": 1}]
}'
```

### 4. Filtering & Searching (GET )

The API supports dynamic querying via URL parameters.

Search by Partial Name/Title:

```Bash
curl -X GET "http://localhost:8080/api/movies/search?title=incept"
curl -X GET "http://localhost:8080/api/actors?name=Leo"
```

Filter Movies by Relationships or Year:

```Bash
curl -X GET "http://localhost:8080/api/movies?genre=1"
curl -X GET "http://localhost:8080/api/movies?actor=3"
curl -X GET "http://localhost:8080/api/movies?year=2010"
```

### 5. Deletions & Constraints

By default, deleting an entity with existing relationships will return a 400 Bad Request to protect data integrity.

Standard Delete (Will fail if relationships exist ):

```Bash
curl -X DELETE http://localhost:8080/api/genres/1
```

Force Delete (Cascades and removes relationships ):

```Bash
curl -X DELETE "http://localhost:8080/api/genres/1?force=true"
```


## Project Structure

```text

open-movies-db/
├── cmd/
│   └── api/
│       └── main.go           # Entry point: wires up dependencies and starts the server
├── internal/
│   ├── models/               # Struct definitions (Movie, Actor, Genre)
│   ├── database/             # SQLite connection setup (db.go)
│   ├── repository/           # Raw SQL queries (genre.go, movie.go, actor.go)
│   ├── service/              # Logic & validation (genre.go, movie.go, actor.go)
│   └── handler/              # HTTP routing, JSON encoding/decoding (genre.go, movie.go, actor.go)
├── pkg/
│   └── errors/               # Custom error definitions (e.g., ErrNotFound)
├── data/                     # Local folder for the SQLite .db file (ignored in git)
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum

```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.