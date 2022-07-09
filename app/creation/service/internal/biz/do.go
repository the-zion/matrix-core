package biz

type Article struct {
	ArticleId int32
	Uuid      string
	Status    int32
}

type ArticleStatistic struct {
	ArticleId int32
	Uuid      string
	Agree     int32
	View      int32
	Collect   int32
	Comment   int32
}

type ArticleStatisticJudge struct {
	Agree   bool
	Collect bool
}

type ArticleDraft struct {
	Id     int32
	Status int32
	Uuid   string
}

type LeaderBoard struct {
	Id   int32
	Uuid string
	Mode string
}

type Collections struct {
	Name      string
	Introduce string
}
