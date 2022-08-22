package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"golang.org/x/sync/errgroup"
)

type AchievementRepo interface {
	SetAchievementAgree(ctx context.Context, uuid, userUuid string) error
	CancelAchievementAgree(ctx context.Context, uuid, userUuid string) error
	SetAchievementView(ctx context.Context, uuid string) error
	SetAchievementCollect(ctx context.Context, uuid string) error
	CancelAchievementCollect(ctx context.Context, uuid string) error
	SetAchievementFollow(ctx context.Context, follow, followed string) error
	SetUserMedalDbAndCache(ctx context.Context, medal, uuid string) error
	CancelAchievementFollow(ctx context.Context, follow, followed string) error
	CancelUserMedalDbAndCache(ctx context.Context, medal, uuid string) error
	AddAchievementScore(ctx context.Context, uuid string, score int32) error
	GetUserAchievement(ctx context.Context, uuid string) (int32, int32, error)
	GetUserActive(ctx context.Context, uuid string) (int32, error)
	AccessUserMedalDbAndCache(ctx context.Context, medal, uuid string) error
}

type AchievementCase struct {
	repo         AchievementRepo
	creationRepo CreationRepo
	commentRepo  CommentRepo
	re           Recovery
	log          *log.Helper
}

func NewAchievementUseCase(repo AchievementRepo, creationRepo CreationRepo, commentRepo CommentRepo, re Recovery, logger log.Logger) *AchievementCase {
	return &AchievementCase{
		repo:         repo,
		creationRepo: creationRepo,
		commentRepo:  commentRepo,
		re:           re,
		log:          log.NewHelper(log.With(logger, "module", "message/biz/achievementUseCase")),
	}
}

func (r *AchievementCase) SetAchievementAgree(ctx context.Context, uuid, userUuid string) error {
	return r.repo.SetAchievementAgree(ctx, uuid, userUuid)
}

func (r *AchievementCase) CancelAchievementAgree(ctx context.Context, uuid, userUuid string) error {
	return r.repo.CancelAchievementAgree(ctx, uuid, userUuid)
}

func (r *AchievementCase) SetAchievementView(ctx context.Context, uuid string) error {
	return r.repo.SetAchievementView(ctx, uuid)
}

func (r *AchievementCase) SetAchievementCollect(ctx context.Context, uuid string) error {
	return r.repo.SetAchievementCollect(ctx, uuid)
}

func (r *AchievementCase) SetUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	return r.repo.SetUserMedalDbAndCache(ctx, medal, uuid)
}

func (r *AchievementCase) CancelAchievementCollect(ctx context.Context, uuid string) error {
	return r.repo.CancelAchievementCollect(ctx, uuid)
}

func (r *AchievementCase) SetAchievementFollow(ctx context.Context, follow, followed string) error {
	return r.repo.SetAchievementFollow(ctx, follow, followed)
}

func (r *AchievementCase) CancelAchievementFollow(ctx context.Context, follow, followed string) error {
	return r.repo.CancelAchievementFollow(ctx, follow, followed)
}

func (r *AchievementCase) CancelUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	return r.repo.CancelUserMedalDbAndCache(ctx, medal, uuid)
}

func (r *AchievementCase) AccessUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	medalMap, err := r.getUserMedalProgress(ctx, uuid)
	if err != nil {
		return nil
	}
	switch medal {
	case "creation1":
		if medalMap["creation"] < 1 {
			return v1.ErrorAccessUserMedalFailed("creation less than 1")
		}
		break
	case "creation2":
		if medalMap["creation"] < 10 {
			return v1.ErrorAccessUserMedalFailed("creation less than 10")
		}
		break
	case "creation3":
		if medalMap["creation"] < 30 {
			return v1.ErrorAccessUserMedalFailed("creation less than 30")
		}
		break
	case "creation4":
		if medalMap["creation"] < 50 {
			return v1.ErrorAccessUserMedalFailed("creation less than 50")
		}
		break
	case "creation5":
		if medalMap["creation"] < 100 {
			return v1.ErrorAccessUserMedalFailed("creation less than 100")
		}
		break
	case "creation6":
		if medalMap["creation"] < 300 {
			return v1.ErrorAccessUserMedalFailed("creation less than 300")
		}
		break
	case "creation7":
		if medalMap["creation"] < 1000 {
			return v1.ErrorAccessUserMedalFailed("creation less than 1000")
		}
		break
	case "agree1":
		if medalMap["activeAgree"] < 30 {
			return v1.ErrorAccessUserMedalFailed("agree less than 30")
		}
		break
	case "agree2":
		if medalMap["activeAgree"] < 100 {
			return v1.ErrorAccessUserMedalFailed("agree less than 100")
		}
		break
	case "agree3":
		if medalMap["activeAgree"] < 300 {
			return v1.ErrorAccessUserMedalFailed("agree less than 300")
		}
		break
	case "agree4":
		if medalMap["agree"] < 50 {
			return v1.ErrorAccessUserMedalFailed("agree less than 50")
		}
		break
	case "agree5":
		if medalMap["agree"] < 300 {
			return v1.ErrorAccessUserMedalFailed("agree less than 300")
		}
		break
	case "agree6":
		if medalMap["agree"] < 1000 {
			return v1.ErrorAccessUserMedalFailed("agree less than 1000")
		}
		break
	case "view1":
		if medalMap["view"] < 10000 {
			return v1.ErrorAccessUserMedalFailed("view less than 10000")
		}
		break
	case "view2":
		if medalMap["view"] < 50000 {
			return v1.ErrorAccessUserMedalFailed("view less than 50000")
		}
		break
	case "view3":
		if medalMap["view"] < 100000 {
			return v1.ErrorAccessUserMedalFailed("view less than 100000")
		}
		break
	case "comment1":
		if medalMap["comment"] < 10 {
			return v1.ErrorAccessUserMedalFailed("comment less than 10")
		}
		break
	case "comment2":
		if medalMap["comment"] < 50 {
			return v1.ErrorAccessUserMedalFailed("comment less than 50")
		}
		break
	case "comment3":
		if medalMap["comment"] < 100 {
			return v1.ErrorAccessUserMedalFailed("comment less than 100")
		}
		break
	case "collect1":
		if medalMap["collect"] < 30 {
			return v1.ErrorAccessUserMedalFailed("collect less than 30")
		}
		break
	case "collect2":
		if medalMap["collect"] < 100 {
			return v1.ErrorAccessUserMedalFailed("collect less than 100")
		}
		break
	case "collect3":
		if medalMap["collect"] < 500 {
			return v1.ErrorAccessUserMedalFailed("collect less than 500")
		}
		break
	}
	return r.repo.AccessUserMedalDbAndCache(ctx, medal, uuid)
}

func (r *AchievementCase) getUserMedalProgress(ctx context.Context, uuid string) (map[string]int32, error) {
	m := make(map[string]int32, 0)
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		agree, view, err := r.repo.GetUserAchievement(ctx, uuid)
		if err != nil {
			return err
		}
		m["agree"] = agree
		m["view"] = view
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		agree, err := r.repo.GetUserActive(ctx, uuid)
		if err != nil {
			return err
		}
		m["activeAgree"] = agree
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		article, talk, collect, err := r.creationRepo.GetCreationUser(ctx, uuid)
		if err != nil {
			return err
		}
		m["creation"] = article + talk
		m["collect"] = collect
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		comment, err := r.commentRepo.GetCommentUser(ctx, uuid)
		if err != nil {
			return err
		}
		m["comment"] = comment
		return nil
	}))
	err := g.Wait()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *AchievementCase) AddAchievementScore(ctx context.Context, uuid string, score int32) error {
	return r.repo.AddAchievementScore(ctx, uuid, score)
}
