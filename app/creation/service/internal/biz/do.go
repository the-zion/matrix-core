package biz

type Article struct {
	ArticleId int32
	Uuid      string
	Status    int32
	Auth      int32
}

type Talk struct {
	TalkId int32
	Uuid   string
	Status int32
	Auth   int32
}

type ArticleStatistic struct {
	ArticleId int32
	Uuid      string
	Agree     int32
	View      int32
	Collect   int32
	Comment   int32
}

type TalkStatistic struct {
	TalkId  int32
	Uuid    string
	Agree   int32
	View    int32
	Collect int32
	Comment int32
}

type ArticleStatisticJudge struct {
	Agree   bool
	Collect bool
}

type TalkStatisticJudge struct {
	Agree   bool
	Collect bool
}

type ArticleDraft struct {
	Id     int32
	Status int32
	Uuid   string
}

type TalkDraft struct {
	Id     int32
	Status int32
	Uuid   string
}

type ArticleReview struct {
	Uuid string
	Id   int32
	Mode string
}

type TalkReview struct {
	Uuid string
	Id   int32
	Mode string
}

type LeaderBoard struct {
	Id   int32
	Uuid string
	Mode string
}

type Collections struct {
	Id        int32
	Uuid      string
	Name      string
	Introduce string
	Auth      int32
}

type Column struct {
	ColumnId int32
	Uuid     string
	Auth     int32
}

type ColumnDraft struct {
	Id     int32
	Status int32
	Uuid   string
}

type ColumnReview struct {
	Uuid string
	Id   int32
	Mode string
}

type ColumnStatistic struct {
	ColumnId int32
	Uuid     string
	Agree    int32
	View     int32
	Collect  int32
	Auth     int32
}

type ColumnStatisticJudge struct {
	Agree   bool
	Collect bool
}
