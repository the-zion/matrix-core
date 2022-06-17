package biz

type Profile struct {
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
	Status int64
}

type Credentials struct {
	TmpSecretID  string
	TmpSecretKey string
	SessionToken string
	StartTime    int64
	ExpiredTime  int64
}
