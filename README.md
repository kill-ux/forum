# Web Forum Project

## Objectives
This project aims to develop a web forum application with the following features:

- **Communication Between Users:** Users can create posts and comments to interact with each other.
- **Categories:** Posts can be associated with one or more categories.
- **Likes and Dislikes:** Users can like or dislike posts and comments.
- **Filtering:** Posts can be filtered based on categories, created posts, or liked posts.


### SQLite Overview:
SQLite is a lightweight embedded database often used for local storage in application software. You can learn more about SQLite on the [SQLite official page](https://www.sqlite.org/).

## Authentication
The forum must allow users to:

1. **Register:**
   - Users must provide:
     - Email (must be unique; return an error if already taken).
     - Username.
     - Password.
   - Password encryption is a bonus task.

2. **Login:**
   - Users must authenticate with their credentials to access forum features.
   - Validate credentials against the database (email and password).
   - Return an error if authentication fails.

3. **Sessions:**
   - Use cookies to maintain user sessions.
   - Sessions must have an expiration date.
   - Implement UUID for session management as a bonus task.

## Communication
To enable communication among users:

- Only **registered users** can create posts and comments.
- Posts can be tagged with one or more categories. The choice and implementation of categories are up to you.
- **Visibility:**
  - Posts and comments are visible to all users, including non-registered users.
  - Non-registered users can only view posts and comments but cannot interact (e.g., create or like).

## Likes and Dislikes

- Only **registered users** can like or dislike posts and comments.
- The number of likes and dislikes must be visible to all users, regardless of registration status.

## Filtering

Implement a filtering mechanism for posts based on:

1. **Categories:** Display posts by their associated category (similar to subforums).
2. **Created Posts:** Show posts created by the logged-in user.
3. **Liked Posts:** Display posts liked by the logged-in user.

> Note: The last two filters are only available to logged-in users.

## Getting Started

### Prerequisites:
- [Docker](https://www.docker.com/) installed.
- go.

### Setup:
1. Clone the repository.
    ```bash
        git clone https://learn.zone01oujda.ma/git/muboutoub/forum.git
    ``` 
2. Build the Docker image:
   ```bash
   docker build -t forum-app .
   ```
3. Run the application:
   ```bash
   docker run -p 8080:8080 forum-app
   ```
4. Access the forum in your web browser at `http://localhost:8080`.

---

