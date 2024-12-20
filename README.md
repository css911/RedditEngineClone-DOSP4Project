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
2. Start the server:
   ```bash
   go run main.go
3. To Check the server, start client(Client Code is available in a parallel repo in github)
   ```bash
   go run client.go

##Project Structure
  redditMessagesType.go: Defines the structures for messages between actors.
  redditEngineActor.go: Contains the core engine logic using actor-based concurrency.
  simulator.go: Simulates user activity on the platform.
  main.go: Sets up the REST API server and routes HTTP requests.

##Performance
   The system scales efficiently with user activity. Below are some performance metrics:

   Number of Users	Number of Activities	Time Taken (seconds)
   20	20	6.41
   1000	1000	256
   2000	2000	522
A Zipf distribution model is used to simulate subreddit popularity.

##Demo Video
Watch the project in action: YouTube Link

Contributors
Chetan Shinde (UFID: 92044184)
Manoj Deo (UFID: 83224893)
License
This project is licensed under the MIT License. See the LICENSE file for details.
