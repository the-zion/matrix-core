package biz

type UserProfile struct {
	Uuid      string
	Username  string
	Avatar    string
	School    string
	Company   string
	Homepage  string
	Introduce string
}

type Credentials struct {
	TmpSecretID  string
	TmpSecretKey string
	SessionToken string
	StartTime    int64
	ExpiredTime  int64
}
