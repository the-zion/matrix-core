package biz

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
