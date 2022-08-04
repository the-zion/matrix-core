package biz

type CommentDraft struct {
	Id   int32
	Uuid string
}

type CommentReview struct {
	Uuid string
	Id   int32
}
