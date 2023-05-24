# app package

This is the entry point for the app.

## Files

- `app.go` has the constructor for the app, and defines setup and tear down of app dependencies.
- `router.go` defines the API routes and middleware for the entire app.

## How to use

This package is for the most high-level app behavior. Most business logic should happen elsewhere. If the app were a book, this package would be the front matter and back matter (e.g. preface, table of contents, bibliography), rather than the principal text. As you add new sections (features) to the book, you probably need to update the table of contents (e.g. the API router), but you likely won't be changing more than a couple lines at a time.

