package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo interface {
	UserRegister(ctx context.Context, email, password, code string) error
	LoginByPassword(ctx context.Context, account, password, mode string) (string, error)
	LoginByCode(ctx context.Context, phone, code string) (string, error)
	LoginPasswordReset(ctx context.Context, account, password, code, mode string) error
	SendPhoneCode(ctx context.Context, template, phone string) error
	SendEmailCode(ctx context.Context, template, email string) error
	GetCosSessionKey(ctx context.Context, uuid string) (*Credentials, error)
	GetAccount(ctx context.Context, uuid string) (*UserAccount, error)
	GetProfile(ctx context.Context, uuid string) (*UserProfile, error)
	GetProfileList(ctx context.Context, uuids []string) ([]*UserProfile, error)
	GetUserInfo(ctx context.Context, uuid string) (*UserProfile, error)
	GetProfileUpdate(ctx context.Context, uuid string) (*UserProfileUpdate, error)
	GetFollowList(ctx context.Context, page int32, uuid string) ([]*Follow, error)
	GetFollowListCount(ctx context.Context, uuid string) (int32, error)
	GetFollowedList(ctx context.Context, page int32, uuid string) ([]*Follow, error)
	GetFollowedListCount(ctx context.Context, uuid string) (int32, error)
	GetUserFollow(ctx context.Context, uuid, userUuid string) (bool, error)
	GetUserFollows(ctx context.Context, userId string, uuids []string) ([]*Follows, error)
	SetProfileUpdate(ctx context.Context, profile *UserProfileUpdate) error
	SetUserPhone(ctx context.Context, uuid, phone, code string) error
	SetUserPassword(ctx context.Context, uuid, password string) error
	SetUserEmail(ctx context.Context, uuid, email, code string) error
	SetUserFollow(ctx context.Context, uuid, userUuid string) error
	CancelUserFollow(ctx context.Context, uuid, userUuid string) error
	ChangeUserPassword(ctx context.Context, uuid, oldpassword, password string) error
	UnbindUserPhone(ctx context.Context, uuid, phone, code string) error
	UnbindUserEmail(ctx context.Context, uuid, email, code string) error
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/UserUseCase")),
	}
}

func (r *UserUseCase) UserRegister(ctx context.Context, email, password, code string) error {
	return r.repo.UserRegister(ctx, email, password, code)
}

func (r *UserUseCase) LoginByPassword(ctx context.Context, account, password, mode string) (string, error) {
	return r.repo.LoginByPassword(ctx, account, password, mode)
}

func (r *UserUseCase) LoginByCode(ctx context.Context, phone, code string) (string, error) {
	return r.repo.LoginByCode(ctx, phone, code)
}

func (r *UserUseCase) LoginPasswordReset(ctx context.Context, account, password, code, mode string) error {
	return r.repo.LoginPasswordReset(ctx, account, password, code, mode)
}

func (r *UserUseCase) SendPhoneCode(ctx context.Context, template, phone string) error {
	return r.repo.SendPhoneCode(ctx, template, phone)
}

func (r *UserUseCase) SendEmailCode(ctx context.Context, template, email string) error {
	return r.repo.SendEmailCode(ctx, template, email)
}

func (r *UserUseCase) GetCosSessionKey(ctx context.Context) (*Credentials, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetCosSessionKey(ctx, uuid)
}

func (r *UserUseCase) GetAccount(ctx context.Context) (*UserAccount, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetAccount(ctx, uuid)
}

func (r *UserUseCase) GetProfile(ctx context.Context) (*UserProfile, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetProfile(ctx, uuid)
}

func (r *UserUseCase) GetProfileList(ctx context.Context, uuids []string) ([]*UserProfile, error) {
	return r.repo.GetProfileList(ctx, uuids)
}

func (r *UserUseCase) GetUserInfo(ctx context.Context, uuid string) (*UserProfile, error) {
	return r.repo.GetUserInfo(ctx, uuid)
}

func (r *UserUseCase) GetProfileUpdate(ctx context.Context) (*UserProfileUpdate, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetProfileUpdate(ctx, uuid)
}

func (r *UserUseCase) GetUserFollow(ctx context.Context, uuid string) (bool, error) {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.GetUserFollow(ctx, uuid, userUuid)
}

func (r *UserUseCase) GetFollowList(ctx context.Context, page int32, uuid string) ([]*Follow, error) {
	return r.repo.GetFollowList(ctx, page, uuid)
}

func (r *UserUseCase) GetFollowListCount(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetFollowListCount(ctx, uuid)
}

func (r *UserUseCase) GetFollowedList(ctx context.Context, page int32, uuid string) ([]*Follow, error) {
	return r.repo.GetFollowedList(ctx, page, uuid)
}

func (r *UserUseCase) GetFollowedListCount(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetFollowedListCount(ctx, uuid)
}

func (r *UserUseCase) GetUserFollows(ctx context.Context, uuids []string) ([]*Follows, error) {
	userId := ctx.Value("uuid").(string)
	return r.repo.GetUserFollows(ctx, userId, uuids)
}

func (r *UserUseCase) SetUserProfile(ctx context.Context, profile *UserProfileUpdate) error {
	uuid := ctx.Value("uuid").(string)
	profile.Uuid = uuid
	return r.repo.SetProfileUpdate(ctx, profile)
}

func (r *UserUseCase) SetUserPhone(ctx context.Context, phone, code string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SetUserPhone(ctx, uuid, phone, code)
}

func (r *UserUseCase) SetUserEmail(ctx context.Context, email, code string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SetUserEmail(ctx, uuid, email, code)
}

func (r *UserUseCase) SetUserPassword(ctx context.Context, password string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SetUserPassword(ctx, uuid, password)
}

func (r *UserUseCase) SetUserFollow(ctx context.Context, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetUserFollow(ctx, uuid, userUuid)
}

func (r *UserUseCase) CancelUserFollow(ctx context.Context, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelUserFollow(ctx, uuid, userUuid)
}

func (r *UserUseCase) ChangeUserPassword(ctx context.Context, oldpassword, password string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.ChangeUserPassword(ctx, uuid, oldpassword, password)
}

func (r *UserUseCase) UnbindUserPhone(ctx context.Context, phone, code string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.UnbindUserPhone(ctx, uuid, phone, code)
}

func (r *UserUseCase) UnbindUserEmail(ctx context.Context, email, code string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.UnbindUserEmail(ctx, uuid, email, code)
}
