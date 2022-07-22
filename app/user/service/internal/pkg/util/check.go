package util

import "github.com/duke-git/lancet/validator"

func Test() {
	validator.IsStrongPassword("qdq", 12)
}
