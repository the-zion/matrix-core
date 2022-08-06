package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
)

type CommentRepo interface {
	GetLastCommentDraft(ctx context.Context, uuid string) (*CommentDraft, error)
	GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	GetCommentListStatistic(ctx context.Context, page, creationId, creationType int32, key string, commentList []*Comment) ([]*CommentStatistic, error)
	GetUserProfileList(ctx context.Context, page, creationId, creationType int32, key string, commentList []*Comment) ([]*UserProfile, error)
	CreateCommentDraft(ctx context.Context, uuid string) (int32, error)
	SendComment(ctx context.Context, id int32, uuid, ip string) error
}

type CommentUseCase struct {
	repo CommentRepo
	log  *log.Helper
}

func NewCommentUseCase(repo CommentRepo, logger log.Logger) *CommentUseCase {
	return &CommentUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/CommentUseCase")),
	}
}

func (r *CommentUseCase) CreateCommentDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateCommentDraft(ctx, uuid)
}

func (r *CommentUseCase) GetLastCommentDraft(ctx context.Context) (*CommentDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetLastCommentDraft(ctx, uuid)
}

func (r *CommentUseCase) GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error) {
	commentList, err := r.repo.GetCommentList(ctx, page, creationId, creationType)
	if err != nil {
		return nil, err
	}
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		commentListStatistic, err := r.repo.GetCommentListStatistic(ctx, page, creationId, creationType, "comment_statistic", commentList)
		if err != nil {
			return err
		}
		for _, item := range commentListStatistic {
			for index, listItem := range commentList {
				if listItem.Id == item.Id {
					commentList[index].Agree = item.Agree
					commentList[index].Comment = item.Comment
				}
			}
		}
		return nil
	})
	g.Go(func() error {
		userProfileList, err := r.repo.GetUserProfileList(ctx, page, creationId, creationType, "comment_user_profile_list", commentList)
		if err != nil {
			return err
		}
		for _, item := range userProfileList {
			for index, listItem := range commentList {
				if listItem.Uuid == item.Uuid {
					commentList[index].UserName = item.Username
				}
			}
		}
		return nil
	})
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return commentList, nil
}

func (r *CommentUseCase) GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error) {
	commentList, err := r.repo.GetCommentListHot(ctx, page, creationId, creationType)
	if err != nil {
		return nil, err
	}
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		commentListStatistic, err := r.repo.GetCommentListStatistic(ctx, page, creationId, creationType, "comment_statistic_hot", commentList)
		if err != nil {
			return err
		}
		for _, item := range commentListStatistic {
			for index, listItem := range commentList {
				if listItem.Id == item.Id {
					commentList[index].Agree = item.Agree
					commentList[index].Comment = item.Comment
				}
			}
		}
		return nil
	})
	g.Go(func() error {
		userProfileList, err := r.repo.GetUserProfileList(ctx, page, creationId, creationType, "comment_user_profile_list_hot", commentList)
		if err != nil {
			return err
		}
		for _, item := range userProfileList {
			for index, listItem := range commentList {
				if listItem.Uuid == item.Uuid {
					commentList[index].UserName = item.Username
				}
			}
		}
		return nil
	})
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return commentList, nil
}

func (r *CommentUseCase) SendComment(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendComment(ctx, id, uuid, ip)
}
