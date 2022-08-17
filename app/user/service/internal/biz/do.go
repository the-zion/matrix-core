package biz

import "time"

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
