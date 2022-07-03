package biz

type Article struct {
	ArticleId int32
	Uuid      string
	Status    int32
}

type ArticleDraft struct {
	Id     int32
	Status int32
	Uuid   string
}
