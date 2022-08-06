package biz

type CommentDraft struct {
	Id     int32
	Uuid   string
	Status int32
}

type CommentReview struct {
	Uuid string
	Id   int32
}

type Comment struct {
	CommentId    int32
	CreationId   int32
	CreationType int32
	Uuid         string
	Agree        int32
	Comment      int32
}

type CommentStatistic struct {
	CommentId int32
	Agree     int32
	Comment   int32
}
