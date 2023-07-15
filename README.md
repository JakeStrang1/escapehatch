# escapehatch

### Try it
The project is alive and well at [escapehatch.ca](https://escapehatch.ca)!

### 
Escapehatch is a platform where users track and share their favorite books and shows. The site is built and operated by Jake Strang and launched for early users in June 2023, and is in active development towards our v1 milestone.

### The Ideal User
You love reading, watching shows, and diving into the worlds that others have created. You feel a sense of achievement when you finish a new book or show. You feel like the content you consume says something about who you are and you want to showcase this to others. You can think of recent conversations where having a list of your favourites on hand would have come in handy. You want to see your achievements and progress laid out in front of you, especially if you can organize them. You like making lists and the time you spend organizing your hobbies feels like a hobby itself. You aren't daunted that most of your content needs to be added manually by you to the site, in fact you feel excited. You're a completionist and the goal of being the one to contribute entries that others will see and find useful is a thrill. You could sit for hours just filling in missing content on the site in order to be a top contributor.

## Roadmap
**Current phase**: Building v1

### V1 - The bare minimum feature set for the ideal user to find value

- [x] Passwordless login
- [x] Account setup wizard
- [x] Profile page displays user's shelves (Movies, TV Shows, Books)
- [x] Search page to find items to add to shelves
- [x] Followers page to view followers, following, and search for users
- [ ] Item details page with additional action buttons (remove from shelf)
- [x] Form to add new items to the DB
- [ ] Form to edit or remove item from DB

### V2 - Features to increase delight and stickiness for the ideal user

- [ ] Create/edit/delete custom shelves
- [ ] Set additional item attributes (read/unread, progress tracking, dates) 
- [ ] View shelf as full page
- [ ] Rearrange items on shelves (drag)
- [ ] New follower notification
- [ ] Publish leaderboard of user contributions
...

## Development

### Setup

Unless you're me, I don't really recommend trying to run this locally. Usually this repo is private, but I've opened it up temporarily to showcase the code. But if you're super determined to run it yourself here are the steps that should get you most of the way.

1. Make sure you have Go, MongoDB and npm installed.
2. Create `/backend/.env` file

My `/backend/.env` file looks like this:

```
#!/bin/bash

export MONGO_DB_NAME="escapehatch"
export MONGO_HOST="mongodb://localhost:27017"
export ORIGIN="http://localhost:3000"
export USE_SENDGRID="true" # I still use the live emails when running locally, but try setting this to false to avoid needing the Sendgrid integration
export SENDGRID_API_KEY="<Your own Sendgrid API Key>"
export SENDGRID_FROM_EMAIL="<Your own Sendgrid From email>"
export FRONTEND_HOST="http://localhost:3000"
export PRODUCTION="false"
export USE_GCS="false"
export GCS_BUCKET_NAME=""
export STATIC_URL_ROOT="http://localhost:8080/local-static"
export USE_ATLAS_SEARCH="false"
export TEST_EMAIL="<Your personal email>" # Test accounts will be generated that you can log into
export TEST_MONGO_DB_NAME="escapehatch-test"
```

3. Create `/frontend/.env.local` file

My `frontend/.env.local` looks like this:

```
REACT_APP_API_HOST="http://localhost:8080"
```

### Run locally

In `/backend`, run `go run main.go` to start the API server.

In `/frontend`, run `npm start` to start the frontend on `localhost:3000`.

Note, there are some issues running the frontend on higher versions of node. If you have problems try `nvm use 15.14.0`, may need to use an elevated console.

See `/backend/.env` and `/frontend/.env.local` to set environment variables.

### Testing

In `/backend` run `go test ./...`

This will run unit tests. It will also run an end-to-end test suite using a temporary database. The database is preserved after tests are run allowing it to serve as a source of seed data for manual testing. See `.env` file to configure.

Please add appropriate unit and end-to-end tests for all new features.
