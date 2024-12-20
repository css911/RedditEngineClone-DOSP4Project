package redditMessagesType

type DecreaseVote struct {
	UserId      int
	UserNameVal string
}

type IncreaseVote struct {
	UserId      int
	UserNameVal string
}

type DirectMessage struct {
	UserId        int
	UserNameVal   string
	MessageString string
}

type SignUpUser struct {
	UserId         int
	UserNameVal    string
	SubRedditArray []string
}

type UserAccountDetails struct {
	UserId         int
	DirectMessages []DirectMessage
	KarmaPointsVal int
	UserNameVal    string
	SubRedditArray []string
}

type PostStruct struct {
	PostId          int
	PostContent     string
	UserNameVal     string
	UpVoteVal       int
	DownVoteVal     int
	CommentsSection []Comments
}

type Comments struct {
	CommentId           int
	CommentString       string
	UserNameVal         string
	UpVoteVal           int
	DownVoteVal         int
	ReplytoCommentArray []Comments
}

type SubRedditStruct struct {
	ArrayofPostStruct []PostStruct
}

type ReplyToDirectMessage struct {
	UserId       int
	UserNameVal  string
	ReplyMessage string
}

type ReplytoCommentStruct struct {
	CommentId        int
	UserId           int
	UserNameVal      string
	ReplyString      string
	PostId           int
	SubRedditNameVal string
}

type newSubRedditStruct struct {
	SubRedditName string
}

type AddUsertoSubRedditStruct struct {
	UserId        int
	UserNameVal   string
	SubredditName string
}

type RemoveUserFromSubRedditStruct struct {
	UserId           int
	SubRedditNameVal string
	UserNameVal      string
}

type NewPostStruct struct {
	UserNameVal   string
	SubRedditName string
	ContentVar    PostStruct
}

type SendDirectMessagetoUserStruct struct {
	UserNameVal string
	ContentVar  string
}
