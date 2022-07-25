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
	Uuid      string
	Username  string
	Avatar    string
	School    string
	Company   string
	Job       string
	Homepage  string
	Introduce string
	Created   string
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

type AvatarReview struct {
	Code       string
	Message    string
	JobId      string
	State      string
	Object     string
	Label      string
	Result     int32
	Category   string
	BucketId   string
	Region     string
	CosHeaders map[string]string
	EventName  string
}

type Follow struct {
	Follow   string
	Followed string
	Status   int32
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
	Id   int32
	Uuid string
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
	Id   int32
	Uuid string
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
	Id        int32
	Uuid      string
	Name      string
	Introduce string
	Auth      int32
}

type Column struct {
	Id   int32
	Uuid string
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

type Achievement struct {
	Uuid     string
	Agree    int32
	Collect  int32
	View     int32
	Follow   int32
	Followed int32
}

type Subscribe struct {
	ColumnId int32
	AuthorId string
	Uuid     string
	Status   int32
}
