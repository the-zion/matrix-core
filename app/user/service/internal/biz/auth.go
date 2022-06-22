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
	"time"
)

type AuthRepo interface {
	FindUserByPhone(ctx context.Context, phone string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUserWithPhone(ctx context.Context, phone string) (*User, error)
	CreateUserWithEmail(ctx context.Context, email, password string) (*User, error)
	CreateUserProfile(ctx context.Context, account, uuid string) error
	CreateUserProfileUpdate(ctx context.Context, account, uuid string) error
	CreateProfileUpdateRetry(ctx context.Context, account, uuid string) error
	SendPhoneCode(ctx context.Context, template, phone string) error
	SendEmailCode(ctx context.Context, template, phone string) error
	VerifyPhoneCode(ctx context.Context, phone, code string) error
	VerifyEmailCode(ctx context.Context, email, code string) error
	VerifyPassword(ctx context.Context, account, password, mode string) (*User, error)
	PasswordResetByPhone(ctx context.Context, phone, password string) error
	PasswordResetByEmail(ctx context.Context, email, password string) error
	GetCosSessionKey(ctx context.Context) (*Credentials, error)
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
		err = r.repo.CreateProfileUpdateRetry(ctx, email, user.Uuid)
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
			err = r.repo.CreateProfileUpdateRetry(ctx, phone, user.Uuid)
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

func (r *AuthUseCase) GetCosSessionKey(ctx context.Context) (*Credentials, error) {
	credentials, err := r.repo.GetCosSessionKey(ctx)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}
