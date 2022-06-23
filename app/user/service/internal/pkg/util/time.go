package util

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func StringToTime(datetime string) (time.Time, error) {
	formatTime, err := time.Parse("2006-01-02 15:04:05", datetime)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, fmt.Sprintf("fail to parse time: time(%v)", datetime))
	}
	return formatTime, nil
}
