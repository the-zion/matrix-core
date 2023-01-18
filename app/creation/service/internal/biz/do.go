package biz

type TimeLine struct {
	Id         int32
	Uuid       string
	CreationId int32
	Mode       int32
}

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
	Auth      int32
}

type TalkStatistic struct {
	TalkId  int32
	Uuid    string
	Agree   int32
	View    int32
	Collect int32
	Comment int32
	Auth    int32
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

//easyjson:json
type ArticleReview struct {
	Uuid string
	Id   int32
	Mode string
}

//easyjson:json
type TalkReview struct {
	Uuid string
	Id   int32
	Mode string
}

type LeaderBoard struct {
	Id    int32
	Agree int32
	Uuid  string
	Mode  string
}

type Collections struct {
	CollectionsId int32
	Uuid          string
	Article       int32
	Column        int32
	Talk          int32
	Auth          int32
}

type CollectionsDraft struct {
	Id     int32
	Status int32
	Uuid   string
}

//easyjson:json
type CollectionsReview struct {
	Uuid string
	Id   int32
	Mode string
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

//easyjson:json
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

type Subscribe struct {
	ColumnId int32
	AuthorId string
	Uuid     string
	Status   int32
}

type News struct {
	Id     string
	Update string
	Title  string
	Author string
	Text   string
	Tags   string
	Cover  string
	Url    string
}

type ArticleSearch struct {
	Id     int32
	Total  int32
	Title  string
	Tags   string
	Uuid   string
	Text   string
	Cover  string
	Update string
}

type TalkSearch struct {
	Id     int32
	Total  int32
	Title  string
	Tags   string
	Uuid   string
	Text   string
	Cover  string
	Update string
}

type ColumnSearch struct {
	Id        int32
	Total     int32
	Name      string
	Tags      string
	Uuid      string
	Introduce string
	Cover     string
	Update    string
}

type NewsSearch struct {
	Id      string
	Update  string
	Title   string
	Author  string
	Tags    string
	Url     string
	Content string
	Cover   string
}

type CreationUser struct {
	Article     int32
	Column      int32
	Talk        int32
	Collections int32
	Collect     int32
	Subscribe   int32
}

//easyjson:json
type ImageReview struct {
	Id         int32
	CreationId int32
	Kind       string
	Uid        string
	CreateAt   string
	Uuid       string
	JobId      string
	Url        string
	Label      string
	Result     int32
	Category   string
	SubLabel   string
	Mode       string
	Score      int32
}

//easyjson:json
type TextReview struct {
	Id         int32
	CreationId int32
	CreateAt   string
	Title      string
	Kind       string
	JobId      string
	Label      string
	Result     int32
	Uuid       string
	Mode       string
	Section    string
}

//easyjson:json
type SendCollectionMap struct {
	Uuid string
	Id   int32
	Auth int32
	Mode string
}

//easyjson:json
type SendScoreMap struct {
	Uuid  string
	Score int32
	Mode  string
}

//easyjson:json
type SendArticleMap struct {
	Uuid string
	Id   int32
	Auth int32
	Mode string
}

//easyjson:json
type SendStatisticMap struct {
	Id            int32
	CollectionsId int32
	Uuid          string
	UserUuid      string
	Mode          string
}

//easyjson:json
type SendArticleStatisticMap struct {
	Uuid     string
	UserUuid string
	Mode     string
}

//easyjson:json
type SendColumnMap struct {
	Uuid string
	Id   int32
	Auth int32
	Mode string
}

//easyjson:json
type SendColumnIncludesMap struct {
	Id        int32
	ArticleId int32
	Uuid      string
	Mode      string
}

//easyjson:json
type SendColumnSubscribeMap struct {
	Id   int32
	Uuid string
	Mode string
}

//easyjson:json
type SendColumnStatisticMap struct {
	Uuid     string
	UserUuid string
	Mode     string
}

//easyjson:json
type SendTalkMap struct {
	Uuid string
	Id   int32
	Auth int32
	Mode string
}

//easyjson:json
type SendTalkStatisticMap struct {
	Uuid     string
	UserUuid string
	Mode     string
}
