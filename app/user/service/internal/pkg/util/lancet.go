package util

import "github.com/duke-git/lancet/random"

func UUIdV4() (string, error) {
	uuid, err := random.UUIdV4()
	if err != nil {
		return "", err
	}
	return uuid, nil
}
