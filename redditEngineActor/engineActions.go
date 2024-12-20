package redditEngineActor

import (
	"DOSP4Project/redditMessagesType"
	"fmt"
	"math/rand"
	"time"
)

func (engineVal *Reddit) UserSignUp(userName string) {
	if _, existsVal := engineVal.UserDataVal[userName]; existsVal {
		fmt.Println("The user already exists in the system, %s", userName)
		return
	}

	engineVal.UserDataVal[userName] = redditMessagesType.UserAccountDetails{
		SubRedditArray: []string{},
		UserNameVal:    userName,
	}
	fmt.Println("User %s registered!!\n", userName)
}

func (engineVal *Reddit) SendDirectMessagetoUsers(senderValue string, contentVal string) {
	recipient, err := engineVal.getRandomRecipient(senderValue)
	if err != nil {
		fmt.Println(err)
		return
	}
	directMessagge := redditMessagesType.DirectMessage{
		UserNameVal: senderValue, MessageString: contentVal,
	}

	engineVal.addDMToUser(recipient, directMessagge)
	fmt.Printf("User '%s' is sending Direct Messages to '%s': %s\n", senderValue, recipient, contentVal)
}

func (engineVal *Reddit) ReplyToAllDMs(sender string, replyContent string) {
	senderData, senderExists := engineVal.UserDataVal[sender]
	if !senderExists {
		fmt.Printf("User '%s' does not exist in the system\n", sender)
		return
	}

	if len(senderData.DirectMessages) == 0 {
		fmt.Printf("User '%s' has no direct messages to respond to\n", sender)
		return
	}

	for _, dm := range senderData.DirectMessages {
		recipient := dm.UserNameVal

		replyDM := redditMessagesType.DirectMessage{
			UserNameVal:   sender,
			MessageString: fmt.Sprintf("Respond to the message: %s", dm.MessageString),
		}

		engineVal.addDMToUser(recipient, replyDM)
		fmt.Printf("User '%s' replied to DM from '%s': %s\n", sender, recipient, replyDM.MessageString)
	}

	engineVal.clearUserDMs(sender)
	fmt.Printf("User '%s' has replied to all DMs and cleared their DM list.\n", sender)
}

func (engineVal *Reddit) getRandomRecipient(excludedUser string) (string, error) {
	var potentialRecipients []string
	for userName := range engineVal.UserDataVal {
		if userName != excludedUser {
			potentialRecipients = append(potentialRecipients, userName)
		}
	}

	if len(potentialRecipients) == 0 {
		return "", fmt.Errorf("no users available to send a DM to")
	}

	rand.Seed(rand.Int63()) // Better randomization
	return potentialRecipients[rand.Intn(len(potentialRecipients))], nil
}

func (engineVal *Reddit) addDMToUser(userName string, dm redditMessagesType.DirectMessage) {
	userData, userExists := engineVal.UserDataVal[userName]
	if !userExists {
		fmt.Printf("Recipient '%s' does not exist.\n", userName)
		return
	}

	userData.DirectMessages = append(userData.DirectMessages, dm)
	engineVal.UserDataVal[userName] = userData
}

func (engineVal *Reddit) clearUserDMs(userName string) {
	userData, doesUserExists := engineVal.UserDataVal[userName]
	if !doesUserExists {
		fmt.Printf("User '%s' does not exist.\n", userName)
		return
	}

	userData.DirectMessages = []redditMessagesType.DirectMessage{} // Clear DMs
	engineVal.UserDataVal[userName] = userData
}

func (engineVal *Reddit) GetCommentsForPost(subredditName string, postId int) ([]redditMessagesType.Comments, error) {
	subreddit, exists := engineVal.SubReddits[subredditName]
	if !exists {
		return nil, fmt.Errorf("subreddit '%s' does not exist", subredditName)
	}

	post := findPost(subreddit.ArrayofPostStruct, postId)
	if post == nil {
		return nil, fmt.Errorf("post with ID %d does not exist in subreddit '%s'", postId, subredditName)
	}

	return post.CommentsSection, nil
}

func (engineVal *Reddit) UpvoteRandomPost(userName string) {
	post, subreddit, err := engineVal.selectRandomPost(userName)
	if err != nil {
		fmt.Println(err)
		return
	}

	post.UpVoteVal++
	fmt.Printf("User '%s' upvoted post ID %d in subreddit '%s'. Total upvotes: %d\n",
		userName, post.PostId, subreddit, post.UpVoteVal)

	engineVal.updateKarma(post.UserNameVal, 1)
}

func (engineVal *Reddit) DownvoteRandomPost(userName string) {
	post, subreddit, err := engineVal.selectRandomPost(userName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Downvote the selected post (ensuring it doesn't go below zero)
	if post.UpVoteVal > 0 {
		post.UpVoteVal--
		fmt.Printf("User '%s' downvoted post ID %d in subreddit '%s'. Total upvotes: %d\n",
			userName, post.PostId, subreddit, post.UpVoteVal)

		// Update karma points for the post's creator
		engineVal.updateKarma(post.UserNameVal, -1)
	} else {
		fmt.Printf("User '%s' attempted to downvote post ID %d in subreddit '%s', but upvotes are already at zero.\n",
			userName, post.PostId, subreddit)
	}
}

func (engineVal *Reddit) selectRandomPost(userName string) (*redditMessagesType.PostStruct, string, error) {
	user, userExists := engineVal.UserDataVal[userName]
	if !userExists {
		return nil, "", fmt.Errorf("user '%s' is not registered and cannot vote on posts", userName)
	}

	if len(user.SubRedditArray) == 0 {
		return nil, "", fmt.Errorf("user '%s' is not subscribed to any subreddits", userName)
	}

	rand.Seed(time.Now().UnixNano())
	selectedSubReddit := user.SubRedditArray[rand.Intn(len(user.SubRedditArray))]

	subreddit, subredditExists := engineVal.SubReddits[selectedSubReddit]
	if !subredditExists || len(subreddit.ArrayofPostStruct) == 0 {
		return nil, "", fmt.Errorf("subreddit '%s' does not exist or has no posts", selectedSubReddit)
	}

	selectedPost := &subreddit.ArrayofPostStruct[rand.Intn(len(subreddit.ArrayofPostStruct))]
	return selectedPost, selectedSubReddit, nil
}

func (engineVal *Reddit) updateKarma(userName string, delta int) {
	user, exists := engineVal.UserDataVal[userName]
	if !exists {
		return
	}

	user.KarmaPointsVal += delta
	if user.KarmaPointsVal < 0 {
		user.KarmaPointsVal = 0 // Ensure karma doesn't go negative
	}

	engineVal.UserDataVal[userName] = user
	fmt.Printf("Updated karma for user '%s'. New karma: %d\n", userName, user.KarmaPointsVal)
}

func (engineVal *Reddit) SubredditSpecificOp(action, userName, subRedditName string) {
	userData, userExists := engineVal.UserDataVal[userName]
	if !userExists {
		fmt.Printf("User %s is not registered and cannot join or leave subreddits.\n", userName)
		return
	}

	if _, exists := engineVal.SubReddits[subRedditName]; !exists {
		engineVal.SubReddits[subRedditName] = redditMessagesType.SubRedditStruct{ArrayofPostStruct: []redditMessagesType.PostStruct{}}
		fmt.Printf("Subreddit %s did not exist and was created.\n", subRedditName)
	}

	switch action {
	case "join":
		if isUserInSubreddit(userData.SubRedditArray, subRedditName) {
			fmt.Printf("User %s is already a member of %s subreddit.\n", userName, subRedditName)
			return
		}
		userData.SubRedditArray = append(userData.SubRedditArray, subRedditName)
		engineVal.UserDataVal[userName] = userData
		fmt.Printf("User %s joined %s subreddit.\n", userName, subRedditName)

	case "leave":
		if !isUserInSubreddit(userData.SubRedditArray, subRedditName) {
			fmt.Printf("User %s is not a member of %s subreddit.\n", userName, subRedditName)
			return
		}
		userData.SubRedditArray = removeSubreddit(userData.SubRedditArray, subRedditName)
		engineVal.UserDataVal[userName] = userData
		fmt.Printf("User %s left %s subreddit.\n", userName, subRedditName)

	default:
		fmt.Printf("Unknown action: %s. Supported actions are 'join' and 'leave'.\n", action)
	}
}

func (engineVal *Reddit) CreatePost(userName, subRedditName string, content redditMessagesType.PostStruct) {
	if _, userExists := engineVal.UserDataVal[userName]; !userExists {
		fmt.Printf("User %s is not registered and cannot create posts.\n", userName)
		return
	}

	subreddit, subredditExists := engineVal.SubReddits[subRedditName]
	if !subredditExists {
		fmt.Printf("Subreddit %s does not exist.\n", subRedditName)
		return
	}

	subreddit.ArrayofPostStruct = append(subreddit.ArrayofPostStruct, content)
	engineVal.SubReddits[subRedditName] = subreddit
	fmt.Printf("User %s created a post in subreddit %s: %s\n", userName, subRedditName, content.PostContent)

	if !isUserInSubreddit(engineVal.UserDataVal[userName].SubRedditArray, subRedditName) {
		engineVal.SubredditSpecificOp("join", userName, subRedditName)
	}
}

func (engineVal *Reddit) ReplyToComment(userName, subRedditName string, postID, commentID int, replyContent string) {
	// Check if the user exists
	if _, userExists := engineVal.UserDataVal[userName]; !userExists {
		fmt.Printf("User %s is not registered and cannot reply to comments.\n", userName)
		return
	}

	// Check if the subreddit exists
	subreddit, subredditExists := engineVal.SubReddits[subRedditName]
	if !subredditExists {
		fmt.Printf("Subreddit %s does not exist.\n", subRedditName)
		return
	}

	// Check if the post exists
	post := findPost(subreddit.ArrayofPostStruct, postID)
	if post == nil {
		fmt.Printf("Post with ID %d does not exist in subreddit %s.\n", postID, subRedditName)
		return
	}

	// Create a new comment or reply
	reply := redditMessagesType.Comments{
		CommentId:           generateSequentialCommentId(post.CommentsSection),
		CommentString:       replyContent,
		UserNameVal:         userName,
		UpVoteVal:           0,
		DownVoteVal:         0,
		ReplytoCommentArray: []redditMessagesType.Comments{},
	}

	if commentID > 0 {
		// Try to add the reply to an existing comment
		if addReplyToComment(post.CommentsSection, commentID, reply) {
			fmt.Printf("User %s replied to comment %d in post %d in subreddit %s.\n", userName, commentID, postID, subRedditName)
		} else {
			fmt.Printf("Comment ID %d not found. Cannot add reply.\n", commentID)
			return
		}
	} else {
		// Add as a top-level comment
		post.CommentsSection = append(post.CommentsSection, reply)
		fmt.Printf("User %s added a new top-level comment to post %d in subreddit %s.\n", userName, postID, subRedditName)
	}

	// Save the updated post back into the subreddit
	for i, p := range subreddit.ArrayofPostStruct {
		if p.PostId == postID {
			subreddit.ArrayofPostStruct[i] = *post
			break
		}
	}
	engineVal.SubReddits[subRedditName] = subreddit
}

func generateSequentialCommentId(commentsArray []redditMessagesType.Comments) int {
	maxId := 0
	for _, comment := range commentsArray {
		if comment.CommentId > maxId {
			maxId = comment.CommentId
		}
		// Check replies recursively
		maxReplyId := generateSequentialCommentId(comment.ReplytoCommentArray)
		if maxReplyId > maxId {
			maxId = maxReplyId
		}
	}
	return maxId + 1
}

func isUserInSubreddit(subreddits []string, subredditName string) bool {
	for _, sub := range subreddits {
		if sub == subredditName {
			return true
		}
	}
	return false
}

func removeSubreddit(subredditsArray []string, subredditName string) []string {
	for i, sub := range subredditsArray {
		if sub == subredditName {
			return append(subredditsArray[:i], subredditsArray[i+1:]...)
		}
	}
	return subredditsArray
}

func findPost(posts []redditMessagesType.PostStruct, postID int) *redditMessagesType.PostStruct {
	for i := range posts {
		if posts[i].PostId == postID {

			println("Post is found")
			return &posts[i]
		}
	}
	return nil
}

func addReplyToComment(commentsArray []redditMessagesType.Comments, targetID int, reply redditMessagesType.Comments) bool {
	for i := range commentsArray {
		if commentsArray[i].CommentId == targetID {
			// Found the target comment; add reply to its ReplytoCommentArray
			reply.CommentId = generateSequentialCommentId(commentsArray[i].ReplytoCommentArray)
			commentsArray[i].ReplytoCommentArray = append(commentsArray[i].ReplytoCommentArray, reply)
			return true
		}
		// Recursively check replies to the current comment
		if addReplyToComment(commentsArray[i].ReplytoCommentArray, targetID, reply) {
			return true
		}
	}
	return false
}
