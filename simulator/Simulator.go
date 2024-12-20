package simulator

import (
	"DOSP4Project/redditMessagesType"
	"github.com/asynkron/protoactor-go/actor"
	"math/rand"
	"sync"
	"time"
)

func randomSubredditGenerator() string {
	subreddits := []string{"Scala", "Machine Learning", "Football", "Elon Musk", "Laptop"}
	return subreddits[rand.Intn(len(subreddits))]

}

func randomReplyGenerator() string {
	replies := []string{
		"I completely agree with you.",
		"That's an intriguing perspective.",
		"I appreciate you sharing this.",
		"Could you provide more details?",
		"This is truly thought-provoking.",
	}
	return replies[rand.Intn(len(replies))]
}

func randomContentGenerator() string {
	contents := []string{
		"This post is amazing",
		"What a fantastic post",
		"I absolutely love this post",
		"This is such a brilliant post!",
		"What a great piece of content!",
	}
	return contents[rand.Intn(len(contents))]
}

func SimulateUserActivity(enginePIDVal *actor.PID, username string, actoropool *actor.ActorSystem, wg *sync.WaitGroup) {

	defer wg.Done()

	rand.Seed(time.Now().UnixNano())

	actoropool.Root.Send(enginePIDVal, &redditMessagesType.SignUpUser{UserNameVal: username})

	randomUserActions := []func(){
		func() {
			subreddit := randomSubredditGenerator()
			actoropool.Root.Send(enginePIDVal, &redditMessagesType.AddUsertoSubRedditStruct{UserNameVal: username, SubredditName: subreddit})
		},
		func() {
			subreddit := randomSubredditGenerator()
			post := redditMessagesType.PostStruct{
				PostId:          rand.Intn(1000),
				PostContent:     randomContentGenerator(),
				UserNameVal:     username,
				UpVoteVal:       0,
				DownVoteVal:     0,
				CommentsSection: []redditMessagesType.Comments{},
			}
			actoropool.Root.Send(enginePIDVal, &redditMessagesType.NewPostStruct{UserNameVal: username, SubRedditName: subreddit, ContentVar: post})
		},
		func() {
			subreddit := randomSubredditGenerator()
			actoropool.Root.Send(enginePIDVal, &redditMessagesType.ReplytoCommentStruct{
				UserNameVal:      username,
				SubRedditNameVal: subreddit,
				PostId:           rand.Intn(10) + 1,
				CommentId:        rand.Intn(10) + 1,
				ReplyString:      randomReplyGenerator(),
			})
		},
	}

	performRandomUserActions(randomUserActions)
}

func performRandomUserActions(randomUserActions []func()) {
	for i := 0; i < 10; i++ {
		action := randomUserActions[rand.Intn(len(randomUserActions))]
		action()
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}
