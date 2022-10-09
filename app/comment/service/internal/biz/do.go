package biz

type CommentDraft struct {
	Id     int32
	Uuid   string
	Status int32
}

type CommentReview struct {
	Uuid string
	Id   int32
	Mode string
}

type Comment struct {
	UpdatedAt      int32
	CommentId      int32
	CreationId     int32
	CreationType   int32
	CreationAuthor string
	Uuid           string
	Agree          int32
	Comment        int32
}

type CommentUser struct {
	Uuid              string
	Comment           int32
	ArticleReply      int32
	ArticleReplySub   int32
	TalkReply         int32
	TalkReplySub      int32
	ArticleReplied    int32
	ArticleRepliedSub int32
	TalkReplied       int32
	TalkRepliedSub    int32
}

type SubComment struct {
	CommentId      int32
	CreationId     int32
	RootId         int32
	ParentId       int32
	CreationType   int32
	CreationAuthor string
	RootUser       string
	Uuid           string
	Reply          string
	UserName       string
	ReplyName      string
	RootName       string
	Agree          int32
}

type CommentStatistic struct {
	CommentId int32
	Agree     int32
	Comment   int32
}

type TextReview struct {
	Id        int32
	CommentId int32
	CreateAt  string
	Comment   string
	Kind      string
	JobId     string
	Label     string
	Result    int32
	Uuid      string
	Mode      string
	Section   string
}
