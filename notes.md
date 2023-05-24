# Development Notes



Architecture problems:

- Error management
- Make auth more modular
- How to build with access control in mind?
- How to keep features modular?
- How to spin up features more quickly?



What are the primary user journeys?
- I want to write something without logging in. (We don't need an API for this, since it can just be cached in the browser - yes. Even if we stored it in the cloud, they would lose access when clearing their cookies any way)

- I want to access public content without logging in. This requires APIs.
- I want to change the access level of something I write. (From public to private, vice versa)
- I want to make it so only a custom set of people can access this document.
- I want to invite other users to edit this document
- I want to invite other users to comment on this document
- I want to invite other users to do X action.


- I want to view the resources that are accessible to me