# Reddit Clone with REST API

This project implements a Reddit-like engine with core functionalities such as account registration, subreddit management, posting, commenting and hierarchical comments. It includes a REST API interface for client interactions and supports simulated user activity for performance testing. 

## Features
1. **Account Management**
   - User registration.
2. **Subreddit Management**
   - Create, join, and leave subreddits.
3. **Post Management**
   - Create and retrieve posts in subreddits.
4. **Commenting**
   - Add hierarchical comments to posts.
5. **Simulation**
   - Simulates user activity, including Zipf distribution for subreddit popularity and reposting.

## REST API Endpoints
The REST API provides the following endpoints for interaction:
- **User Registration**
  - `POST /register`: Register a new user.
- **Subreddit Management**
  - `POST /subreddits`: Create a new subreddit.
  - `POST /subreddits/join`: Join an existing subreddit.
- **Post Management**
  - `POST /posts`: Create a new post in a subreddit.
  - `GET /posts/<subreddit>`: Retrieve posts from a subreddit.
- **Commenting**
  - `POST /comments`: Add a comment to a post.

## How to Run
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/reddit-clone.git
   cd reddit-clone
