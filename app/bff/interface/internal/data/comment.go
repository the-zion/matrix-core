package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	commentV1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	userV1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
)

var _ biz.CommentRepo = (*commentRepo)(nil)

type commentRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/comment")),
		sg:   &singleflight.Group{},
	}
}

func (r *commentRepo) GetLastCommentDraft(ctx context.Context, uuid string) (*biz.CommentDraft, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("last_comment_draft_"+uuid), func() (interface{}, error) {
		reply, err := r.data.commc.GetLastCommentDraft(ctx, &commentV1.GetLastCommentDraftReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.CommentDraft{
			Id:     reply.Id,
			Status: reply.Status,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.CommentDraft), nil
}

func (r *commentRepo) GetUserCommentAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_comment_agree_"+uuid), func() (interface{}, error) {
		reply, err := r.data.commc.GetUserCommentAgree(ctx, &commentV1.GetUserCommentAgreeReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return reply.Agree, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[int32]bool), nil
}

func (r *commentRepo) GetCommentUser(ctx context.Context, uuid string) (*biz.CommentUser, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_comment_user_"+uuid), func() (interface{}, error) {
		reply, err := r.data.commc.GetCommentUser(ctx, &commentV1.GetCommentUserReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.CommentUser{
			Comment:           reply.Comment,
			ArticleReply:      reply.ArticleReply,
			ArticleReplySub:   reply.ArticleReplySub,
			TalkReply:         reply.TalkReply,
			TalkReplySub:      reply.TalkReplySub,
			ArticleReplied:    reply.ArticleReplied,
			ArticleRepliedSub: reply.ArticleRepliedSub,
			TalkReplied:       reply.TalkReplied,
			TalkRepliedSub:    reply.TalkRepliedSub,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.CommentUser), nil
}

func (r *commentRepo) GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("comment_%v_%v_%v", creationId, creationType, page), func() (interface{}, error) {
		reply := make([]*biz.Comment, 0)
		commentList, err := r.data.commc.GetCommentList(ctx, &commentV1.GetCommentListReq{
			Page:         page,
			CreationId:   creationId,
			CreationType: creationType,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.Comment {
			reply = append(reply, &biz.Comment{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Comment), nil
}

func (r *commentRepo) GetSubCommentList(ctx context.Context, page, id int32) ([]*biz.SubComment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("sub_comment_%v_%v", id, page), func() (interface{}, error) {
		reply := make([]*biz.SubComment, 0)
		subCommentList, err := r.data.commc.GetSubCommentList(ctx, &commentV1.GetSubCommentListReq{
			Page: page,
			Id:   id,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range subCommentList.Comment {
			reply = append(reply, &biz.SubComment{
				Id:    item.Id,
				Uuid:  item.Uuid,
				Reply: item.Reply,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.SubComment), nil
}

func (r *commentRepo) GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("comment_hot_%v_%v_%v", creationId, creationType, page), func() (interface{}, error) {
		reply := make([]*biz.Comment, 0)
		commentList, err := r.data.commc.GetCommentListHot(ctx, &commentV1.GetCommentListReq{
			Page:         page,
			CreationId:   creationId,
			CreationType: creationType,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.Comment {
			reply = append(reply, &biz.Comment{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Comment), nil
}

func (r *commentRepo) GetCommentListStatistic(ctx context.Context, page, creationId, creationType int32, key string, commentList []*biz.Comment) ([]*biz.CommentStatistic, error) {
	ids := make([]int32, 0)
	for _, item := range commentList {
		ids = append(ids, item.Id)
	}
	result, err, _ := r.sg.Do(fmt.Sprintf("%s_%v_%v_%v", key, creationId, creationType, page), func() (interface{}, error) {
		reply := make([]*biz.CommentStatistic, 0)
		commentListStatistic, err := r.data.commc.GetCommentListStatistic(ctx, &commentV1.GetCommentListStatisticReq{
			Ids: ids,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentListStatistic.Count {
			reply = append(reply, &biz.CommentStatistic{
				Id:      item.Id,
				Agree:   item.Agree,
				Comment: item.Comment,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CommentStatistic), nil
}

func (r *commentRepo) GetSubCommentListStatistic(ctx context.Context, page, id int32, commentList []*biz.SubComment) ([]*biz.CommentStatistic, error) {
	ids := make([]int32, 0)
	for _, item := range commentList {
		ids = append(ids, item.Id)
	}
	result, err, _ := r.sg.Do(fmt.Sprintf("sub_comment_%v_%v", id, page), func() (interface{}, error) {
		reply := make([]*biz.CommentStatistic, 0)
		commentListStatistic, err := r.data.commc.GetSubCommentListStatistic(ctx, &commentV1.GetCommentListStatisticReq{
			Ids: ids,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentListStatistic.Count {
			reply = append(reply, &biz.CommentStatistic{
				Id:    item.Id,
				Agree: item.Agree,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CommentStatistic), nil
}

func (r *commentRepo) GetUserProfileList(ctx context.Context, page, creationId, creationType int32, key string, commentList []*biz.Comment) (map[string]string, error) {
	uuids := make([]string, 0)
	set := make(map[string]bool, 0)
	for _, item := range commentList {
		if _, ok := set[item.Uuid]; !ok {
			uuids = append(uuids, item.Uuid)
			set[item.Uuid] = true
		}
	}
	result, err, _ := r.sg.Do(fmt.Sprintf("%s_%v_%v_%v", key, creationId, creationType, page), func() (interface{}, error) {
		reply := make(map[string]string, 0)
		userProfileList, err := r.data.uc.GetProfileList(ctx, &userV1.GetProfileListReq{
			Uuids: uuids,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range userProfileList.Profile {
			reply[item.Uuid] = item.Username
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]string), nil
}

func (r *commentRepo) GetSubUserProfileList(ctx context.Context, page, id int32, subCommentList []*biz.SubComment) (map[string]string, error) {
	uuids := make([]string, 0)
	set := make(map[string]bool, 0)
	for _, item := range subCommentList {
		if _, ok := set[item.Uuid]; !ok {
			uuids = append(uuids, item.Uuid)
			set[item.Uuid] = true
		}

		if _, ok := set[item.Reply]; item.Reply != "" && !ok {
			uuids = append(uuids, item.Reply)
			set[item.Reply] = true
		}
	}
	result, err, _ := r.sg.Do(fmt.Sprintf("sub_comment_user_profile_list_%v_%v", id, page), func() (interface{}, error) {
		reply := make(map[string]string, 0)
		userProfileList, err := r.data.uc.GetProfileList(ctx, &userV1.GetProfileListReq{
			Uuids: uuids,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range userProfileList.Profile {
			reply[item.Uuid] = item.Username
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]string), nil
}

func (r *commentRepo) GetUserCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_comment_article_reply_list_%v_%s", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.Comment, 0)
		commentList, err := r.data.commc.GetUserCommentArticleReplyList(ctx, &commentV1.GetUserCommentArticleReplyListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.List {
			reply = append(reply, &biz.Comment{
				Id:             item.Id,
				CreationId:     item.CreationId,
				CreationAuthor: item.CreationAuthor,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Comment), nil
}

func (r *commentRepo) GetUserSubCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_sub_comment_article_reply_list_%v_%s", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.SubComment, 0)
		commentList, err := r.data.commc.GetUserSubCommentArticleReplyList(ctx, &commentV1.GetUserSubCommentArticleReplyListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.List {
			reply = append(reply, &biz.SubComment{
				Id:             item.Id,
				CreationId:     item.CreationId,
				RootId:         item.RootId,
				ParentId:       item.ParentId,
				CreationAuthor: item.CreationAuthor,
				RootUser:       item.RootUser,
				Reply:          item.Reply,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.SubComment), nil
}

func (r *commentRepo) GetUserSubCommentProfileList(ctx context.Context, page int32, key, uuid string, subCommentList []*biz.SubComment) (map[string]string, error) {
	uuids := make([]string, 0)
	set := make(map[string]bool, 0)
	for _, item := range subCommentList {
		if _, ok := set[item.Uuid]; item.Uuid != "" && !ok {
			uuids = append(uuids, item.Uuid)
			set[item.Uuid] = true
		}

		if _, ok := set[item.RootUser]; item.RootUser != "" && !ok {
			uuids = append(uuids, item.RootUser)
			set[item.RootUser] = true
		}

		if _, ok := set[item.Reply]; item.Reply != "" && !ok {
			uuids = append(uuids, item.Reply)
			set[item.Reply] = true
		}
	}
	result, err, _ := r.sg.Do(fmt.Sprintf(key+"%s_%v", uuid, page), func() (interface{}, error) {
		reply := make(map[string]string, 0)
		userProfileList, err := r.data.uc.GetProfileList(ctx, &userV1.GetProfileListReq{
			Uuids: uuids,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range userProfileList.Profile {
			reply[item.Uuid] = item.Username
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]string), nil
}

func (r *commentRepo) GetUserCommentProfileList(ctx context.Context, page int32, key, uuid string, commentList []*biz.Comment) (map[string]string, error) {
	uuids := make([]string, 0)
	set := make(map[string]bool, 0)
	for _, item := range commentList {
		if _, ok := set[item.Uuid]; item.Uuid != "" && !ok {
			uuids = append(uuids, item.Uuid)
			set[item.Uuid] = true
		}
	}
	result, err, _ := r.sg.Do(fmt.Sprintf(key+"%s_%v", uuid, page), func() (interface{}, error) {
		reply := make(map[string]string, 0)
		userProfileList, err := r.data.uc.GetProfileList(ctx, &userV1.GetProfileListReq{
			Uuids: uuids,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range userProfileList.Profile {
			reply[item.Uuid] = item.Username
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string]string), nil
}

func (r *commentRepo) GetUserCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_comment_talk_reply_list_%v_%s", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.Comment, 0)
		commentList, err := r.data.commc.GetUserCommentTalkReplyList(ctx, &commentV1.GetUserCommentTalkReplyListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.List {
			reply = append(reply, &biz.Comment{
				Id:             item.Id,
				CreationId:     item.CreationId,
				CreationAuthor: item.CreationAuthor,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Comment), nil
}

func (r *commentRepo) GetUserSubCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_sub_comment_talk_reply_list_%v_%s", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.SubComment, 0)
		commentList, err := r.data.commc.GetUserSubCommentTalkReplyList(ctx, &commentV1.GetUserSubCommentTalkReplyListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.List {
			reply = append(reply, &biz.SubComment{
				Id:             item.Id,
				CreationId:     item.CreationId,
				RootId:         item.RootId,
				ParentId:       item.ParentId,
				CreationAuthor: item.CreationAuthor,
				RootUser:       item.RootUser,
				Reply:          item.Reply,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.SubComment), nil
}

func (r *commentRepo) GetUserCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_comment_article_replied_list_%v_%s", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.Comment, 0)
		commentList, err := r.data.commc.GetUserCommentArticleRepliedList(ctx, &commentV1.GetUserCommentArticleRepliedListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.List {
			reply = append(reply, &biz.Comment{
				Id:         item.Id,
				CreationId: item.CreationId,
				Uuid:       item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Comment), nil
}

func (r *commentRepo) GetUserSubCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_sub_comment_article_replied_list_%v_%s", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.SubComment, 0)
		commentList, err := r.data.commc.GetUserSubCommentArticleRepliedList(ctx, &commentV1.GetUserSubCommentArticleRepliedListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.List {
			reply = append(reply, &biz.SubComment{
				Id:             item.Id,
				Uuid:           item.Uuid,
				CreationId:     item.CreationId,
				RootId:         item.RootId,
				ParentId:       item.ParentId,
				CreationAuthor: item.CreationAuthor,
				RootUser:       item.RootUser,
				Reply:          item.Reply,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.SubComment), nil
}

func (r *commentRepo) GetUserCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_comment_talk_replied_list_%v_%s", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.Comment, 0)
		commentList, err := r.data.commc.GetUserCommentTalkRepliedList(ctx, &commentV1.GetUserCommentTalkRepliedListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.List {
			reply = append(reply, &biz.Comment{
				Id:         item.Id,
				CreationId: item.CreationId,
				Uuid:       item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Comment), nil
}

func (r *commentRepo) GetUserSubCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_sub_comment_talk_replied_list_%v_%s", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.SubComment, 0)
		commentList, err := r.data.commc.GetUserSubCommentTalkRepliedList(ctx, &commentV1.GetUserSubCommentTalkRepliedListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.List {
			reply = append(reply, &biz.SubComment{
				Id:             item.Id,
				Uuid:           item.Uuid,
				CreationId:     item.CreationId,
				RootId:         item.RootId,
				ParentId:       item.ParentId,
				CreationAuthor: item.CreationAuthor,
				RootUser:       item.RootUser,
				Reply:          item.Reply,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.SubComment), nil
}

func (r *commentRepo) GetCommentContentReview(ctx context.Context, page int32, uuid string) ([]*biz.CommentContentReview, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_comment_content_review_%s_%v", uuid, page), func() (interface{}, error) {
		reply := make([]*biz.CommentContentReview, 0)
		reviewReply, err := r.data.commc.GetCommentContentReview(ctx, &commentV1.GetCommentContentReviewReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range reviewReply.Review {
			reply = append(reply, &biz.CommentContentReview{
				Id:        item.Id,
				CommentId: item.CommentId,
				Comment:   item.Comment,
				Kind:      item.Kind,
				Uuid:      item.Uuid,
				CreateAt:  item.CreateAt,
				JobId:     item.JobId,
				Label:     item.Label,
				Result:    item.Result,
				Section:   item.Section,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CommentContentReview), nil
}

func (r *commentRepo) CreateCommentDraft(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.commc.CreateCommentDraft(ctx, &commentV1.CreateCommentDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}

func (r *commentRepo) SendComment(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.commc.SendComment(ctx, &commentV1.SendCommentReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) SendSubComment(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.commc.SendSubComment(ctx, &commentV1.SendSubCommentReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) RemoveComment(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.commc.RemoveComment(ctx, &commentV1.RemoveCommentReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) RemoveSubComment(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.commc.RemoveSubComment(ctx, &commentV1.RemoveSubCommentReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) SetCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	_, err := r.data.commc.SetCommentAgree(ctx, &commentV1.SetCommentAgreeReq{
		Uuid:         uuid,
		Id:           id,
		UserUuid:     userUuid,
		CreationId:   creationId,
		CreationType: creationType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) SetSubCommentAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.commc.SetSubCommentAgree(ctx, &commentV1.SetSubCommentAgreeReq{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) CancelCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	_, err := r.data.commc.CancelCommentAgree(ctx, &commentV1.CancelCommentAgreeReq{
		Uuid:         uuid,
		Id:           id,
		UserUuid:     userUuid,
		CreationId:   creationId,
		CreationType: creationType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) CancelSubCommentAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.commc.CancelSubCommentAgree(ctx, &commentV1.CancelSubCommentAgreeReq{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}
