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
	CommentId    int32
	CreationId   int32
	CreationType int32
	Uuid         string
	Agree        int32
	Comment      int32
}

type SubComment struct {
	CommentId int32
	RootId    int32
	ParentId  int32
	Uuid      string
	Reply     string
	UserName  string
	ReplyName string
	Agree     int32
}

type CommentStatistic struct {
	CommentId int32
	Agree     int32
	Comment   int32
}
