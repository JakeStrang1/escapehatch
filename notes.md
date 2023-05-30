# Development Notes

## Vision

This is your go-to place to showcase "your favorites" when it comes to books, shows, movies, bands, podcasts, youtubers, and all other sorts of media.

Check out what your friends like, what your crush likes, discover new creators you'll love, and share your latest binge sessions.

## Potential features

- Have an enormous database of content (books, shows, movies, artists) - these will be user-driven for now.
- Integrate with lots of APIs to pull in data like Spotify, Netflix, Goodreads, to pull in data.
- Have a ton of content types:
    - Books
    - TV Shows
    - Movies
    - Music
    - Podcasts
    - Youtubers
    - Subreddits
    - Video games
    - Board games
    - Blogs?
    - Be able to add more content types as needed
- Users add content they like to their profile (a.k.a. escape pod?)
- Users can follow other users, this means their posts will show up in their feed
- Users can connect Facebook to help find existing connections
- Users can connect Instagram to help find existing connections
- Users can have a "Top 5" on their profile, which is their lifetime top 5 all time favourite anything
- Users can post certain types of activity, but I think all the posts will be presets, they won't be able to type any text. I don't think I want the platform to be about sharing messages
- Users have an activity feed (or a "passivity feed" perhaps) that shows posts and other time-based content.
    - Posts from people you are following
    - "Suggested for you" posts from people you aren't following
    - Targeted ads eventually?
- Users can react to posts from other users (a simple "like" or heart sort of reaction). I don't think there will be any comment feature though. Maybe?
- One type of post will be what they're watching "right now" (e.g. show status for the next hour or until changed)
- One type of post will be when a user adds a new content to their profile, they can optionally share it as a post. If its a book they can specify "reading for the first time" or "completed". Same with TV series, maybe other types of series. Not so much for movies or music.
- Users can indicate which content they spent time on today. This wouldn't really be a post, more of a "activity tracker"
- For books/TV shows, users should be able to also update their profile with how far along they are (which page or episode).
- Users can see their own/other's "recent activity" aggregated in some form (e.g. what is that user watching/reading lately)
- If a user doesn't find the book/show they're looking for in the search, then they can add it to the database for all users to see. I think each piece of content needs to have an image and some basic info. We should keep track of which users added which content so we can do badges or something.
- Users can enable notifications for a person/interest combo based on the "right now" feature. That way if that person says their watching The Office, then you'll get notified and you could reach out - maybe you want to watch with them or something? The notification should disappear if it isn't seen within 30 minutes or something.
- We want to give general insights: "a lot of people with similar interests as you are recently watching X"
- I do think we probably need a chat feature. Users won't comment on posts, but they can kick off a chat message (similar to commenting on IG stories)
- User notifications


## UI Planning

The difficulty of designing a UI for a young project is that it takes so much effort to move things around and change it up, but each new feature requires exactly that.

Iterative approach:
1. Select a batch of features to include in this iteration - batching features allows for the UI to evolve intentionally
2. Design version of the UI on paper.
3. Outline every API endpoint needed based on the design.
4. Develop/release backend as features are ready.
5. Develop frontend all at once after the backend is done.
6. Release the new version of the UI.

## Iteration 1: MVP

### Core features

1. Add content to your personal library
2. View your personal library
3. View the libraries of other users
4. Follow other users 

### Supporting features

5. Sign up
6. Log in
7. Log out
8. Search for content
9. Search for users
10. List followers / following
11. Update registry (add/update/delete)
12. Remove content from your personal library
13. View content info page

### UI capabilities

- search for content
- search for users
- update registry
- add content to library
- view library
- view libraries of other users

### Views

- User profile (including your own)
- Search page (content, maybe also people)
- Followers (shows followers, following, and search for users)
- Add/Edit content form
- Content info (image, description, stats, actions)

### Default view: your library

#### Information

- Username
- Real name
- User number
- Number of followers
- Number of following
- Shelves
- Number of entries per shelf
- Images of entries on each shelf

#### Links

- Search view
- Followers view
- Logout
- Follow / Unfollow
- Link to each content

### Search View

#### Information

- Results
- Result image
- Result media type
- Result title (includes series, #)
- Publication year (or year range for shows, where "present" could be an option)
- Author (or main 2 actors for movies / TV)
- Number of people w/ this media on their shelf
- If its on your shelf or not

#### Links

- Home view
- Followers view
- Submit search
- Link to each content
- Add entry to shelf
- "Add new entry" link

### Followers View

#### Information

- Home view
- Search view
- Number of followers
- Number of following
- List of all followers/following
- Result username
- Result real name
- Result number of mutual following (for find page only)
- Result - do you follow them?
- Result - do they follow you?

#### Links

- Submit search (for followers, following, and find new)
- Followers tab
- Following tab
- Find new tab
- Follow user
- Unfollow user
- Remove user
- User library

### Add/Edit Content Form

#### Information

- "Add new content" vs. "Edit content" (Title)
- Media Type (dropdown: book, movie, TV series)
- Cover Image
- Book: Author
- Book: Published year
- Book: Title
- Book: Is Series?
- Book: Series title
- Book: # in Series
- Movie: List at least 2 main actors
- Movie: Release year
- Movie: Title
- TV series: List at least 2 main actors
- TV series: Years running ( ____ to ____) can be "present"
- TV series: Title
- Form error displayed (e.g. missing required field)

#### Links

- Home view (confirm leave page)
- Search view (confirm leave page)
- Followers view (confirm leave page)
- Delete from archive (edit page only) - provide reason for request
- Upload image button (from file and from link)
- Save/Submit
- Cancel

### Content Info View

#### Information

- Title
- Image
- Media Type
- Year
- Author (if applicable)
- Added by X users
- Added by you?

#### Links

- Home view
- Search view
- Followers view
- Add / Remove from personal library
- Edit content

### API Endpoints

- [x] `GET /users`
- [x] `GET /users/:id`
- [x] `GET /users/:id/followers`
- [ ] `GET /users/:id/following`
- [x] `POST /users/:id/follow`
- [ ] `POST /users/:id/unfollow`
- [ ] `POST /users/:id/remove`
- [ ] `POST /items`
- [ ] `GET /items/:id`
- [ ] `PATCH /items/:id`
- [ ] `POST /items/:id/add`
- [ ] `POST /items/:id/remove`
- [ ] `POST /items/:id/delete`
- [ ] `GET /search`