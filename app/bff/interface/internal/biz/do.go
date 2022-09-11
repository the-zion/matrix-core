package biz

type UserAccount struct {
	Phone    string
	Email    string
	Qq       string
	Wechat   string
	Weibo    string
	Github   string
	Password string
}

type UserProfile struct {
	Uuid        string
	Username    string
	Avatar      string
	School      string
	Company     string
	Job         string
	Homepage    string
	Introduce   string
	Created     string
	Score       int32
	Agree       int32
	Collect     int32
	View        int32
	Follow      int32
	Followed    int32
	Article     int32
	Column      int32
	Subscribe   int32
	Talk        int32
	Collections int32
}

type UserSearch struct {
	Uuid        string
	Username    string
	Introduce   string
	Agree       int32
	View        int32
	Collect     int32
	FollowNum   int32
	FollowedNum int32
}

type UserProfileUpdate struct {
	UserProfile
	Status int32
}

type Credentials struct {
	TmpSecretID  string
	TmpSecretKey string
	SessionToken string
	StartTime    int64
	ExpiredTime  int64
}

type UserImageReview struct {
	Id       int32
	Uuid     string
	CreateAt string
	JobId    string
	Url      string
	Label    string
	Result   int32
	Score    int32
	Category string
	SubLabel string
}

type CreationImageReview struct {
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

type CreationContentReview struct {
	Id         int32
	CreationId int32
	Title      string
	Kind       string
	CreateAt   string
	Uuid       string
	JobId      string
	Label      string
	Result     int32
	Section    string
	Mode       string
}

type Follow struct {
	Follow      string
	Followed    string
	Username    string
	Introduce   string
	Agree       int32
	View        int32
	Collect     int32
	FollowNum   int32
	FollowedNum int32
	Status      int32
}

type Follows struct {
	Uuid   string
	Follow int32
}

type PornInfo struct {
	HitFlag int32
	Count   int32
}
type AdsInfo struct {
	HitFlag int32
	Count   int32
}

type IllegalInfo struct {
	HitFlag int32
	Count   int32
}

type AbuseInfo struct {
	HitFlag int32
	Count   int32
}

type Section struct {
	Label       string
	Result      int32
	PornInfo    *SectionPornInfo
	AdsInfo     *SectionAdsInfo
	IllegalInfo *SectionIllegalInfo
	AbuseInfo   *SectionAbuseInfo
}

type SectionPornInfo struct {
	HitFlag  int32
	Score    int32
	Keywords string
}

type SectionAdsInfo struct {
	HitFlag  int32
	Score    int32
	Keywords string
}

type SectionIllegalInfo struct {
	HitFlag  int32
	Score    int32
	Keywords string
}

type SectionAbuseInfo struct {
	HitFlag  int32
	Score    int32
	Keywords string
}

type TextReview struct {
	Code         string
	Message      string
	JobId        string
	DataId       string
	State        string
	CreationTime string
	Object       string
	Label        string
	Result       int32
	PornInfo     *PornInfo
	AdsInfo      *AdsInfo
	IllegalInfo  *IllegalInfo
	AbuseInfo    *AbuseInfo
	BucketId     string
	Region       string
	CosHeaders   map[string]string
	Section      []*Section
}

type LeaderBoard struct {
	Id   int32
	Uuid string
	Mode string
}

type Article struct {
	Id      int32
	Agree   int32
	View    int32
	Collect int32
	Comment int32
	Total   int32
	Title   string
	Tags    string
	Uuid    string
	Text    string
	Cover   string
	Update  string
}

type ArticleStatistic struct {
	Id      int32
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

type ArticleDraft struct {
	Id     int32
	Status int32
}

type Talk struct {
	Id      int32
	Agree   int32
	View    int32
	Collect int32
	Comment int32
	Total   int32
	Title   string
	Tags    string
	Uuid    string
	Text    string
	Cover   string
	Update  string
}

type TalkDraft struct {
	Id     int32
	Status int32
}

type TalkStatistic struct {
	Id      int32
	Uuid    string
	Agree   int32
	View    int32
	Collect int32
	Comment int32
}

type TalkStatisticJudge struct {
	Agree   bool
	Collect bool
}

type Collections struct {
	Id      int32
	Uuid    string
	Auth    int32
	Article int32
	Column  int32
	Talk    int32
}

type CollectionsDraft struct {
	Id     int32
	Status int32
}

type Column struct {
	Id        int32
	Agree     int32
	View      int32
	Collect   int32
	Total     int32
	Name      string
	Tags      string
	Uuid      string
	Introduce string
	Cover     string
	Update    string
}

type ColumnDraft struct {
	Id     int32
	Status int32
}

type ColumnStatistic struct {
	Id      int32
	Uuid    string
	Agree   int32
	View    int32
	Collect int32
}

type ColumnStatisticJudge struct {
	Agree   bool
	Collect bool
}

type CreationUser struct {
	Article     int32
	Column      int32
	Talk        int32
	Collect     int32
	Subscribe   int32
	Collections int32
}

type Achievement struct {
	Uuid     string
	Score    int32
	Agree    int32
	Collect  int32
	View     int32
	Follow   int32
	Followed int32
}

type Active struct {
	Agree int32
}

type Medal struct {
	Creation1 int32
	Creation2 int32
	Creation3 int32
	Creation4 int32
	Creation5 int32
	Creation6 int32
	Creation7 int32
	Agree1    int32
	Agree2    int32
	Agree3    int32
	Agree4    int32
	Agree5    int32
	Agree6    int32
	View1     int32
	View2     int32
	View3     int32
	Comment1  int32
	Comment2  int32
	Comment3  int32
	Collect1  int32
	Collect2  int32
	Collect3  int32
}

type MedalProgress struct {
	Article     int32
	Talk        int32
	Agree       int32
	ActiveAgree int32
	View        int32
	Comment     int32
	Collect     int32
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

type CommentDraft struct {
	Id     int32
	Status int32
}

type Comment struct {
	Id             int32
	Uuid           string
	UserName       string
	CreationAuthor string
	CreationId     int32
	Agree          int32
	Comment        int32
}

type SubComment struct {
	Id             int32
	CreationId     int32
	RootId         int32
	ParentId       int32
	CreationAuthor string
	RootUser       string
	Uuid           string
	Reply          string
	UserName       string
	ReplyName      string
	RootName       string
	Agree          int32
}

type CommentUser struct {
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

type CommentStatistic struct {
	Id      int32
	Agree   int32
	Comment int32
}

type CommentContentReview struct {
	Id        int32
	CommentId int32
	Comment   string
	Kind      string
	CreateAt  string
	Uuid      string
	JobId     string
	Label     string
	Result    int32
	Section   string
	Mode      string
}
