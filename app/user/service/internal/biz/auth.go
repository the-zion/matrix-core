package biz

import (
	"context"
	"fmt"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"github.com/the-zion/matrix-core/app/user/service/internal/pkg/util"
	"time"
)

type AuthRepo interface {
	GetGithubAccessToken(ctx context.Context, code string) (string, error)
	GetGithubUserInfo(ctx context.Context, token string) (map[string]interface{}, error)
	FindUserByPhone(ctx context.Context, phone string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByGithub(ctx context.Context, github int32) (*User, error)
	CreateUserWithPhone(ctx context.Context, phone string) (*User, error)
	CreateUserWithEmail(ctx context.Context, email, password string) (*User, error)
	CreateUserWithGithub(ctx context.Context, github int32) (*User, error)
	CreateUserProfile(ctx context.Context, account, uuid string) error
	CreateUserProfileWithGithub(ctx context.Context, account, github, uuid string) error
	CreateUserProfileUpdate(ctx context.Context, account, uuid string) error
	CreateUserSearch(ctx context.Context, account, uuid string) error
	SetUserPhone(ctx context.Context, uuid, phone string) error
	SetUserEmail(ctx context.Context, uuid, email string) error
	SetUserPassword(ctx context.Context, uuid, password string) error
	SendPhoneCode(ctx context.Context, template, phone string) error
	SendEmailCode(ctx context.Context, template, phone string) error
	VerifyPhoneCode(ctx context.Context, phone, code string) error
	VerifyEmailCode(ctx context.Context, email, code string) error
	VerifyPassword(ctx context.Context, account, password, mode string) (*User, error)
	PasswordResetByPhone(ctx context.Context, phone, password string) error
	PasswordResetByEmail(ctx context.Context, email, password string) error
	GetCosSessionKey(ctx context.Context, uuid string) (*Credentials, error)
	UnbindUserPhone(ctx context.Context, uuid string) error
	UnbindUserEmail(ctx context.Context, uuid string) error
}

type AuthUseCase struct {
	key      string
	repo     AuthRepo
	userRepo UserRepo
	tm       Transaction
	re       Recovery
	log      *log.Helper
}

func NewAuthUseCase(conf *conf.Auth, repo AuthRepo, re Recovery, userRepo UserRepo, tm Transaction, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		key:      conf.ApiKey,
		repo:     repo,
		userRepo: userRepo,
		tm:       tm,
		re:       re,
		log:      log.NewHelper(log.With(logger, "module", "user/biz/authUseCase")),
	}
}

func (r *AuthUseCase) UserRegister(ctx context.Context, email, password, code string) error {
	err := r.repo.VerifyEmailCode(ctx, email, code)
	if err != nil {
		return v1.ErrorVerifyCodeFailed("register failed: %s", err.Error())
	}
	err = r.tm.ExecTx(ctx, func(ctx context.Context) error {
		user, err := r.repo.CreateUserWithEmail(ctx, email, password)
		if err != nil {
			return err
		}
		err = r.repo.CreateUserProfile(ctx, email, user.Uuid)
		if err != nil {
			return err
		}
		err = r.repo.CreateUserProfileUpdate(ctx, email, user.Uuid)
		if err != nil {
			return err
		}
		err = r.repo.CreateUserSearch(ctx, email, user.Uuid)
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
	user, err := r.repo.VerifyPassword(ctx, account, password, mode)
	if err != nil {
		return "", v1.ErrorVerifyPasswordFailed("login failed: %s", err.Error())
	}

	token, err := signToken(user.Uuid, r.key)
	if err != nil {
		return "", v1.ErrorLoginFailed("login failed: %s", err.Error())
	}
	return token, nil
}

func (r *AuthUseCase) LoginByCode(ctx context.Context, phone, code string) (string, error) {
	err := r.repo.VerifyPhoneCode(ctx, phone, code)
	if err != nil {
		return "", v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	}

	user, err := r.repo.FindUserByPhone(ctx, phone)
	if kerrors.IsNotFound(err) {
		err = r.tm.ExecTx(ctx, func(ctx context.Context) error {
			user, err = r.repo.CreateUserWithPhone(ctx, phone)
			if err != nil {
				return err
			}
			err = r.repo.CreateUserProfile(ctx, phone, user.Uuid)
			if err != nil {
				return err
			}
			err = r.repo.CreateUserProfileUpdate(ctx, phone, user.Uuid)
			if err != nil {
				return err
			}
			err = r.repo.CreateUserSearch(ctx, phone, user.Uuid)
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

	token, err := signToken(user.Uuid, r.key)
	if err != nil {
		return "", v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	return token, nil
}

func (r *AuthUseCase) LoginByGithub(ctx context.Context, code string) (string, error) {
	accessToken, err := r.repo.GetGithubAccessToken(ctx, code)
	if err != nil {
		return "", v1.ErrorLoginFailed("fail to get github access token: %s", err.Error())
	}

	userInfo, err := r.repo.GetGithubUserInfo(ctx, accessToken)
	if err != nil {
		return "", v1.ErrorLoginFailed("fail to get github user info: %s", err.Error())
	}

	githubId := int32(userInfo["id"].(float64))
	github := userInfo["html_url"].(string)
	name := fmt.Sprintf("%s-%v", userInfo["name"].(string), githubId)
	user, err := r.repo.FindUserByGithub(ctx, githubId)
	if kerrors.IsNotFound(err) {
		err = r.tm.ExecTx(ctx, func(ctx context.Context) error {
			user, err = r.repo.CreateUserWithGithub(ctx, githubId)
			if err != nil {
				return err
			}
			err = r.repo.CreateUserProfileWithGithub(ctx, name, github, user.Uuid)
			if err != nil {
				return err
			}
			err = r.repo.CreateUserProfileUpdate(ctx, name, user.Uuid)
			if err != nil {
				return err
			}
			err = r.repo.CreateUserSearch(ctx, name, user.Uuid)
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

	token, err := signToken(user.Uuid, r.key)
	if err != nil {
		return "", v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	return token, nil
}

func (r *AuthUseCase) LoginPasswordReset(ctx context.Context, account, password, code, mode string) error {
	if mode == "phone" {
		return r.passwordResetByPhone(ctx, account, password, code)
	} else {
		return r.passwordResetByEmail(ctx, account, password, code)
	}
}

func (r *AuthUseCase) passwordResetByPhone(ctx context.Context, phone, password, code string) error {
	err := r.repo.VerifyPhoneCode(ctx, phone, code)
	if err != nil {
		return v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	}

	err = r.repo.PasswordResetByPhone(ctx, phone, password)
	if err != nil {
		return v1.ErrorResetPasswordFailed("reset password failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) passwordResetByEmail(ctx context.Context, email, password, code string) error {
	err := r.repo.VerifyEmailCode(ctx, email, code)
	if err != nil {
		return v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	}

	err = r.repo.PasswordResetByEmail(ctx, email, password)
	if err != nil {
		return v1.ErrorResetPasswordFailed("reset password failed: %s", err.Error())
	}
	return nil
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

func (r *AuthUseCase) GetCosSessionKey(ctx context.Context, uuid string) (*Credentials, error) {
	credentials, err := r.repo.GetCosSessionKey(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func (r *AuthUseCase) SetUserPhone(ctx context.Context, uuid, phone, code string) error {
	err := r.repo.VerifyPhoneCode(ctx, phone, code)
	if err != nil {
		return v1.ErrorVerifyCodeFailed("set user phone failed: %s", err.Error())
	}

	err = r.repo.SetUserPhone(ctx, uuid, phone)
	if kerrors.IsConflict(err) {
		return v1.ErrorPhoneConflict("set user phone failed: %s", err.Error())
	}
	if err != nil {
		return v1.ErrorSetPhoneFailed("set user phone failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) SetUserEmail(ctx context.Context, uuid, email, code string) error {
	err := r.repo.VerifyEmailCode(ctx, email, code)
	if err != nil {
		return v1.ErrorVerifyCodeFailed("set user email failed: %s", err.Error())
	}

	err = r.repo.SetUserEmail(ctx, uuid, email)
	if kerrors.IsConflict(err) {
		return v1.ErrorEmailConflict("set user email failed: %s", err.Error())
	}
	if err != nil {
		return v1.ErrorSetEmailFailed("set user email failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) SetUserPassword(ctx context.Context, uuid, password string) error {
	err := r.repo.SetUserPassword(ctx, uuid, password)
	if err != nil {
		return v1.ErrorSetPasswordFailed("set user password failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) ChangeUserPassword(ctx context.Context, uuid, oldpassword, password string) error {
	account, err := r.userRepo.GetAccount(ctx, uuid)
	if err != nil {
		return v1.ErrorGetAccountFailed("get user account failed: %s", err.Error())
	}

	pass := util.CheckPasswordHash(oldpassword, account.Password)
	if !pass {
		return v1.ErrorVerifyPasswordFailed("fail to verify password: password(%s)", oldpassword)
	}

	err = r.repo.SetUserPassword(ctx, uuid, password)
	if err != nil {
		return v1.ErrorSetPasswordFailed("set user password failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) UnbindUserPhone(ctx context.Context, uuid, phone, code string) error {
	err := r.repo.VerifyPhoneCode(ctx, phone, code)
	if err != nil {
		return v1.ErrorVerifyCodeFailed("unbind user phone failed: %s", err.Error())
	}

	account, err := r.userRepo.GetAccount(ctx, uuid)
	if err != nil {
		return v1.ErrorUnbindEmailFailed("fail to get account: %s", err.Error())
	}

	if account.Email == "" && account.Qq == "" && account.Weibo == "" && account.Wechat == "" && account.Github == 0 {
		return v1.ErrorUniqueAccount("unbind user email failed: unique account")
	}

	err = r.repo.UnbindUserPhone(ctx, uuid)
	if err != nil {
		return v1.ErrorUnbindPhoneFailed("unbind user phone failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) UnbindUserEmail(ctx context.Context, uuid, email, code string) error {
	err := r.repo.VerifyEmailCode(ctx, email, code)
	if err != nil {
		return v1.ErrorVerifyCodeFailed("unbind user email failed: %s", err.Error())
	}

	account, err := r.userRepo.GetAccount(ctx, uuid)
	if err != nil {
		return v1.ErrorUnbindEmailFailed("fail to get account: %s", err.Error())
	}

	if account.Phone == "" && account.Qq == "" && account.Weibo == "" && account.Wechat == "" && account.Github == 0 {
		return v1.ErrorUniqueAccount("unbind user email failed: unique account")
	}

	err = r.repo.UnbindUserEmail(ctx, uuid)
	if err != nil {
		return v1.ErrorUnbindEmailFailed("unbind user email failed: %s", err.Error())
	}
	return nil
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
