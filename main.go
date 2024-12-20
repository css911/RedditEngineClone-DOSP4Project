package main

import (
	"DOSP4Project/redditEngineActor"
	"DOSP4Project/redditMessagesType"
	"encoding/json"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Engine *redditEngineActor.Reddit
	WG     *sync.WaitGroup
}

func main() {
	app := &App{
		Router: mux.NewRouter(),
		WG:     &sync.WaitGroup{},
	}

	protoactorPool := actor.NewActorSystem()
	app.Engine = redditEngineActor.NewEngineActor(protoactorPool, app.WG)

	app.initializeRoutes()

	fmt.Println("Starting Reddit Simulation API on port 8080...")
	http.ListenAndServe(":8080", app.Router)
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/users", app.createUser).Methods("POST")
	app.Router.HandleFunc("/users/{userName}/sendDM", app.sendDirectMessage).Methods("POST")
	app.Router.HandleFunc("/users/{userName}/replyDMs", app.replyToDirectMessages).Methods("POST")
	app.Router.HandleFunc("/users/{userName}/upvote", app.upvotePost).Methods("POST")
	app.Router.HandleFunc("/users/{userName}/downvote", app.downvotePost).Methods("POST")
	app.Router.HandleFunc("/subreddits/{subredditName}/leave", app.leaveSubreddit).Methods("POST")
	app.Router.HandleFunc("/subreddits", app.createSubreddit).Methods("POST")
	app.Router.HandleFunc("/subreddits/{subredditName}/join", app.joinSubreddit).Methods("POST")
	app.Router.HandleFunc("/subreddits/{subredditName}/posts", app.createPost).Methods("POST")
	app.Router.HandleFunc("/subreddits/{subredditName}/posts/{postId}/comments/{commentId}/reply", app.replyToComment).Methods("POST")
	app.Router.HandleFunc("/subreddits/{subredditName}/posts/{postId}/comments", app.getCommentsForPost).Methods("GET")
	app.Router.HandleFunc("/subreddits/{subredditName}/posts", app.getAllPostsFromSubreddit).Methods("GET")

}

func (app *App) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserName string `json:"userName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Engine.UserSignUp(req.UserName)
	fmt.Printf("[LOG] User %s created successfully\n", req.UserName)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s created\n\n", req.UserName)
}

func (app *App) getCommentsForPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subredditName := vars["subredditName"]
	postId := vars["postId"]

	postIDInt, err := strconv.Atoi(postId)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Fetch comments from the Engine
	comments, err := app.Engine.GetCommentsForPost(subredditName, postIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return comments in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (app *App) getAllPostsFromSubreddit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subredditName := vars["subredditName"]

	subreddit, exists := app.Engine.SubReddits[subredditName]
	if !exists {
		http.Error(w, fmt.Sprintf("Subreddit '%s' does not exist", subredditName), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subreddit.ArrayofPostStruct)
}

func (app *App) leaveSubreddit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subredditName := vars["subredditName"]

	var req struct {
		UserName string `json:"userName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Engine.SubredditSpecificOp("leave", req.UserName, subredditName)
	fmt.Printf("[LOG] User %s left subreddit %s successfully\n", req.UserName, subredditName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s left subreddit %s\n\n", req.UserName, subredditName)
}

func (app *App) createSubreddit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubredditName string `json:"subredditName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Engine.SubredditSpecificOp("create", "", req.SubredditName)
	fmt.Printf("[LOG] Subreddit %s created successfully\n", req.SubredditName)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Subreddit %s created\n\n", req.SubredditName)
}

func (app *App) joinSubreddit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subredditName := vars["subredditName"]

	var req struct {
		UserName string `json:"userName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Engine.SubredditSpecificOp("join", req.UserName, subredditName)
	fmt.Printf("[LOG] User %s joined subreddit %s successfully\n", req.UserName, subredditName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s joined subreddit %s\n\n", req.UserName, subredditName)
}

func (app *App) sendDirectMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Engine.SendDirectMessagetoUsers(userName, req.Content)
	fmt.Printf("[LOG] Direct message sent by %s with content: %s\n", userName, req.Content)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Direct message sent by %s\n\n", userName)
}

func (app *App) replyToDirectMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]

	var req struct {
		ReplyMessage string `json:"replyMessage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Engine.ReplyToAllDMs(userName, req.ReplyMessage)
	fmt.Printf("[LOG] User %s replied to all DMs with message: %s\n", userName, req.ReplyMessage)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All direct messages replied by %s\n\n", userName)
}

func (app *App) upvotePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]

	app.Engine.UpvoteRandomPost(userName)
	fmt.Printf("[LOG] User %s upvoted a random post\n", userName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Post upvoted by %s\n\n", userName)
}

func (app *App) downvotePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]

	app.Engine.DownvoteRandomPost(userName)
	fmt.Printf("[LOG] User %s downvoted a random post\n", userName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Post downvoted by %s\n\n", userName)
}

func (app *App) createPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subredditName := vars["subredditName"]

	var req struct {
		UserName string                        `json:"userName"`
		Content  redditMessagesType.PostStruct `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique postId for the new post
	newPostId := rand.Intn(1000000) // Ensure uniqueness in production with better methods
	req.Content.PostId = newPostId

	app.Engine.CreatePost(req.UserName, subredditName, req.Content)

	fmt.Printf("[LOG] User %s created a post in subreddit %s with content: %s and with postID: %s\n", req.UserName, subredditName, req.Content.PostContent, newPostId)
	w.WriteHeader(http.StatusCreated)

	// Return the generated postId in the response
	response := map[string]interface{}{
		"message": "Post created successfully",
		"postId":  newPostId,
	}
	json.NewEncoder(w).Encode(response)
}

func (app *App) replyToComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subredditName := vars["subredditName"]
	postId := vars["postId"]
	commentId := vars["commentId"]

	var req struct {
		UserName string `json:"userName"`
		Reply    string `json:"reply"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postIDInt, err := strconv.Atoi(postId)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	commentIDInt, err := strconv.Atoi(commentId)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	app.Engine.ReplyToComment(req.UserName, subredditName, postIDInt, commentIDInt, req.Reply)
	fmt.Printf("[LOG] User %s replied to comment %s in post %s in subreddit %s with reply: %s\n", req.UserName, commentId, postId, subredditName, req.Reply)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Reply added to comment %s in post %s in subreddit %s\n\n", commentId, postId, subredditName)
}
