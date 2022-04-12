package util

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomNumber() string {
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	return code
}
