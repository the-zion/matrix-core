package biz

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"time"
)

type AuthRepo interface {
	FindUserByPhone(ctx context.Context, phone string) (string, error)
	CreateUserWithPhone(ctx context.Context, phone string) (string, error)
	CreateUserWithEmail(ctx context.Context, email, password string) (string, error)
	CreateUserProfile(ctx context.Context, account, uuid string) error
	SendPhoneCode(ctx context.Context, template, phone string) error
	SendEmailCode(ctx context.Context, template, phone string) error
	//SendCode(ctx context.Context, template int64, account, mode string) error
	//UserRegister(ctx context.Context, account, mode string) (*User, error)
	VerifyPhoneCode(ctx context.Context, phone, code string) error
	VerifyEmailCode(ctx context.Context, email, code string) error
	SendCode(msgs ...*primitive.MessageExt)
}

type AuthUseCase struct {
	key      string
	repo     AuthRepo
	userRepo UserRepo
	tm       Transaction
	log      *log.Helper
}

func NewAuthUseCase(conf *conf.Auth, repo AuthRepo, userRepo UserRepo, tm Transaction, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		key:      conf.ApiKey,
		repo:     repo,
		userRepo: userRepo,
		tm:       tm,
		log:      log.NewHelper(log.With(logger, "module", "user/biz/authUseCase")),
	}
}

func (r *AuthUseCase) UserRegister(ctx context.Context, email, password, code string) error {
	err := r.repo.VerifyEmailCode(ctx, email, code)
	if err != nil {
		return v1.ErrorVerifyCodeFailed("register failed: %s", err.Error())
	}
	err = r.tm.ExecTx(ctx, func(ctx context.Context) error {
		uuid, err := r.repo.CreateUserWithEmail(ctx, email, password)
		if err != nil {
			return err
		}
		err = r.repo.CreateUserProfile(ctx, email, uuid)
		if err != nil {
			return err
		}
		return nil
	})
	if kerrors.IsConflict(err) {
		return v1.ErrorEmailConflict("register failed: %s", err.Error())
	}
	if err != nil {
		return v1.ErrorRegisterFailed("register failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) LoginByPassword(ctx context.Context, account, password, mode string) (string, error) {
	//user, err := r.userRepo.FindByAccount(ctx, account, mode)
	//if err != nil {
	//	return nil, v1.ErrorGetUserFailed("login failed: %s", err.Error())
	//}
	//
	//err = r.userRepo.VerifyPassword(ctx, user.Id, password)
	//if err != nil {
	//	return nil, v1.ErrorVerifyPasswordFailed("login failed: %s", err.Error())
	//}
	//
	//token, err := signToken(user.Id, r.key)
	//if err != nil {
	//	return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	//}
	//
	//return &Login{
	//	Id:    user.Id,
	//	Token: token,
	//}, nil
	return "", nil
}

func (r *AuthUseCase) LoginByCode(ctx context.Context, phone, code string) (string, error) {
	err := r.repo.VerifyPhoneCode(ctx, phone, code)
	if err != nil {
		return "", v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	}

	uuid, err := r.repo.FindUserByPhone(ctx, phone)
	if kerrors.IsNotFound(err) {
		err = r.tm.ExecTx(ctx, func(ctx context.Context) error {
			uuid, err = r.repo.CreateUserWithPhone(ctx, phone)
			if err != nil {
				return err
			}
			err = r.repo.CreateUserProfile(ctx, phone, uuid)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return "", v1.ErrorLoginFailed("login failed: %s", err.Error())
		}
	}
	if err != nil {
		return "", v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	token, err := signToken(uuid, r.key)
	if err != nil {
		return "", v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	return token, nil
}

//func (r *AuthUseCase) Register(ctx context.Context, account, mode string) (*User, error) {
//	user, err := r.repo.UserRegister(ctx, account, strings.ToUpper(mode[:1])+mode[1:])
//	if err != nil {
//		return nil, err
//	}
//	return user, nil
//}

func (r *AuthUseCase) LoginPasswordForget(ctx context.Context, account, password, code, mode string) (string, error) {
	//err := r.userRepo.VerifyCode(ctx, account, code, mode)
	//if err != nil {
	//	return nil, v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	//}
	//
	//user, err := r.userRepo.FindByAccount(ctx, account, mode)
	//if err != nil {
	//	return nil, v1.ErrorGetUserFailed("login failed: %s", err.Error())
	//}
	//
	//err = r.userRepo.PasswordModify(ctx, user.Id, password)
	//if err != nil {
	//	return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	//}
	//
	//token, err := signToken(user.Id, r.key)
	//if err != nil {
	//	return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	//}
	//
	//return &Login{
	//	Id:    user.Id,
	//	Token: token,
	//}, nil
	return "", nil
}

func (r *AuthUseCase) SendPhoneCode(ctx context.Context, template, phone string) error {
	err := r.repo.SendPhoneCode(ctx, template, phone)
	if err != nil {
		return v1.ErrorSendCodeFailed("send code failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) SendEmailCode(ctx context.Context, template, email string) error {
	err := r.repo.SendEmailCode(ctx, template, email)
	if err != nil {
		return v1.ErrorSendCodeFailed("send code failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) SendCode(msgs ...*primitive.MessageExt) {
	r.repo.SendCode(msgs...)
}

func signToken(uuid, key string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(604800 * time.Second)
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
