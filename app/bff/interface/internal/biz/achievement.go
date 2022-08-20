package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
)

type AchievementRepo interface {
	GetAchievementList(ctx context.Context, uuids []string) ([]*Achievement, error)
	GetUserAchievement(ctx context.Context, uuid string) (*Achievement, error)
	GetUserActive(ctx context.Context, uuid string) (*Active, error)
	GetUserMedal(ctx context.Context, uuid string) (*Medal, error)
	SetUserMedal(ctx context.Context, medal, uuid string) error
	CancelUserMedalSet(ctx context.Context, medal, uuid string) error
	AccessUserMedal(ctx context.Context, medal, uuid string) error
}

type AchievementUseCase struct {
	repo         AchievementRepo
	creationRepo CreationRepo
	commentRepo  CommentRepo
	re           Recovery
	log          *log.Helper
}

func NewAchievementUseCase(repo AchievementRepo, creationRepo CreationRepo, commentRepo CommentRepo, re Recovery, logger log.Logger) *AchievementUseCase {
	return &AchievementUseCase{
		repo:         repo,
		creationRepo: creationRepo,
		commentRepo:  commentRepo,
		re:           re,
		log:          log.NewHelper(log.With(logger, "module", "bff/biz/AchievementUseCase")),
	}
}

func (r *AchievementUseCase) GetAchievementList(ctx context.Context, uuids []string) ([]*Achievement, error) {
	return r.repo.GetAchievementList(ctx, uuids)
}

func (r *AchievementUseCase) GetUserAchievement(ctx context.Context, uuid string) (*Achievement, error) {
	return r.repo.GetUserAchievement(ctx, uuid)
}

func (r *AchievementUseCase) GetUserMedal(ctx context.Context, uuid string) (*Medal, error) {
	return r.repo.GetUserMedal(ctx, uuid)
}

func (r *AchievementUseCase) GetUserMedalProgress(ctx context.Context) (*MedalProgress, error) {
	uuid := ctx.Value("uuid").(string)
	mp := &MedalProgress{}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		achievement, err := r.repo.GetUserAchievement(ctx, uuid)
		if err != nil {
			return err
		}
		mp.Agree = achievement.Agree
		mp.View = achievement.View
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		active, err := r.repo.GetUserActive(ctx, uuid)
		if err != nil {
			return err
		}
		mp.ActiveAgree = active.Agree
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		creation, err := r.creationRepo.GetCreationUser(ctx, uuid)
		if err != nil {
			return err
		}
		mp.Article = creation.Article
		mp.Talk = creation.Talk
		mp.Collect = creation.Collect
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		comment, err := r.commentRepo.GetCommentUser(ctx, uuid)
		if err != nil {
			return err
		}
		mp.Comment = comment
		return nil
	}))
	err := g.Wait()
	if err != nil {
		return nil, err
	}
	return mp, nil
}

func (r *AchievementUseCase) SetUserMedal(ctx context.Context, medal string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SetUserMedal(ctx, medal, uuid)
}

func (r *AchievementUseCase) CancelUserMedalSet(ctx context.Context, medal string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CancelUserMedalSet(ctx, medal, uuid)
}

func (r *AchievementUseCase) AccessUserMedal(ctx context.Context, medal string) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.AccessUserMedal(ctx, medal, uuid)
}
