package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
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
	GetFollowProfileList(ctx context.Context, page int32, uuid string, followList []*Follow) ([]*UserProfile, error)
	GetFollowAchievementList(ctx context.Context, page int32, uuid string, followList []*Follow) ([]*Achievement, error)
	GetFollowedProfileList(ctx context.Context, page int32, uuid string, followedList []*Follow) ([]*UserProfile, error)
	GetFollowedAchievementList(ctx context.Context, page int32, uuid string, followedList []*Follow) ([]*Achievement, error)
	GetSearchAchievementList(ctx context.Context, searchList []*UserSearch) ([]*Achievement, error)
	GetFollowListCount(ctx context.Context, uuid string) (int32, error)
	GetFollowedList(ctx context.Context, page int32, uuid string) ([]*Follow, error)
	GetFollowedListCount(ctx context.Context, uuid string) (int32, error)
	GetUserFollow(ctx context.Context, uuid, userUuid string) (bool, error)
	GetUserFollows(ctx context.Context, uuid string) (map[string]bool, error)
	GetUserSearch(ctx context.Context, page int32, search string) ([]*UserSearch, int32, error)
	GetAvatarReview(ctx context.Context, page int32, uuid string) ([]*UserImageReview, error)
	GetCoverReview(ctx context.Context, page int32, uuid string) ([]*UserImageReview, error)
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
	repo         UserRepo
	achRepo      AchievementRepo
	creationRepo CreationRepo
	re           Recovery
	log          *log.Helper
}

func NewUserUseCase(repo UserRepo, achRepo AchievementRepo, creationRepo CreationRepo, re Recovery, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo:         repo,
		achRepo:      achRepo,
		creationRepo: creationRepo,
		re:           re,
		log:          log.NewHelper(log.With(logger, "module", "bff/biz/UserUseCase")),
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

func (r *UserUseCase) GetUserInfo(ctx context.Context) (*UserProfile, error) {
	uuid := ctx.Value("uuid").(string)
	userProfile := &UserProfile{}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		profile, err := r.repo.GetUserInfo(ctx, uuid)
		if err != nil {
			return err
		}
		userProfile.Uuid = profile.Uuid
		userProfile.Username = profile.Username
		userProfile.School = profile.School
		userProfile.Company = profile.Company
		userProfile.Job = profile.Job
		userProfile.Homepage = profile.Homepage
		userProfile.Introduce = profile.Introduce
		userProfile.Created = profile.Created
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		ach, err := r.achRepo.GetUserAchievement(ctx, uuid)
		if err != nil {
			return err
		}
		userProfile.Score = ach.Score
		userProfile.Agree = ach.Agree
		userProfile.Collect = ach.Collect
		userProfile.View = ach.View
		userProfile.Follow = ach.Follow
		userProfile.Followed = ach.Followed
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		creation, err := r.creationRepo.GetCreationUser(ctx, uuid)
		if err != nil {
			return err
		}
		userProfile.Article = creation.Article
		userProfile.Talk = creation.Talk
		userProfile.Column = creation.Column
		userProfile.Collections = creation.Collections
		userProfile.Subscribe = creation.Subscribe
		return nil
	}))
	err := g.Wait()
	if err != nil {
		return nil, err
	}
	return userProfile, err
}

func (r *UserUseCase) GetUserInfoVisitor(ctx context.Context, uuid string) (*UserProfile, error) {
	userProfile := &UserProfile{}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		profile, err := r.repo.GetUserInfo(ctx, uuid)
		if err != nil {
			return err
		}
		userProfile.Uuid = profile.Uuid
		userProfile.Username = profile.Username
		userProfile.School = profile.School
		userProfile.Company = profile.Company
		userProfile.Job = profile.Job
		userProfile.Homepage = profile.Homepage
		userProfile.Introduce = profile.Introduce
		userProfile.Created = profile.Created
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		ach, err := r.achRepo.GetUserAchievement(ctx, uuid)
		if err != nil {
			return err
		}
		userProfile.Score = ach.Score
		userProfile.Agree = ach.Agree
		userProfile.Collect = ach.Collect
		userProfile.View = ach.View
		userProfile.Follow = ach.Follow
		userProfile.Followed = ach.Followed
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		creation, err := r.creationRepo.GetCreationUserVisitor(ctx, uuid)
		if err != nil {
			return err
		}
		userProfile.Article = creation.Article
		userProfile.Talk = creation.Talk
		userProfile.Column = creation.Column
		userProfile.Collections = creation.Collections
		return nil
	}))
	err := g.Wait()
	if err != nil {
		return nil, err
	}
	return userProfile, err
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
	followList, err := r.repo.GetFollowList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		userProfileList, err := r.repo.GetFollowProfileList(ctx, page, uuid, followList)
		if err != nil {
			return err
		}
		for _, item := range userProfileList {
			for index, listItem := range followList {
				if listItem.Follow == item.Uuid {
					followList[index].Username = item.Username
					followList[index].Introduce = item.Introduce
				}
			}
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		userAchievementList, err := r.repo.GetFollowAchievementList(ctx, page, uuid, followList)
		if err != nil {
			return err
		}
		for _, item := range userAchievementList {
			for index, listItem := range followList {
				if listItem.Follow == item.Uuid {
					followList[index].Agree = item.Agree
					followList[index].View = item.View
					followList[index].FollowNum = item.Follow
					followList[index].FollowedNum = item.Followed
				}
			}
		}
		return nil
	}))
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return followList, nil
}

func (r *UserUseCase) GetFollowListCount(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetFollowListCount(ctx, uuid)
}

func (r *UserUseCase) GetFollowedList(ctx context.Context, page int32, uuid string) ([]*Follow, error) {
	followedList, err := r.repo.GetFollowedList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		userProfileList, err := r.repo.GetFollowedProfileList(ctx, page, uuid, followedList)
		if err != nil {
			return err
		}
		for _, item := range userProfileList {
			for index, listItem := range followedList {
				if listItem.Followed == item.Uuid {
					followedList[index].Username = item.Username
					followedList[index].Introduce = item.Introduce
				}
			}
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		userAchievementList, err := r.repo.GetFollowedAchievementList(ctx, page, uuid, followedList)
		if err != nil {
			return err
		}
		for _, item := range userAchievementList {
			for index, listItem := range followedList {
				if listItem.Followed == item.Uuid {
					followedList[index].Agree = item.Agree
					followedList[index].View = item.View
					followedList[index].FollowNum = item.Follow
					followedList[index].FollowedNum = item.Followed
				}
			}
		}
		return nil
	}))
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return followedList, nil
}

func (r *UserUseCase) GetFollowedListCount(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetFollowedListCount(ctx, uuid)
}

func (r *UserUseCase) GetUserFollows(ctx context.Context) (map[string]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserFollows(ctx, uuid)
}

func (r *UserUseCase) GetTimeLineUsers(ctx context.Context) ([]*TimeLIneFollows, error) {
	uuid := ctx.Value("uuid").(string)
	follows, err := r.repo.GetUserFollows(ctx, uuid)
	if err != nil {
		return nil, err
	}

	followList := make([]*TimeLIneFollows, 0)
	uuids := make([]string, 0)
	for key := range follows {
		followList = append(followList, &TimeLIneFollows{
			Uuid: key,
		})
		uuids = append(uuids, key)
	}

	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		userProfileList, err := r.repo.GetProfileList(ctx, uuids)
		if err != nil {
			return err
		}
		for _, item := range userProfileList {
			for index, listItem := range followList {
				if listItem.Uuid == item.Uuid {
					followList[index].Username = item.Username
				}
			}
		}
		return nil
	}))
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return followList, nil
}

func (r *UserUseCase) GetUserSearch(ctx context.Context, page int32, search string) ([]*UserSearch, int32, error) {
	searchList, total, err := r.repo.GetUserSearch(ctx, page, search)
	if err != nil {
		return nil, 0, err
	}
	userAchievementList, err := r.repo.GetSearchAchievementList(ctx, searchList)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range userAchievementList {
		for index, listItem := range searchList {
			if listItem.Uuid == item.Uuid {
				searchList[index].Agree = item.Agree
				searchList[index].View = item.View
				searchList[index].FollowNum = item.Follow
				searchList[index].FollowedNum = item.Followed
			}
		}
	}
	return searchList, total, nil
}

func (r *UserUseCase) GetAvatarReview(ctx context.Context, page int32) ([]*UserImageReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetAvatarReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}

func (r *UserUseCase) GetCoverReview(ctx context.Context, page int32) ([]*UserImageReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetCoverReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
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
