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

type validateLoginPassword struct {
	Account  string `validate:"required,e164|email"`
	Password string `validate:"required"`
	Code     string `validate:"required,len=6,number"`
	Mode     string `validate:"required,eq=phone|eq=email"`
}

type validateSendCode struct {
	Account  string `validate:"required,e164|email"`
	Mode     string `validate:"required,eq=phone|eq=email"`
	Template int64  `validate:"required,gte=1,lte=3"`
}

type validateGetUser struct {
	Id int64 `validate:"required"`
}

func loginVerify(validate *validator.Validate, log *log.Helper, account, password, code, mode string) bool {
	return verify(validate, log, &validateLogin{
		Account:  account,
		Password: password,
		Code:     code,
		Mode:     mode,
	})
}

func loginPasswordVerify(validate *validator.Validate, log *log.Helper, account, password, code, mode string) bool {
	return verify(validate, log, &validateLoginPassword{
		Account:  account,
		Password: password,
		Code:     code,
		Mode:     mode,
	})
}

func sendCodeVerify(validate *validator.Validate, log *log.Helper, template int64, account, mode string) bool {
	return verify(validate, log, &validateSendCode{
		Account:  account,
		Mode:     mode,
		Template: template,
	})
}

func getUserVerify(validate *validator.Validate, log *log.Helper, id int64) bool {
	return verify(validate, log, &validateGetUser{
		Id: id,
	})
}

func verify(validate *validator.Validate, log *log.Helper, s interface{}) bool {
	err := validate.Struct(s)
	if err != nil {
		log.Errorf("fail to verify login params: error(%v)", err)
		return false
	}
	return true
}
