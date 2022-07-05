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

type ArticleDraft struct {
	Id     int32
	Status int32
}
