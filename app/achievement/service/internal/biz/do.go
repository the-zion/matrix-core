package biz

type Achievement struct {
	Uuid     string
	Agree    int32
	Collect  int32
	View     int32
	Follow   int32
	Followed int32
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

type Active struct {
	Agree int32
}
