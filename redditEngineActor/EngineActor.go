package redditEngineActor

import (
	"DOSP4Project/redditMessagesType"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"sync"
)

type Reddit struct {
	UserDataVal map[string]redditMessagesType.UserAccountDetails
	waitgroup   *sync.WaitGroup
	sys         *actor.ActorSystem
	SubReddits  map[string]redditMessagesType.SubRedditStruct
}

func NewEngineActor(systemActor *actor.ActorSystem, wg *sync.WaitGroup) *Reddit {
	fmt.Println("New Reddit Engine Actor has been initialized.")
	return &Reddit{
		UserDataVal: make(map[string]redditMessagesType.UserAccountDetails),
		waitgroup:   wg,
		sys:         systemActor,
		SubReddits:  make(map[string]redditMessagesType.SubRedditStruct),
	}
}

func (engineVal *Reddit) Receive(ctx actor.Context) {
	fmt.Println("New Reddit Engine Actor initialized; Action to be performed: %v", ctx.Message())
	switch messageVar := ctx.Message().(type) {
	case *redditMessagesType.SignUpUser:
		fmt.Println("Reddit: Doing User Sign Up")
		engineVal.UserSignUp(messageVar.UserNameVal)

	case *redditMessagesType.IncreaseVote:
		fmt.Println("Reddit: Increasing the vote count.\n")
		engineVal.UpvoteRandomPost(messageVar.UserNameVal)

	case *redditMessagesType.NewPostStruct:
		fmt.Println("Reddit: New post creation in progress.")
		engineVal.CreatePost(messageVar.UserNameVal, messageVar.SubRedditName, messageVar.ContentVar)

	case *redditMessagesType.AddUsertoSubRedditStruct:
		fmt.Println("Reddit: Adding user to subreddit.")
		engineVal.SubredditSpecificOp("join", messageVar.UserNameVal, messageVar.SubredditName)

	case *redditMessagesType.RemoveUserFromSubRedditStruct:
		fmt.Println("Reddit: Operation to remove user from subreddit in progress.\n")
		engineVal.SubredditSpecificOp("leave", messageVar.UserNameVal, messageVar.SubRedditNameVal)

	case *redditMessagesType.ReplytoCommentStruct:
		fmt.Println("Reddit: Executing the comment reply operation.")
		engineVal.ReplyToComment(messageVar.UserNameVal, messageVar.SubRedditNameVal, messageVar.PostId, messageVar.CommentId, messageVar.ReplyString)

	case *redditMessagesType.SendDirectMessagetoUserStruct:
		fmt.Println("Reddit: Initiating direct message to the user.")
		engineVal.SendDirectMessagetoUsers(messageVar.UserNameVal, messageVar.ContentVar)

	case *redditMessagesType.ReplyToDirectMessage:
		fmt.Println("Reddit: Replying to Direct Message Operation")
		engineVal.ReplyToAllDMs(messageVar.UserNameVal, messageVar.ReplyMessage)

	case *redditMessagesType.DecreaseVote:
		fmt.Println("Reddit: Decrementing the vote count by 1.")
		engineVal.DownvoteRandomPost(messageVar.UserNameVal)

	default:
		fmt.Println("Reddit: An unrecognized messageVar type was received.")
	}
	fmt.Println("Reddit: Operation execution Completed\n\n\n")
}
