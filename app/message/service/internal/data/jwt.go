package data

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"github.com/the-zion/matrix-core/app/message/service/internal/conf"
)

var (
	LOG                       = log.NewHelper(log.With(log.GetLogger(), "source", "accesslog"))
	ErrMissingJwtToken        = errors.Unauthorized("TOKEN_MISSING", "JWT token is missing")
	ErrTokenInvalid           = errors.Unauthorized("TOKEN_INVALID", "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("TOKEN_EXPIRED", "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized("TOKEN_PARSE_FAIL", "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("UN_SUPPORT_SIGNING_METHOD", "Wrong signing method")
)

func claimsFunc() jwt.Claims {
	return &jwt.MapClaims{}
}

func (d *Data) JwtCheck(jwtToken string) (string, error) {
	var (
		tokenInfo *jwt.Token
		err       error
	)
	tokenInfo, err = jwt.ParseWithClaims(jwtToken, claimsFunc(), func(token *jwt.Token) (interface{}, error) {
		return []byte(d.jwt.key), nil
	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if !ok {
			return "", errors.Unauthorized("UNAUTHORIZED", err.Error())
		}
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "", ErrTokenInvalid
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return "", ErrTokenExpired
		}
		return "", ErrTokenParseFail
	}
	if !tokenInfo.Valid {
		return "", ErrTokenInvalid
	}
	if tokenInfo.Method != jwt.SigningMethodHS256 {
		return "", ErrUnSupportSigningMethod
	}
	token := *((tokenInfo.Claims.(jwt.Claims)).(*jwt.MapClaims))
	return token["uuid"].(string), nil
}

func NewJwt(d *Data) biz.Jwt {
	return d
}

func NewJwtClient(conf *conf.Data) Jwt {
	return Jwt{
		key: conf.Jwt.Key,
	}
}
