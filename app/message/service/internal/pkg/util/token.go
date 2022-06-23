package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"time"
)

func SignToken(uuid, key string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(172800 * time.Second)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"uuid": uuid,
		"exp":  expireTime.Unix(),
		"iss":  "matrix",
	})
	signedString, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to sign token: uuid(%v)", uuid))
	}
	return signedString, nil
}
