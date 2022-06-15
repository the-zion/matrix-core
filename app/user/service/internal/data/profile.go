package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
)

var _ biz.ProfileRepo = (*profileRepo)(nil)

var profileCacheKey = func(id string) string {
	return "profile_" + id
}

type profileRepo struct {
	data *Data
	log  *log.Helper
}

func NewProfileRepo(data *Data, logger log.Logger) biz.ProfileRepo {
	return &profileRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/profile")),
	}
}

// SetProfile set redis ?
func (r *profileRepo) SetProfile(ctx context.Context, id int64, sex, introduce, industry, address, profile, tag string) error {
	//p := Profile{
	//	Sex:             sex,
	//	Introduce:       introduce,
	//	Industry:        industry,
	//	Address:         address,
	//	PersonalProfile: profile,
	//	Tag:             tag,
	//}
	//err := r.data.db.WithContext(ctx).Model(&Profile{}).Where("user_id = ?", id).Select("Sex", "Introduce", "Industry", "Address", "PersonalProfile", "Tag").Updates(p).Error
	//if err != nil {
	//	return errors.Wrapf(err, fmt.Sprintf("fail to set user profile to db: profile(%v), user_id(%v)", p, id))
	//}
	return nil
}

// SetName set redis ?
func (r *profileRepo) SetName(ctx context.Context, id int64, name string) error {
	err := r.data.db.WithContext(ctx).Model(&Profile{}).Where("user_id = ?", id).Update("username", name).Error
	if err != nil {
		//r.log.Errorf("fail to set user name to db:name(%v) error(%v)", name, err)
		return errors.Wrapf(err, fmt.Sprintf("fail to set user name to db: name(%s), user_id(%v)", name, id))
	}
	return nil
}
