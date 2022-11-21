package biz

type ImageReview struct {
	Code       string
	Message    string
	JobId      string
	State      string
	Object     string
	Url        string
	Label      string
	Result     int32
	Score      int32
	Category   string
	SubLabel   string
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
	Section      string
}

type MailBox struct {
	Time int32
}

type Notification struct {
	Timeline           map[string]int32
	Comment            int32
	SubComment         int32
	SystemNotification int32
}

//easyjson:json
type SystemNotification struct {
	Id               int32
	ContentId        int32
	CreatedAt        string
	NotificationType string
	Title            string
	Uid              string
	Uuid             string
	Label            string
	Result           int32
	Section          string
	Text             string
	Comment          string
}
