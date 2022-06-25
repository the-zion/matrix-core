package biz

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
