package biz

import (
	"time"
)

type User struct {
	Uuid     string
	Phone    string
	Email    string
	Qq       string
	Wechat   string
	Weibo    string
	Github   string
	Password string
}

type Profile struct {
	Created   string
	Updated   string
	Uuid      string
	Username  string
	Avatar    string
	School    string
	Company   string
	Job       string
	Homepage  string
	Introduce string
}

type UserSearch struct {
	Uuid      string
	Username  string
	Introduce string
}

//easyjson:json
type ProfileUpdate struct {
	Profile
	Status int32
}

type Credentials struct {
	TmpSecretID  string
	TmpSecretKey string
	SessionToken string
	StartTime    int64
	ExpiredTime  int64
}

type Follow struct {
	Id       int32
	Update   time.Time
	Follow   string
	Followed string
	Status   int32
}

type Follows struct {
	Uuid   string
	Follow int32
}

//easyjson:json
type ImageReview struct {
	Id       int32
	CreateAt string
	Uuid     string
	JobId    string
	Url      string
	Label    string
	Result   int32
	Category string
	SubLabel string
	Mode     string
	Score    int32
}

//easyjson:json
type UserSearchMap struct {
	Username  string
	Introduce string
}

//easyjson:json
type ProfileUpdateMap struct {
	Username  string
	Introduce string
}

//easyjson:json
type SendUserStatisticMap struct {
	Follow   string
	Followed string
	Mode     string
}

//easyjson:json
type SetFollowMap struct {
	Uuid   string
	UserId string
	Mode   string
}
