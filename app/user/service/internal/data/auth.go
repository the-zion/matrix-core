package data

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	utilService "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"github.com/the-zion/matrix-core/app/user/service/internal/pkg/util"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var _ biz.AuthRepo = (*authRepo)(nil)

type authRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/auth")),
	}
}

func (r *authRepo) FindUserByPhone(ctx context.Context, phone string) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Select("uuid", "password", "phone", "email", "wechat", "qq", "github", "gitee").Where("phone = ?", phone).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("phone not found from db", fmt.Sprintf("phone(%s)", phone))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to find user by phone: phone(%s)", phone))
	}
	return &biz.User{
		Uuid:     user.Uuid,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Wechat:   user.Wechat,
		Qq:       user.Qq,
		Github:   user.Github,
		Gitee:    user.Gitee,
	}, nil
}

func (r *authRepo) FindUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Select("uuid", "password", "phone", "email", "wechat", "qq", "github", "gitee").Where("email = ?", email).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("email not found from db", fmt.Sprintf("email(%s)", email))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to find user by email: email(%s)", email))
	}
	return &biz.User{
		Uuid:     user.Uuid,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Wechat:   user.Wechat,
		Qq:       user.Qq,
		Github:   user.Github,
		Gitee:    user.Gitee,
	}, nil
}

func (r *authRepo) FindUserByWechat(ctx context.Context, wechat string) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Select("uuid", "password", "phone", "email", "wechat", "qq", "github", "gitee").Where("wechat = ?", wechat).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("wechat not found from db", fmt.Sprintf("wechat(%s)", wechat))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to find user by wechat: wechat(%s)", wechat))
	}
	return &biz.User{
		Uuid:     user.Uuid,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Wechat:   user.Wechat,
		Qq:       user.Qq,
		Github:   user.Github,
		Gitee:    user.Gitee,
	}, nil
}

func (r *authRepo) FindUserByQQ(ctx context.Context, qq string) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Select("uuid", "password", "phone", "email", "wechat", "qq", "github", "gitee").Where("qq = ?", qq).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("qq not found from db", fmt.Sprintf("qq(%s)", qq))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to find user by qq: qq(%s)", qq))
	}
	return &biz.User{
		Uuid:     user.Uuid,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Wechat:   user.Wechat,
		Qq:       user.Qq,
		Github:   user.Github,
		Gitee:    user.Gitee,
	}, nil
}

func (r *authRepo) FindUserByGithub(ctx context.Context, github int32) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Select("uuid", "password", "phone", "email", "wechat", "qq", "github", "gitee").Where("github = ?", github).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("github not found from db", fmt.Sprintf("github(%v)", github))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to find user by github: github(%v)", github))
	}
	return &biz.User{
		Uuid:     user.Uuid,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Wechat:   user.Wechat,
		Qq:       user.Qq,
		Github:   user.Github,
		Gitee:    user.Gitee,
	}, nil
}

func (r *authRepo) FindUserByGitee(ctx context.Context, gitee int32) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Select("uuid", "password", "phone", "email", "wechat", "qq", "github", "gitee").Where("gitee = ?", gitee).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("gitee not found from db", fmt.Sprintf("gitee(%v)", gitee))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to find user by gitee: gitee(%v)", gitee))
	}
	return &biz.User{
		Uuid:     user.Uuid,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Wechat:   user.Wechat,
		Qq:       user.Qq,
		Github:   user.Github,
		Gitee:    user.Gitee,
	}, nil
}

func (r *authRepo) CreateUserWithPhone(ctx context.Context, phone string) (*biz.User, error) {
	uuid := xid.New().String()
	user := &User{
		Uuid:  uuid,
		Phone: phone,
	}
	err := r.data.DB(ctx).Select("Phone", "Uuid").Create(user).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to create a user: phone(%s)", phone))
	}

	return &biz.User{
		Uuid:  uuid,
		Phone: phone,
	}, nil
}

func (r *authRepo) CreateUserWithEmail(ctx context.Context, email, password string) (*biz.User, error) {
	uuid := xid.New().String()

	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to hash password: password(%s)", password))
	}

	user := &User{
		Uuid:     uuid,
		Email:    email,
		Password: hashPassword,
	}
	err = r.data.DB(ctx).Select("Email", "Uuid", "Password").Create(user).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil, kerrors.Conflict("email conflict", fmt.Sprintf("email(%s)", email))
		} else {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to create a user: email(%s)", email))
		}
	}

	return &biz.User{
		Uuid:     uuid,
		Email:    email,
		Password: hashPassword,
	}, nil
}

func (r *authRepo) CreateUserWithWechat(ctx context.Context, wechat string) (*biz.User, error) {
	uuid := xid.New().String()
	user := &User{
		Uuid:   uuid,
		Wechat: wechat,
	}
	err := r.data.DB(ctx).Select("Wechat", "Uuid").Create(user).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil, kerrors.Conflict("github conflict", fmt.Sprintf("wechat(%v)", wechat))
		} else {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to create a user: wechat(%v)", wechat))
		}
	}

	return &biz.User{
		Uuid:   uuid,
		Wechat: wechat,
	}, nil
}

func (r *authRepo) CreateUserWithQQ(ctx context.Context, qq string) (*biz.User, error) {
	uuid := xid.New().String()

	user := &User{
		Uuid: uuid,
		Qq:   qq,
	}
	err := r.data.DB(ctx).Select("Qq", "Uuid").Create(user).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil, kerrors.Conflict("github conflict", fmt.Sprintf("qq(%s)", qq))
		} else {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to create a user: qq(%s)", qq))
		}
	}

	return &biz.User{
		Uuid: uuid,
		Qq:   qq,
	}, nil
}

func (r *authRepo) CreateUserWithGithub(ctx context.Context, github int32) (*biz.User, error) {
	uuid := xid.New().String()

	user := &User{
		Uuid:   uuid,
		Github: github,
	}
	err := r.data.DB(ctx).Select("Github", "Uuid").Create(user).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil, kerrors.Conflict("github conflict", fmt.Sprintf("github(%v)", github))
		} else {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to create a user: github(%v)", github))
		}
	}

	return &biz.User{
		Uuid:   uuid,
		Github: github,
	}, nil
}

func (r *authRepo) CreateUserWithGitee(ctx context.Context, gitee int32) (*biz.User, error) {
	uuid := xid.New().String()

	user := &User{
		Uuid:  uuid,
		Gitee: gitee,
	}
	err := r.data.DB(ctx).Select("Gitee", "Uuid").Create(user).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil, kerrors.Conflict("github conflict", fmt.Sprintf("gitee(%v)", gitee))
		} else {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to create a user: gitee(%v)", gitee))
		}
	}

	return &biz.User{
		Uuid:   uuid,
		Github: gitee,
	}, nil
}

func (r *authRepo) CreateUserProfile(ctx context.Context, account, uuid string) error {
	p := &Profile{
		Uuid:     uuid,
		Username: account,
		Updated:  time.Now().Unix(),
	}
	err := r.data.DB(ctx).Select("Uuid", "Username", "Updated").Create(p).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to register a profile: uuid(%s)", uuid))
	}
	return nil
}

func (r *authRepo) CreateUserProfileWithGithub(ctx context.Context, account, github, uuid string) error {
	p := &Profile{
		Uuid:     uuid,
		Username: account,
		Github:   github,
		Updated:  time.Now().Unix(),
	}
	err := r.data.DB(ctx).Select("Uuid", "Username", "Github", "Updated").Create(p).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to register a profile: uuid(%s)", uuid))
	}
	return nil
}

func (r *authRepo) CreateUserProfileWithGitee(ctx context.Context, account, gitee, uuid string) error {
	p := &Profile{
		Uuid:     uuid,
		Username: account,
		Gitee:    gitee,
		Updated:  time.Now().Unix(),
	}
	err := r.data.DB(ctx).Select("Uuid", "Username", "Gitee", "Updated").Create(p).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to register a profile: uuid(%s)", uuid))
	}
	return nil
}

func (r *authRepo) CreateUserProfileUpdate(ctx context.Context, account, uuid string) error {
	pu := &ProfileUpdate{}
	pu.Uuid = uuid
	pu.Username = account
	pu.Updated = time.Now().Unix()
	err := r.data.DB(ctx).Select("Uuid", "Username", "Updated").Create(pu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create table: profile_update, uuid(%s)", uuid))
	}
	return nil
}

func (r *authRepo) CreateUserProfileUpdateWithGithub(ctx context.Context, account, uuid, github string) error {
	pu := &ProfileUpdate{}
	pu.Uuid = uuid
	pu.Username = account
	pu.Github = github
	pu.Updated = time.Now().Unix()
	err := r.data.DB(ctx).Select("Uuid", "Username", "Github", "Updated").Create(pu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create table: profile_update, uuid(%s)", uuid))
	}
	return nil
}

func (r *authRepo) CreateUserProfileUpdateWithGitee(ctx context.Context, account, uuid, gitee string) error {
	pu := &ProfileUpdate{}
	pu.Uuid = uuid
	pu.Username = account
	pu.Gitee = gitee
	pu.Updated = time.Now().Unix()
	err := r.data.DB(ctx).Select("Uuid", "Username", "Gitee", "Updated").Create(pu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create table: profile_update, uuid(%s)", uuid))
	}
	return nil
}

func (r *authRepo) CreateUserSearch(ctx context.Context, account, uuid string) error {
	user := &biz.UserSearchMap{
		Username:  account,
		Introduce: "",
	}
	body, err := user.MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("error marshaling document: account(%s), uuid(%s)", account, uuid))
	}

	req := esapi.IndexRequest{
		Index:      "user",
		DocumentID: uuid,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("error getting user search create response: account(%s), uuid(%s)", account, uuid))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("error parsing the response body: account(%s), uuid(%s)", account, uuid))
		} else {
			return errors.Errorf(fmt.Sprintf("error indexing document to es: reason(%v), account(%s), uuid(%s)", e, account, uuid))
		}
	}
	return nil
}

func (r *authRepo) SetUserPhone(ctx context.Context, uuid, phone string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("uuid = ?", uuid).Update("phone", phone).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return kerrors.Conflict("phone conflict", fmt.Sprintf("uuid(%s), phone(%s)", uuid, phone))
		} else {
			return errors.Wrapf(err, fmt.Sprintf("fail to set user uphone: uuid(%s), phone(%s)", uuid, phone))
		}
	}
	return nil
}

func (r *authRepo) SetUserEmail(ctx context.Context, uuid, email string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("uuid = ?", uuid).Update("email", email).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return kerrors.Conflict("email conflict", fmt.Sprintf("uuid(%s), email(%s)", uuid, email))
		} else {
			return errors.Wrapf(err, fmt.Sprintf("fail to set user email: uuid(%s), email(%s)", uuid, email))
		}
	}
	return nil
}

func (r *authRepo) SetUserWechat(ctx context.Context, uuid, wechat string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("uuid = ?", uuid).Update("wechat", wechat).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return kerrors.Conflict("wechat conflict", fmt.Sprintf("uuid(%s), wechat(%s)", uuid, wechat))
		} else {
			return errors.Wrapf(err, fmt.Sprintf("fail to set user wechat: uuid(%s), wechat(%s)", uuid, wechat))
		}
	}
	return nil
}

func (r *authRepo) SetUserQQ(ctx context.Context, uuid, qq string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("uuid = ?", uuid).Update("qq", qq).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return kerrors.Conflict("qq conflict", fmt.Sprintf("uuid(%s), qq(%s)", uuid, qq))
		} else {
			return errors.Wrapf(err, fmt.Sprintf("fail to set user qq: uuid(%s), qq(%s)", uuid, qq))
		}
	}
	return nil
}

func (r *authRepo) SetUserGitee(ctx context.Context, uuid string, gitee int32) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("uuid = ?", uuid).Update("gitee", gitee).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return kerrors.Conflict("gitee conflict", fmt.Sprintf("uuid(%s), gitee(%s)", uuid, gitee))
		} else {
			return errors.Wrapf(err, fmt.Sprintf("fail to set user gitee: uuid(%s), gitee(%s)", uuid, gitee))
		}
	}
	return nil
}

func (r *authRepo) SetUserGithub(ctx context.Context, uuid string, github int32) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("uuid = ?", uuid).Update("github", github).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return kerrors.Conflict("github conflict", fmt.Sprintf("uuid(%s), github(%s)", uuid, github))
		} else {
			return errors.Wrapf(err, fmt.Sprintf("fail to set user github: uuid(%s), github(%s)", uuid, github))
		}
	}
	return nil
}

func (r *authRepo) SetUserPassword(ctx context.Context, uuid, password string) error {
	p, err := util.HashPassword(password)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to hash password: %s", password))
	}

	err = r.data.db.WithContext(ctx).Model(&User{}).Where("uuid = ?", uuid).Update("password", p).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user password: uuid(%s), password(%s)", uuid, password))
	}
	return nil
}

func (r *authRepo) SetUserAvatar(_ context.Context, uuid, avatar string) {
	if avatar == "" {
		return
	}

	newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
	go r.data.Recover(newCtx, func(ctx context.Context) {
		r.userAvatarUpload(ctx, uuid, avatar)
	})()
}

func (r *authRepo) userAvatarUpload(ctx context.Context, uuid, avatar string) {
	avatarUrl := avatar
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, avatarUrl, nil)
	if err != nil {
		r.log.Errorf("fail to new an avatar request, uuid(%s), avatar_url(%s), error(%v)", uuid, avatarUrl, err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		r.log.Errorf("fail to get github user avatar, uuid(%s), avatar_url(%s), error(%v)", uuid, avatarUrl, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		r.log.Errorf("fail to get github user avatar, uuid(%s), avatar_url(%s), status_code(%v)", uuid, avatarUrl, res.StatusCode)
		return
	}

	key := "avatar/" + uuid + "/avatar.png"
	operation := "imageMogr2/format/webp/interlace/0/quality/100"
	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				FileId: "avatar.webp",
				Rule:   operation,
			},
		},
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   "image/png",
			XOptionHeader: &http.Header{},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	_, err = r.data.cosCli.Object.Put(context.Background(), key, res.Body, opt)
	if err != nil {
		r.log.Errorf("fail to send github user avatar to cos, uuid(%s), avatar_url(%s), err(%v)", uuid, avatarUrl, err)
	}
	return
}

func (r *authRepo) SendPhoneCode(ctx context.Context, template, phone string) error {
	code := util.RandomNumber()
	err := r.setCodeToCache(ctx, "phone_"+phone, code)
	if err != nil {
		return errors.Wrapf(err, "fail to send user phone code: phone(%s), code(%s), template(%s)", phone, code, template)
	}
	go r.data.Recover(ctx, func(ctx context.Context) {
		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			SignName:      tea.String(r.data.aliCode.signName),
			TemplateCode:  tea.String(util.GetPhoneTemplate(template)),
			PhoneNumbers:  tea.String(phone),
			TemplateParam: tea.String(fmt.Sprintf("{code:%s}", code)),
		}
		runtime := &utilService.RuntimeOptions{}
		_, tryErr := r.data.aliCode.client.SendSmsWithOptions(sendSmsRequest, runtime)
		if tryErr != nil {
			var err = &tea.SDKError{}
			if _t, ok := tryErr.(*tea.SDKError); ok {
				err = _t
			} else {
				err.Message = tea.String(tryErr.Error())
			}
			msg, _err := utilService.AssertAsString(err.Message)
			if _err != nil {
				r.log.Errorf("fail to assert message as string: phone(%s), code(%s), template(%s), err(%v)", phone, code, template, err.Error())
			}
			r.log.Errorf("fail to send user phone code: phone(%s), code(%s), template(%s), err(%v)", phone, code, template, *msg)
		}
	})()
	return nil
}

func (r *authRepo) SendEmailCode(ctx context.Context, template, email string) error {
	code := util.RandomNumber()
	err := r.setCodeToCache(ctx, "email_"+email, code)
	if err != nil {
		return errors.Wrapf(err, "fail to send email code: email(%s) code(%s) template(%s)", email, code, template)
	}
	go r.data.Recover(ctx, func(ctx context.Context) {
		m := r.data.mail.message
		d := r.data.mail.dialer
		m.SetHeader("To", email)
		m.SetHeader("Subject", "matrix 魔方技术")
		m.SetBody("text/html", util.GetEmailTemplate(template, code))
		err = d.DialAndSend(m)
		if err != nil {
			r.log.Errorf("fail to send email code: email(%s) code(%s) template(%s) error: %v", email, code, template, err.Error())
		}
	})()
	return nil
}

func (r *authRepo) VerifyPhoneCode(ctx context.Context, phone, code string) error {
	key := "phone_" + phone
	return r.verifyCode(ctx, key, code)
}

func (r *authRepo) VerifyEmailCode(ctx context.Context, email, code string) error {
	key := "email_" + email
	return r.verifyCode(ctx, key, code)
}

func (r *authRepo) verifyCode(ctx context.Context, key, code string) error {
	codeInCache, err := r.getCodeFromCache(ctx, key)
	if err != nil {
		return err
	}
	if code != codeInCache {
		return errors.Errorf("code error: code(%s)", code)
	}
	r.removeCodeFromCache(ctx, key)
	return nil
}

func (r *authRepo) VerifyPassword(ctx context.Context, account, password, mode string) (*biz.User, error) {
	var err error
	var user *biz.User
	if mode == "phone" {
		user, err = r.FindUserByPhone(ctx, account)
	} else {
		user, err = r.FindUserByEmail(ctx, account)
	}
	if err != nil {
		return nil, err
	}

	pass := util.CheckPasswordHash(password, user.Password)
	if !pass {
		return nil, errors.Errorf("password error: password(%s)", password)
	}
	return user, nil
}

func (r *authRepo) PasswordResetByPhone(ctx context.Context, phone, password string) error {
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to hash password: password(%s)", password))
	}

	err = r.data.db.WithContext(ctx).Model(&User{}).Where("phone = ?", phone).Update("password", hashPassword).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reset password: password(%s)", password))
	}
	return nil
}

func (r *authRepo) PasswordResetByEmail(ctx context.Context, email, password string) error {
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to hash password: password(%s)", password))
	}

	err = r.data.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Update("password", hashPassword).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reset password: password(%s)", password))
	}
	return nil
}

func (r *authRepo) GetCosSessionKey(_ context.Context, uuid string) (*biz.Credentials, error) {
	c := r.data.cos.client
	opt := r.data.cos.opt

	optNew := &sts.CredentialOptions{}
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(opt)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to encoder cos credential options")
	}

	err = gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(optNew)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to decode cos credential options")
	}

	for _, statement := range optNew.Policy.Statement {
		for index, _ := range statement.Resource {
			statement.Resource[index] += uuid + "/*"
		}
	}
	res, err := c.GetCredential(optNew)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get cos session key")
	}
	return &biz.Credentials{
		TmpSecretKey: res.Credentials.TmpSecretKey,
		TmpSecretID:  res.Credentials.TmpSecretID,
		SessionToken: res.Credentials.SessionToken,
		StartTime:    int64(res.StartTime),
		ExpiredTime:  int64(res.ExpiredTime),
	}, nil
}

func (r *authRepo) UnbindUserPhone(ctx context.Context, uuid string) error {
	return r.unbindUserAccount(ctx, "phone", uuid)
}

func (r *authRepo) UnbindUserEmail(ctx context.Context, uuid string) error {
	return r.unbindUserAccount(ctx, "email", uuid)
}

func (r *authRepo) UnbindUserWechat(ctx context.Context, uuid string) error {
	return r.unbindUserAccount(ctx, "wechat", uuid)
}

func (r *authRepo) UnbindUserQQ(ctx context.Context, uuid string) error {
	return r.unbindUserAccount(ctx, "qq", uuid)
}

func (r *authRepo) UnbindUserGitee(ctx context.Context, uuid string) error {
	return r.unbindUserAccount(ctx, "gitee", uuid)
}

func (r *authRepo) UnbindUserGithub(ctx context.Context, uuid string) error {
	return r.unbindUserAccount(ctx, "github", uuid)
}

func (r *authRepo) unbindUserAccount(ctx context.Context, account, uuid string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("uuid = ?", uuid).Update(account, nil).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to unbind user %s: uuid(%s)", account, uuid))
	}
	return nil
}

func (r *authRepo) GetWechatAccessToken(ctx context.Context, code string) (string, string, error) {
	url := r.data.wechat.accessTokenUrl + "?appid=" + r.data.wechat.appid + "&secret=" + r.data.wechat.secret + "&code=" + code + "&grant_type=" + r.data.wechat.grantType
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return "", "", errors.Wrapf(err, fmt.Sprintf("fail to new wechat access token request: code(%s)", code))
	}

	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", "", errors.Wrapf(err, fmt.Sprintf("fail to get wechat access token: code(%s)", code))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]interface{}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", errors.Wrapf(err, fmt.Sprintf("fail to  read body from responce data: code(%s)", code))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", "", errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: code(%s)", code))
	}
	token := data["access_token"].(string)
	openId := data["openid"].(string)
	return token, openId, nil
}

func (r *authRepo) GetQQAccessToken(ctx context.Context, code string) (string, error) {
	url := r.data.qq.accessTokenUrl + "?client_id=" + r.data.qq.clientId + "&client_secret=" + r.data.qq.clientSecret + "&code=" + code + "&grant_type=" + r.data.qq.grantType + "&redirect_uri=" + r.data.qq.redirectUri + "&fmt=json"
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to new qq access token request: code(%s)", code))
	}

	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get qq access token: code(%s)", code))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]interface{}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to read body from responce data: code(%s)", code))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: code(%s)", code))
	}
	token := data["access_token"].(string)
	return token, nil
}

func (r *authRepo) GetQQOpenId(ctx context.Context, token string) (string, error) {
	url := r.data.qq.openIdUrl + "?access_token=" + token + "&fmt=json"
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to new qq open id request: token(%s)", token))
	}

	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get qq openid: token(%s)", token))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]string{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to read body from responce data: token(%s)", token))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: token(%s)", token))
	}
	openId := data["openid"]
	return openId, nil
}

func (r *authRepo) GetGithubAccessToken(ctx context.Context, code string) (string, error) {
	url := r.data.github.accessTokenUrl + "?client_id=" + r.data.github.clientId + "&client_secret=" + r.data.github.clientSecret + "&code=" + code
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to new github access token request: code(%s)", code))
	}

	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get github access token: code(%s)", code))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]string{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to  read body from responce data: code(%s)", code))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: code(%s)", code))
	}
	token := data["access_token"]
	return token, nil
}

func (r *authRepo) GetGiteeAccessToken(ctx context.Context, code, redirectUrl string) (string, error) {
	if redirectUrl == "" {
		redirectUrl = r.data.gitee.redirectUri
	}
	url := r.data.gitee.accessTokenUrl + "?client_id=" + r.data.gitee.clientId + "&client_secret=" + r.data.gitee.clientSecret + "&code=" + code + "&redirect_uri=" + redirectUrl + "&grant_type=" + r.data.gitee.grantType
	method := "POST"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to new gitee access token request: code(%s)", code))
	}

	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get gitee access token: code(%s)", code))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]interface{}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to read body from responce data: code(%s)", code))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: code(%s)", code))
	}
	token := data["access_token"].(string)
	return token, nil
}

func (r *authRepo) GetWechatUserInfo(ctx context.Context, token, openid string) (map[string]interface{}, error) {
	url := r.data.wechat.userInfoUrl + "?access_token=" + token + "&openid=" + openid
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to new request: token(%s)", token))
	}

	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get wechat user info: token(%s)", token))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]interface{}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to  read body from responce data: token(%s)", token))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: token(%s)", token))
	}
	return data, nil
}

func (r *authRepo) GetQQUserInfo(ctx context.Context, token, openid string) (map[string]interface{}, error) {
	url := r.data.qq.userInfoUrl + "?access_token=" + token + "&oauth_consumer_key=" + r.data.qq.clientId + "&openid=" + openid
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to new request: token(%s) openid(%s)", token, openid))
	}

	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get qq user info: token(%s) openid(%s)", token, openid))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]interface{}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to  read body from responce data: token(%s) openid(%s)", token, openid))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: token(%s) openid(%s)", token, openid))
	}
	return data, nil
}

func (r *authRepo) GetGithubUserInfo(ctx context.Context, token string) (map[string]interface{}, error) {
	url := r.data.github.userInfoUrl
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to new request: token(%s)", token))
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+token)
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get github user info: token(%s)", token))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]interface{}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to  read body from responce data: token(%s)", token))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: token(%s)", token))
	}
	return data, nil
}

func (r *authRepo) GetGiteeUserInfo(ctx context.Context, token string) (map[string]interface{}, error) {
	url := r.data.gitee.userInfoUrl + "?access_token=" + token
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to new request: token(%s)", token))
	}

	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get gitee user info: token(%s)", token))
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	data := map[string]interface{}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to read body from responce data: token(%s)", token))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to unmarshal responce data: token(%s)", token))
	}
	return data, nil
}

func (r *authRepo) setCodeToCache(ctx context.Context, key, code string) error {
	err := r.data.redisCli.Set(ctx, key, code, time.Minute*2).Err()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set code to cache: redis.Set(%v), code(%s)", key, code))
	}
	return nil
}

func (r *authRepo) getCodeFromCache(ctx context.Context, key string) (string, error) {
	code, err := r.data.redisCli.Get(ctx, key).Result()
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get code from cache: redis.Get(%v)", key))
	}
	return code, nil
}

func (r *authRepo) removeCodeFromCache(ctx context.Context, key string) {
	_, err := r.data.redisCli.Del(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to delete code from cache: redis.Del(key, %v), error(%v)", key, err)
	}
}
