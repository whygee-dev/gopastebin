# GoLang Paste Bin REST API

This project provides a simple Paste Bin REST API implemented in GoLang using Mux for routing and SQLite3 for data storage.

## Features

- Signup
- Login
- Create a new paste with a given content and expiry time.
- Update an existing paste with new content.
- Retrieve an existing paste by its ID, with click count.
- Retrieve global stats

## Routes

- Signup: POST /user/signup [email, password]
- Login: POST /user/login [email, password]

[Authenticated routes, require a valid JWT token in the Authorization header, e.g. `Bearer <token>`]

- Create a new paste: PUT /paste/create => {content, expiry}
- Update an existing paste: PATCH /paste/update => {id, content}
- Retrieve an existing paste: GET /paste/{id}
- Retrieve global stats: GET /stats

## Installation

To run this project locally, ensure you have Go installed on your machine.

1. Clone the repository:

   ```
   git clone https://github.com/gopastebin/gopastebin.git
   cd gopastebin
   ```

2. Run the server:

   ```
   docker compose up
   ```

The server will start running on `http://localhost:3333`.

## Testing

To run tests:

```
go test ./...
```

Coverage:

```
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

or

```
sh test.sh
```
