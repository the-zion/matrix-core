package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
)

type CommentRepo interface {
	GetLastCommentDraft(ctx context.Context, uuid string) (*CommentDraft, error)
	GetUserCommentAgree(ctx context.Context, uuid string) (map[int32]bool, error)
	GetCommentUser(ctx context.Context, uuid string) (*CommentUser, error)
	GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	GetSubCommentList(ctx context.Context, page, id int32) ([]*SubComment, error)
	GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	GetCommentListStatistic(ctx context.Context, page, creationId, creationType int32, key string, commentList []*Comment) ([]*CommentStatistic, error)
	GetSubCommentListStatistic(ctx context.Context, page, id int32, commentList []*SubComment) ([]*CommentStatistic, error)
	GetUserProfileList(ctx context.Context, page, creationId, creationType int32, key string, commentList []*Comment) (map[string]string, error)
	GetSubUserProfileList(ctx context.Context, page, id int32, commentList []*SubComment) (map[string]string, error)
	GetUserCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*Comment, error)
	GetUserSubCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*SubComment, error)
	GetUserSubCommentProfileList(ctx context.Context, page int32, key, uuid string, commentList []*SubComment) (map[string]string, error)
	GetUserCommentProfileList(ctx context.Context, page int32, key, uuid string, commentList []*Comment) (map[string]string, error)
	GetUserCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*Comment, error)
	GetUserSubCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*SubComment, error)
	GetUserCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*Comment, error)
	GetUserSubCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*SubComment, error)
	GetUserCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*Comment, error)
	GetUserSubCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*SubComment, error)
	CreateCommentDraft(ctx context.Context, uuid string) (int32, error)
	SendComment(ctx context.Context, id int32, uuid, ip string) error
	SendSubComment(ctx context.Context, id int32, uuid, ip string) error
	RemoveComment(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	RemoveSubComment(ctx context.Context, id, rootId int32, uuid, userUuid, reply string) error
	SetCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	SetSubCommentAgree(ctx context.Context, id int32, uuid, userUuid string) error
	CancelCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	CancelSubCommentAgree(ctx context.Context, id int32, uuid, userUuid string) error
}

type CommentUseCase struct {
	repo CommentRepo
	re   Recovery
	log  *log.Helper
}

func NewCommentUseCase(repo CommentRepo, re Recovery, logger log.Logger) *CommentUseCase {
	return &CommentUseCase{
		repo: repo,
		re:   re,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/CommentUseCase")),
	}
}

func (r *CommentUseCase) CreateCommentDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateCommentDraft(ctx, uuid)
}

func (r *CommentUseCase) GetUserCommentAgree(ctx context.Context) (map[int32]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserCommentAgree(ctx, uuid)
}

func (r *CommentUseCase) GetCommentUser(ctx context.Context) (*CommentUser, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetCommentUser(ctx, uuid)
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
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
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
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		userProfileMap, err := r.repo.GetUserProfileList(ctx, page, creationId, creationType, "comment_user_profile_list", commentList)
		if err != nil {
			return err
		}
		for index, listItem := range commentList {
			if value, ok := userProfileMap[listItem.Uuid]; ok {
				commentList[index].UserName = value
			}
		}
		return nil
	}))
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return commentList, nil
}

func (r *CommentUseCase) GetSubCommentList(ctx context.Context, page, id int32) ([]*SubComment, error) {
	subCommentList, err := r.repo.GetSubCommentList(ctx, page, id)
	if err != nil {
		return nil, err
	}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		commentListStatistic, err := r.repo.GetSubCommentListStatistic(ctx, page, id, subCommentList)
		if err != nil {
			return err
		}
		for _, item := range commentListStatistic {
			for index, listItem := range subCommentList {
				if listItem.Id == item.Id {
					subCommentList[index].Agree = item.Agree
				}
			}
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		userProfileMap, err := r.repo.GetSubUserProfileList(ctx, page, id, subCommentList)
		if err != nil {
			return err
		}
		for index, listItem := range subCommentList {
			if value, ok := userProfileMap[listItem.Uuid]; ok {
				subCommentList[index].UserName = value
			}
			if value, ok := userProfileMap[listItem.Reply]; ok {
				subCommentList[index].ReplyName = value
			}
		}
		return nil
	}))
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return subCommentList, nil
}

func (r *CommentUseCase) GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error) {
	commentList, err := r.repo.GetCommentListHot(ctx, page, creationId, creationType)
	if err != nil {
		return nil, err
	}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
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
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		userProfileMap, err := r.repo.GetUserProfileList(ctx, page, creationId, creationType, "comment_user_profile_list_hot", commentList)
		if err != nil {
			return err
		}
		for index, listItem := range commentList {
			if value, ok := userProfileMap[listItem.Uuid]; ok {
				commentList[index].UserName = value
			}
		}
		return nil
	}))
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserCommentArticleReplyList(ctx context.Context, page int32) ([]*Comment, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserCommentArticleReplyList(ctx, page, uuid)
}

func (r *CommentUseCase) GetUserSubCommentArticleReplyList(ctx context.Context, page int32) ([]*SubComment, error) {
	uuid := ctx.Value("uuid").(string)
	commentList, err := r.repo.GetUserSubCommentArticleReplyList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	userProfileMap, err := r.repo.GetUserSubCommentProfileList(ctx, page, "get_sub_comment_article_reply_profile_list_", uuid, commentList)
	if err != nil {
		return nil, err
	}

	for index, listItem := range commentList {
		if value, ok := userProfileMap[listItem.Reply]; ok {
			commentList[index].ReplyName = value
		}

		if value, ok := userProfileMap[listItem.RootUser]; ok {
			commentList[index].RootName = value
		}
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserCommentTalkReplyList(ctx context.Context, page int32) ([]*Comment, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserCommentTalkReplyList(ctx, page, uuid)
}

func (r *CommentUseCase) GetUserSubCommentTalkReplyList(ctx context.Context, page int32) ([]*SubComment, error) {
	uuid := ctx.Value("uuid").(string)
	commentList, err := r.repo.GetUserSubCommentTalkReplyList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	userProfileMap, err := r.repo.GetUserSubCommentProfileList(ctx, page, "get_sub_comment_talk_reply_profile_list_", uuid, commentList)
	if err != nil {
		return nil, err
	}

	for index, listItem := range commentList {
		if value, ok := userProfileMap[listItem.Reply]; ok {
			commentList[index].ReplyName = value
		}

		if value, ok := userProfileMap[listItem.RootUser]; ok {
			commentList[index].RootName = value
		}
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserCommentArticleRepliedList(ctx context.Context, page int32) ([]*Comment, error) {
	uuid := ctx.Value("uuid").(string)
	commentList, err := r.repo.GetUserCommentArticleRepliedList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	userProfileMap, err := r.repo.GetUserCommentProfileList(ctx, page, "get_comment_article_replied_profile_list_", uuid, commentList)
	if err != nil {
		return nil, err
	}

	for index, listItem := range commentList {
		if value, ok := userProfileMap[listItem.Uuid]; ok {
			commentList[index].UserName = value
		}
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserSubCommentArticleRepliedList(ctx context.Context, page int32) ([]*SubComment, error) {
	uuid := ctx.Value("uuid").(string)
	commentList, err := r.repo.GetUserSubCommentArticleRepliedList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	userProfileMap, err := r.repo.GetUserSubCommentProfileList(ctx, page, "get_sub_comment_article_replied_profile_list_", uuid, commentList)
	if err != nil {
		return nil, err
	}

	for index, listItem := range commentList {
		if value, ok := userProfileMap[listItem.Uuid]; ok {
			commentList[index].UserName = value
		}

		if value, ok := userProfileMap[listItem.Reply]; ok {
			commentList[index].ReplyName = value
		}

		if value, ok := userProfileMap[listItem.RootUser]; ok {
			commentList[index].RootName = value
		}
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserCommentTalkRepliedList(ctx context.Context, page int32) ([]*Comment, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserCommentTalkRepliedList(ctx, page, uuid)
}

func (r *CommentUseCase) GetUserSubCommentTalkRepliedList(ctx context.Context, page int32) ([]*SubComment, error) {
	uuid := ctx.Value("uuid").(string)
	commentList, err := r.repo.GetUserSubCommentTalkRepliedList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	userProfileMap, err := r.repo.GetUserSubCommentProfileList(ctx, page, "get_sub_comment_talk_replied_profile_list_", uuid, commentList)
	if err != nil {
		return nil, err
	}

	for index, listItem := range commentList {
		if value, ok := userProfileMap[listItem.Uuid]; ok {
			commentList[index].UserName = value
		}

		if value, ok := userProfileMap[listItem.Reply]; ok {
			commentList[index].ReplyName = value
		}

		if value, ok := userProfileMap[listItem.RootUser]; ok {
			commentList[index].RootName = value
		}
	}
	return commentList, nil
}

func (r *CommentUseCase) SendComment(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendComment(ctx, id, uuid, ip)
}

func (r *CommentUseCase) SendSubComment(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendSubComment(ctx, id, uuid, ip)
}

func (r *CommentUseCase) RemoveComment(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.RemoveComment(ctx, id, creationId, creationType, uuid, userUuid)
}

func (r *CommentUseCase) RemoveSubComment(ctx context.Context, id, rootId int32, uuid, reply string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.RemoveSubComment(ctx, id, rootId, uuid, userUuid, reply)
}

func (r *CommentUseCase) SetCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetCommentAgree(ctx, id, creationId, creationType, uuid, userUuid)
}

func (r *CommentUseCase) SetSubCommentAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetSubCommentAgree(ctx, id, uuid, userUuid)
}

func (r *CommentUseCase) CancelCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelCommentAgree(ctx, id, creationId, creationType, uuid, userUuid)
}

func (r *CommentUseCase) CancelSubCommentAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelSubCommentAgree(ctx, id, uuid, userUuid)
}
