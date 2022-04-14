package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-playground/validator/v10"
)

type validateLogin struct {
	Account  string `validate:"required,e164|email"`
	Password string `validate:"required_without=Code"`
	Code     string `validate:"required_without=Password,omitempty,len=6,number"`
	Mode     string `validate:"required,eq=phone|eq=email"`
}

func loginVerify(validate *validator.Validate, log *log.Helper, account, password, code, mode string) bool {
	v := &validateLogin{
		Account:  account,
		Password: password,
		Code:     code,
		Mode:     mode,
	}
	err := validate.Struct(v)
	if err != nil {
		log.Errorf("fail to verify login params: error(%v)", err)
		return false
	}
	return true
}
