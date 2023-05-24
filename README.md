# escapehatch

A platform where users track and share their favorite books and shows.

### Usage

In `/backend`, run `go run main.go` to start the API server.

In `/frontend`, run `npm start` to start the frontend on `localhost:3000`.

Note, there are some issues running the frontend on higher versions of node. If you have problems try `nvm use 15.14.0`, may need to use an elevated console.

See `/backend/.env` and `/frontend/.env` to set environment variables.

### Finished features

- Sign up / Sign in flows
- Passwordless authentication
- Email abuse reporting and IP flagging
- Backend test suite
- Post API
- Comment API

### Feature roadmap

- Remove comment API lol
- Add book content type
- Add TV show content type
- Add movie content type

### Testing

In `/backend` run `go test ./...`

This will run unit tests. It will also run an end-to-end test suite using a temporary database. The database is preserved after tests are run allowing it to serve as a source of seed data for manual testing. See `.env` file to configure.

Please add appropriate unit and end-to-end tests for all new features.