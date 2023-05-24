# saas-template

A generic saas backend and frontend used as a template for new projects.

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
- Stories
    - Create - backend only
    - Get Many - backend only
    - Get One - beckend only

### What is this site?

It's a blogging platform.

You can sign up as an individual or create an organization.

Write blog posts. Comment and react to posts. Subscribe to authors.

### Potential features

- Organizations
- Google SSO
- Facebook SSO
- Microsoft SSO
- Password-based login
- Oauth2 for custom apps
- User settings
- User photo
- User bio
- Organization settings
- Organization logo
- Organization info
- Role based access control
- API documentation
- Delete user account
- Delete organization
- Change email address
- Email notifications
- Sign out
- Search
- Subscription settings
- (I should pick a basic use case for the saas template so I can do stuff with data models)

### Testing

In `/backend` run `go test ./...`

This will run unit tests. It will also run an end-to-end test suite using a temporary database. The database is preserved after tests are run allowing it to serve as a source of seed data for manual testing. See `.env` file to configure.

Please add appropriate unit and end-to-end tests for all new features.