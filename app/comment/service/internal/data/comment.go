package data

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	creationV1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/comment/service/internal/biz"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

var _ biz.CommentRepo = (*commentRepo)(nil)

type commentRepo struct {
	data *Data
	log  *log.Helper
}

func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "comment/data/comment")),
	}
}

func (r *commentRepo) GetLastCommentDraft(ctx context.Context, uuid string) (*biz.CommentDraft, error) {
	draft := &CommentDraft{}
	err := r.data.db.WithContext(ctx).Select("id", "status").Where("uuid = ? and status = ?", uuid, 1).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("comment draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get last comment draft: uuid(%s)", uuid))
	}
	return &biz.CommentDraft{
		Id:     int32(draft.ID),
		Status: draft.Status,
	}, nil
}

func (r *commentRepo) GetUserCommentAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	exist, err := r.userCommentAgreeExist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		return r.getUserCommentAgreeFromCache(ctx, uuid)
	} else {
		return r.getUserCommentAgreeFromDb(ctx, uuid)
	}
}

func (r *commentRepo) userCommentAgreeExist(ctx context.Context, uuid string) (int32, error) {
	exist, err := r.data.redisCli.Exists(ctx, "user_comment_agree_"+uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to check if user comment agree exist from cache: uuid(%s)", uuid))
	}
	return int32(exist), nil
}

func (r *commentRepo) getUserCommentAgreeFromCache(ctx context.Context, uuid string) (map[int32]bool, error) {
	key := "user_comment_agree_" + uuid
	agreeSet, err := r.data.redisCli.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment agree from cache: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range agreeSet {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s), uuid(%s), key(%s)", id, uuid, key))
		}
		agreeMap[int32(id)] = true
	}
	return agreeMap, nil
}

func (r *commentRepo) getUserCommentAgreeFromDb(ctx context.Context, uuid string) (map[int32]bool, error) {
	list := make([]*CommentAgree, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id").Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment agree from db: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range list {
		agreeMap[item.CommentId] = true
	}
	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserCommentAgreeToCache(uuid, list)
		})()
	}
	return agreeMap, nil
}

func (r *commentRepo) setUserCommentAgreeToCache(uuid string, agreeList []*CommentAgree) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0, len(agreeList))
		key := "user_comment_agree_" + uuid
		for _, item := range agreeList {
			set = append(set, item.CommentId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user comment agree to cache: uuid(%s), agreeList(%v), error(%s) ", uuid, agreeList, err.Error())
	}
}

func (r *commentRepo) GetCommentUser(ctx context.Context, uuid string) (*biz.CommentUser, error) {
	var commentUser *biz.CommentUser
	key := "comment_user_" + uuid
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		commentUser, err = r.getCommentUserFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return commentUser, nil
	}

	commentUser, err = r.getCommentUserFromDB(ctx, uuid)
	if err != nil {
		return nil, err
	}

	newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
	go r.data.Recover(newCtx, func(ctx context.Context) {
		r.setCommentUserToCache(key, commentUser)
	})()

	return commentUser, nil
}

func (r *commentRepo) getCommentUserFromCache(ctx context.Context, key string) (*biz.CommentUser, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "comment", "article_reply", "article_reply_sub", "talk_reply", "talk_reply_sub", "article_replied", "article_replied_sub", "talk_replied", "talk_replied_sub").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment user form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _index, count := range statistic {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	return &biz.CommentUser{
		Comment:           val[0],
		ArticleReply:      val[1],
		ArticleReplySub:   val[2],
		TalkReply:         val[3],
		TalkReplySub:      val[4],
		ArticleReplied:    val[5],
		ArticleRepliedSub: val[6],
		TalkReplied:       val[7],
		TalkRepliedSub:    val[8],
	}, nil
}

func (r *commentRepo) getCommentUserFromDB(ctx context.Context, uuid string) (*biz.CommentUser, error) {
	cu := &CommentUser{}
	err := r.data.db.WithContext(ctx).Select("comment", "article_reply", "article_reply_sub", "talk_reply", "talk_reply_sub", "article_replied", "article_replied_sub", "talk_replied", "talk_replied_sub").Where("uuid = ?", uuid).First(cu).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get comment user from db: uuid(%v)", uuid))
	}
	return &biz.CommentUser{
		Comment:           cu.Comment,
		ArticleReply:      cu.ArticleReply,
		ArticleReplySub:   cu.ArticleReplySub,
		TalkReply:         cu.TalkReply,
		TalkReplySub:      cu.TalkReplySub,
		ArticleReplied:    cu.ArticleReplied,
		ArticleRepliedSub: cu.ArticleRepliedSub,
		TalkReplied:       cu.TalkReplied,
		TalkRepliedSub:    cu.TalkRepliedSub,
	}, nil
}

func (r *commentRepo) setCommentUserToCache(key string, commentUser *biz.CommentUser) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(context.Background(), key, "comment", commentUser.Comment, "article_reply", commentUser.ArticleReply, "article_reply_sub", commentUser.ArticleReplySub, "talk_reply", commentUser.TalkReply, "talk_reply_sub", commentUser.TalkReplySub, "article_replied", commentUser.ArticleReplied, "article_replied_sub", commentUser.ArticleRepliedSub, "talk_replied", commentUser.TalkReplied, "talk_replied_sub", commentUser.TalkRepliedSub)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set comment user to cache, err(%s)", err.Error())
	}
}

func (r *commentRepo) GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	comment, err := r.getCommentFromCache(ctx, page, creationId, creationType)
	if err != nil {
		return nil, err
	}

	size := len(comment)
	if size != 0 {
		return comment, nil
	}

	comment, err = r.getCommentFromDB(ctx, page, creationId, creationType)
	if err != nil {
		return nil, err
	}

	size = len(comment)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setCommentToCache(creationId, creationType, comment)
		})()
	}
	return comment, nil
}

func (r *commentRepo) GetSubCommentList(ctx context.Context, page, id int32) ([]*biz.SubComment, error) {
	subComment, err := r.getSubCommentFromCache(ctx, page, id)
	if err != nil {
		return nil, err
	}

	size := len(subComment)
	if size != 0 {
		return subComment, nil
	}

	subComment, err = r.getSubCommentFromDB(ctx, page, id)
	if err != nil {
		return nil, err
	}

	size = len(subComment)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setSubCommentToCache(id, subComment)
		})()
	}
	return subComment, nil
}

func (r *commentRepo) getCommentFromCache(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("comment_%v_%v", creationId, creationType)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		comment = append(comment, &biz.Comment{
			CommentId: int32(id),
			Uuid:      member[1],
		})
	}
	return comment, nil
}

func (r *commentRepo) getSubCommentFromCache(ctx context.Context, page, id int32) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("sub_comment_%v", id)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get sub comment from cache: key(%s), page(%v)", key, page))
	}

	subComment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		subComment = append(subComment, &biz.SubComment{
			CommentId: int32(id),
			Uuid:      member[1],
			Reply:     member[2],
		})
	}
	return subComment, nil
}

func (r *commentRepo) getCommentFromDB(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Comment, 0)
	err := r.data.db.WithContext(ctx).Select("updated_at", "comment_id", "uuid").Where("creation_id = ? and creation_type = ?", creationId, creationType).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.Comment{
			UpdatedAt: int32(item.UpdatedAt.Unix()),
			CommentId: item.CommentId,
			Uuid:      item.Uuid,
		})
	}
	return comment, nil
}

func (r *commentRepo) getSubCommentFromDB(ctx context.Context, page, id int32) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*SubComment, 0)
	err := r.data.db.WithContext(ctx).Select("updated_at", "comment_id", "uuid", "reply").Where("root_id = ?", id).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get sub comment from db: page(%v)", page))
	}

	subComment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		subComment = append(subComment, &biz.SubComment{
			CommentId: item.CommentId,
			UpdatedAt: int32(item.UpdatedAt.Unix()),
			Uuid:      item.Uuid,
			Reply:     item.Reply,
		})
	}
	return subComment, nil
}

func (r *commentRepo) setCommentToCache(creationId, creationType int32, comment []*biz.Comment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("comment_%v_%v", creationId, creationType)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set comment to cache: comment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) setSubCommentToCache(id int32, subComment []*biz.SubComment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(subComment))
		key := fmt.Sprintf("sub_comment_%v", id)
		for _, item := range subComment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + item.Uuid + "%" + item.Reply,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set sub comment to cache: comment(%v), err(%v)", subComment, err)
	}
}

func (r *commentRepo) GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	comment, err := r.getCommentHotFromCache(ctx, page, creationId, creationType)
	if err != nil {
		return nil, err
	}

	size := len(comment)
	if size != 0 {
		return comment, nil
	}

	comment, err = r.getCommentHotFromDB(ctx, page, creationId, creationType)
	if err != nil {
		return nil, err
	}

	size = len(comment)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setCommentHotToCache(creationId, creationType, comment)
		})()
	}
	return comment, nil
}

func (r *commentRepo) getCommentHotFromCache(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("comment_%v_%v_hot", creationId, creationType)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		comment = append(comment, &biz.Comment{
			CommentId: int32(id),
			Uuid:      member[1],
		})
	}
	return comment, nil
}

func (r *commentRepo) getCommentHotFromDB(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Comment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "uuid", "agree").Where("creation_id = ? and creation_type = ?", creationId, creationType).Order("agree desc, comment_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment hot from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.Comment{
			CommentId: item.CommentId,
			Uuid:      item.Uuid,
			Agree:     item.Agree,
		})
	}
	return comment, nil
}

func (r *commentRepo) setCommentHotToCache(creationId, creationType int32, comment []*biz.Comment) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("comment_%v_%v_hot", creationId, creationType)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		pipe.Expire(context.Background(), key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set comment hot to cache: comment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetCommentListStatistic(ctx context.Context, ids []int32) ([]*biz.CommentStatistic, error) {
	exists, unExists, err := r.commentStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

	commentListStatistic := make([]*biz.CommentStatistic, 0, cap(exists))
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(exists) == 0 {
			return nil
		}
		return r.getCommentStatisticFromCache(ctx, exists, &commentListStatistic)
	}))
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(unExists) == 0 {
			return nil
		}
		return r.getCommentStatisticFromDb(ctx, unExists, &commentListStatistic)
	}))

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return commentListStatistic, nil
}

func (r *commentRepo) GetSubCommentListStatistic(ctx context.Context, ids []int32) ([]*biz.CommentStatistic, error) {
	exists, unExists, err := r.commentStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

	commentListStatistic := make([]*biz.CommentStatistic, 0, cap(exists))
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(exists) == 0 {
			return nil
		}
		return r.getSubCommentStatisticFromCache(ctx, exists, &commentListStatistic)
	}))
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(unExists) == 0 {
			return nil
		}
		return r.getSubCommentStatisticFromDb(ctx, unExists, &commentListStatistic)
	}))

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return commentListStatistic, nil
}

func (r *commentRepo) GetUserCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	commentList, err := r.getUserCommentArticleReplyListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserCommentArticleReplyListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserCommentArticleReplyListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserCommentArticleReplyListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_comment_article_reply_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment article reply list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		comment = append(comment, &biz.Comment{
			CommentId:      int32(id),
			CreationId:     int32(creationId),
			CreationAuthor: member[2],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserCommentArticleReplyListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Comment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "creation_id", "creation_author").Where("uuid = ? and creation_type = ?", uuid, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment article reply list from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.Comment{
			CommentId:      item.CommentId,
			UpdatedAt:      int32(item.UpdatedAt.Unix()),
			CreationId:     item.CreationId,
			CreationAuthor: item.CreationAuthor,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserCommentArticleReplyListToCache(uuid string, comment []*biz.Comment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_comment_article_reply_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + item.CreationAuthor,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user comment article reply list to cache: comment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserSubCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	commentList, err := r.getUserSubCommentArticleReplyListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserSubCommentArticleReplyListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserSubCommentArticleReplyListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserSubCommentArticleReplyListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_sub_comment_article_reply_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment article reply list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		rootId, err := strconv.ParseInt(member[2], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[2]))
		}
		parentId, err := strconv.ParseInt(member[3], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[3]))
		}
		comment = append(comment, &biz.SubComment{
			CommentId:      int32(id),
			CreationId:     int32(creationId),
			RootId:         int32(rootId),
			ParentId:       int32(parentId),
			CreationAuthor: member[4],
			RootUser:       member[5],
			Reply:          member[6],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserSubCommentArticleReplyListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*SubComment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "creation_id", "root_id", "parent_id", "creation_author", "root_user", "reply").Where("uuid = ? and creation_type = ?", uuid, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment article reply list from db: page(%v)", page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.SubComment{
			CommentId:      item.CommentId,
			UpdatedAt:      int32(item.UpdatedAt.Unix()),
			CreationId:     item.CreationId,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserSubCommentArticleReplyListToCache(uuid string, comment []*biz.SubComment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_sub_comment_article_reply_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + strconv.Itoa(int(item.RootId)) + "%" + strconv.Itoa(int(item.ParentId)) + "%" + item.CreationAuthor + "%" + item.RootUser + "%" + item.Reply,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user sub comment article reply list to cache: subComment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	commentList, err := r.getUserCommentTalkReplyListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserCommentTalkReplyListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserCommentTalkReplyListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserCommentTalkReplyListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_comment_talk_reply_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment talk reply list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		comment = append(comment, &biz.Comment{
			CommentId:      int32(id),
			CreationId:     int32(creationId),
			CreationAuthor: member[2],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserCommentTalkReplyListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Comment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "creation_id", "creation_author").Where("uuid = ? and creation_type = ?", uuid, 3).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment talk reply list from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.Comment{
			CommentId:      item.CommentId,
			UpdatedAt:      int32(item.UpdatedAt.Unix()),
			CreationId:     item.CreationId,
			CreationAuthor: item.CreationAuthor,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserCommentTalkReplyListToCache(uuid string, comment []*biz.Comment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_comment_talk_reply_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + item.CreationAuthor,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user comment talk reply list to cache: comment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserSubCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	commentList, err := r.getUserSubCommentTalkReplyListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserSubCommentTalkReplyListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserSubCommentTalkReplyListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserSubCommentTalkReplyListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_sub_comment_talk_reply_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment talk reply list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		rootId, err := strconv.ParseInt(member[2], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[2]))
		}
		parentId, err := strconv.ParseInt(member[3], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[3]))
		}
		comment = append(comment, &biz.SubComment{
			CommentId:      int32(id),
			CreationId:     int32(creationId),
			RootId:         int32(rootId),
			ParentId:       int32(parentId),
			CreationAuthor: member[4],
			RootUser:       member[5],
			Reply:          member[6],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserSubCommentTalkReplyListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*SubComment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "creation_id", "root_id", "parent_id", "creation_author", "root_user", "reply").Where("uuid = ? and creation_type = ?", uuid, 3).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment talk reply list from db: page(%v)", page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.SubComment{
			CommentId:      item.CommentId,
			UpdatedAt:      int32(item.UpdatedAt.Unix()),
			CreationId:     item.CreationId,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserSubCommentTalkReplyListToCache(uuid string, comment []*biz.SubComment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_sub_comment_talk_reply_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + strconv.Itoa(int(item.RootId)) + "%" + strconv.Itoa(int(item.ParentId)) + "%" + item.CreationAuthor + "%" + item.RootUser + "%" + item.Reply,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user sub comment talk reply list to cache: subComment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	commentList, err := r.getUserCommentArticleRepliedListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserCommentArticleRepliedListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserCommentArticleRepliedListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserCommentArticleRepliedListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_comment_article_replied_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment article replied list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		comment = append(comment, &biz.Comment{
			CommentId:  int32(id),
			CreationId: int32(creationId),
			Uuid:       member[2],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserCommentArticleRepliedListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Comment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "creation_id", "uuid").Where("creation_author = ? and creation_type = ?", uuid, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment article replied list from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.Comment{
			CommentId:  item.CommentId,
			UpdatedAt:  int32(item.UpdatedAt.Unix()),
			CreationId: item.CreationId,
			Uuid:       item.Uuid,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserCommentArticleRepliedListToCache(uuid string, comment []*biz.Comment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_comment_article_replied_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user comment article replied list to cache: comment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserSubCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	commentList, err := r.getUserSubCommentArticleRepliedListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserSubCommentArticleRepliedListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserSubCommentArticleRepliedListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserSubCommentArticleRepliedListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_sub_comment_article_replied_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment article replied list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		rootId, err := strconv.ParseInt(member[2], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[2]))
		}
		parentId, err := strconv.ParseInt(member[3], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[3]))
		}
		comment = append(comment, &biz.SubComment{
			CommentId:      int32(id),
			CreationId:     int32(creationId),
			RootId:         int32(rootId),
			ParentId:       int32(parentId),
			Uuid:           member[4],
			CreationAuthor: member[5],
			RootUser:       member[6],
			Reply:          member[7],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserSubCommentArticleRepliedListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*SubComment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "uuid", "creation_id", "root_id", "parent_id", "creation_author", "root_user", "reply").Where("(root_user = ? or reply = ?) and creation_type = ?", uuid, uuid, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment article replied list from db: page(%v)", page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.SubComment{
			CommentId:      item.CommentId,
			UpdatedAt:      int32(item.UpdatedAt.Unix()),
			Uuid:           item.Uuid,
			CreationId:     item.CreationId,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserSubCommentArticleRepliedListToCache(uuid string, comment []*biz.SubComment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_sub_comment_article_replied_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + strconv.Itoa(int(item.RootId)) + "%" + strconv.Itoa(int(item.ParentId)) + "%" + item.Uuid + "%" + item.CreationAuthor + "%" + item.RootUser + "%" + item.Reply,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user sub comment article replied list to cache: subComment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	commentList, err := r.getUserCommentTalkRepliedListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserCommentTalkRepliedListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserCommentTalkRepliedListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserCommentTalkRepliedListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_comment_talk_replied_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment talk replied list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		comment = append(comment, &biz.Comment{
			CommentId:  int32(id),
			CreationId: int32(creationId),
			Uuid:       member[2],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserCommentTalkRepliedListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Comment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "creation_id", "uuid").Where("creation_author = ? and creation_type = ?", uuid, 3).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment talk replied list from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.Comment{
			CommentId:  item.CommentId,
			UpdatedAt:  int32(item.UpdatedAt.Unix()),
			CreationId: item.CreationId,
			Uuid:       item.Uuid,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserCommentTalkRepliedListToCache(uuid string, comment []*biz.Comment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_comment_talk_replied_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user comment talk replied list to cache: comment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserSubCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	commentList, err := r.getUserSubCommentTalkRepliedListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserSubCommentTalkRepliedListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserSubCommentTalkRepliedListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserSubCommentTalkRepliedListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_sub_comment_talk_replied_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment talk replied list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		rootId, err := strconv.ParseInt(member[2], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[2]))
		}
		parentId, err := strconv.ParseInt(member[3], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[3]))
		}
		comment = append(comment, &biz.SubComment{
			CommentId:      int32(id),
			CreationId:     int32(creationId),
			RootId:         int32(rootId),
			ParentId:       int32(parentId),
			Uuid:           member[4],
			CreationAuthor: member[5],
			RootUser:       member[6],
			Reply:          member[7],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserSubCommentTalkRepliedListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*SubComment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "uuid", "creation_id", "root_id", "parent_id", "creation_author", "root_user", "reply").Where("(root_user = ? or reply = ?) and creation_type = ?", uuid, uuid, 3).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment talk replied list from db: page(%v)", page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.SubComment{
			CommentId:      item.CommentId,
			UpdatedAt:      int32(item.UpdatedAt.Unix()),
			Uuid:           item.Uuid,
			CreationId:     item.CreationId,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserSubCommentTalkRepliedListToCache(uuid string, comment []*biz.SubComment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_sub_comment_talk_replied_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + strconv.Itoa(int(item.RootId)) + "%" + strconv.Itoa(int(item.ParentId)) + "%" + item.Uuid + "%" + item.CreationAuthor + "%" + item.RootUser + "%" + item.Reply,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user sub comment talk replied list to cache: subComment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserCommentRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	commentList, err := r.getUserCommentRepliedListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserCommentRepliedListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserCommentRepliedListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserCommentRepliedListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_comment_replied_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment replied list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		creationType, err := strconv.ParseInt(member[2], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[2]))
		}
		comment = append(comment, &biz.Comment{
			CommentId:    int32(id),
			CreationId:   int32(creationId),
			CreationType: int32(creationType),
			Uuid:         member[3],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserCommentRepliedListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Comment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Comment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "creation_id", "creation_type", "uuid").Where("creation_author = ?", uuid).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment replied list from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.Comment{
			CommentId:    item.CommentId,
			UpdatedAt:    int32(item.UpdatedAt.Unix()),
			CreationId:   item.CreationId,
			CreationType: item.CreationType,
			Uuid:         item.Uuid,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserCommentRepliedListToCache(uuid string, comment []*biz.Comment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_comment_replied_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + strconv.Itoa(int(item.CreationType)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user comment replied list to cache: comment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetUserSubCommentRepliedList(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	commentList, err := r.getUserSubCommentRepliedListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(commentList)
	if size != 0 {
		return commentList, nil
	}

	commentList, err = r.getUserSubCommentRepliedListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(commentList)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserSubCommentRepliedListToCache(uuid, commentList)
		})()
	}
	return commentList, nil
}

func (r *commentRepo) getUserSubCommentRepliedListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	key := fmt.Sprintf("user_sub_comment_replied_list_%s", uuid)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment replied list from cache: key(%s), page(%v)", key, page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[1]))
		}
		creationType, err := strconv.ParseInt(member[2], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[2]))
		}
		rootId, err := strconv.ParseInt(member[3], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[3]))
		}
		parentId, err := strconv.ParseInt(member[4], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[4]))
		}
		comment = append(comment, &biz.SubComment{
			CommentId:      int32(id),
			CreationId:     int32(creationId),
			CreationType:   int32(creationType),
			RootId:         int32(rootId),
			ParentId:       int32(parentId),
			Uuid:           member[5],
			CreationAuthor: member[6],
			RootUser:       member[7],
			Reply:          member[8],
		})
	}
	return comment, nil
}

func (r *commentRepo) getUserSubCommentRepliedListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.SubComment, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*SubComment, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "updated_at", "uuid", "creation_id", "creation_type", "root_id", "parent_id", "creation_author", "root_user", "reply").Where("root_user = ? or reply = ?", uuid, uuid).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user sub comment replied list from db: page(%v)", page))
	}

	comment := make([]*biz.SubComment, 0, len(list))
	for _, item := range list {
		comment = append(comment, &biz.SubComment{
			CommentId:      item.CommentId,
			UpdatedAt:      int32(item.UpdatedAt.Unix()),
			Uuid:           item.Uuid,
			CreationId:     item.CreationId,
			CreationType:   item.CreationType,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
		})
	}
	return comment, nil
}

func (r *commentRepo) setUserSubCommentRepliedListToCache(uuid string, comment []*biz.SubComment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(comment))
		key := fmt.Sprintf("user_sub_comment_replied_list_%s", uuid)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.UpdatedAt),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + strconv.Itoa(int(item.CreationType)) + "%" + strconv.Itoa(int(item.RootId)) + "%" + strconv.Itoa(int(item.ParentId)) + "%" + item.Uuid + "%" + item.CreationAuthor + "%" + item.RootUser + "%" + item.Reply,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user sub comment replied list to cache: subComment(%v), err(%v)", comment, err)
	}
}

func (r *commentRepo) GetCommentContentReview(ctx context.Context, page int32, uuid string) ([]*biz.TextReview, error) {
	key := "comment_content_irregular_" + uuid
	review, err := r.getCommentContentReviewFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(review)
	if size != 0 {
		return review, nil
	}

	review, err = r.getCommentContentReviewFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(review)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setCommentContentReviewToCache(key, review)
		})()
	}
	return review, nil
}

func (r *commentRepo) getCommentContentReviewFromCache(ctx context.Context, page int32, key string) ([]*biz.TextReview, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.LRange(ctx, key, index*20, index*20+19).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment content irregular list from cache: key(%s), page(%v)", key, page))
	}

	review := make([]*biz.TextReview, 0, len(list))
	for _index, item := range list {
		var textReview = &biz.TextReview{}
		err = textReview.UnmarshalJSON([]byte(item))
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: contentReview(%v)", item))
		}
		review = append(review, &biz.TextReview{
			Id:        int32(_index+1) + (page-1)*20,
			CommentId: textReview.CommentId,
			Comment:   textReview.Comment,
			Kind:      textReview.Kind,
			Uuid:      textReview.Uuid,
			CreateAt:  textReview.CreateAt,
			JobId:     textReview.JobId,
			Label:     textReview.Label,
			Result:    textReview.Result,
			Section:   textReview.Section,
		})
	}
	return review, nil
}

func (r *commentRepo) getCommentContentReviewFromDB(ctx context.Context, page int32, uuid string) ([]*biz.TextReview, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*CommentContentReview, 0)
	err := r.data.db.WithContext(ctx).Select("comment_id", "kind", "comment", "uuid", "job_id", "create_at", "label", "result", "section").Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment content review from db: page(%v), uuid(%s)", page, uuid))
	}

	review := make([]*biz.TextReview, 0, len(list))
	for _index, item := range list {
		review = append(review, &biz.TextReview{
			Id:        int32(_index+1) + (page-1)*20,
			CommentId: item.CommentId,
			Kind:      item.Kind,
			Comment:   item.Comment,
			Uuid:      item.Uuid,
			JobId:     item.JobId,
			CreateAt:  item.CreatedAt.Format("2006-01-02"),
			Label:     item.Label,
			Result:    item.Result,
			Section:   item.Section,
		})
	}
	return review, nil
}

func (r *commentRepo) setCommentContentReviewToCache(key string, review []*biz.TextReview) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		list := make([]interface{}, 0, len(review))
		for _, item := range review {
			m, err := item.MarshalJSON()
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to marshal avatar review: contentReview(%v)", review))
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set comment content review to cache: contentReview(%v), err(%v)", review, err)
	}
}

func (r *commentRepo) GetRootComment(ctx context.Context, id int32) (*biz.Comment, error) {
	sc := &Comment{}
	err := r.data.DB(ctx).Select("creation_id", "creation_type", "creation_author", "uuid").Where("comment_id = ?", id).First(sc).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get root comment: rootId(%v)", id))
	}
	return &biz.Comment{
		CreationId:     sc.CreationId,
		CreationType:   sc.CreationType,
		CreationAuthor: sc.CreationAuthor,
		Uuid:           sc.Uuid,
	}, nil
}

func (r *commentRepo) GetSubComment(ctx context.Context, id int32) (*biz.SubComment, error) {
	sc := &SubComment{}
	err := r.data.DB(ctx).Select("creation_id", "root_id", "parent_id", "creation_type", "creation_author", "root_user", "uuid", "reply").Where("comment_id = ?", id).First(sc).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get sun comment: rootId(%v)", id))
	}
	return &biz.SubComment{
		CreationId:     sc.CreationId,
		RootId:         sc.RootId,
		ParentId:       sc.ParentId,
		CreationType:   sc.CreationType,
		CreationAuthor: sc.CreationAuthor,
		RootUser:       sc.RootUser,
		Uuid:           sc.Uuid,
		Reply:          sc.Reply,
	}, nil
}

func (r *commentRepo) GetParentCommentUserId(ctx context.Context, id, rootId int32) (string, error) {
	sc := &SubComment{}
	err := r.data.DB(ctx).Select("uuid").Where("comment_id = ? and root_id = ?", id, rootId).First(sc).Error
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get parent comment user id: parentId(%v)", id))
	}
	return sc.Uuid, nil
}

func (r *commentRepo) GetSubCommentReply(ctx context.Context, id int32, uuid string) (string, error) {
	sc := &SubComment{}
	err := r.data.DB(ctx).Select("reply").Where("comment_id = ? and uuid = ?", id, uuid).First(sc).Error
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get sub comment reply: parentId(%v)", id))
	}
	return sc.Reply, nil
}

func (r *commentRepo) GetArticleAuthor(ctx context.Context, id int32) (string, error) {
	result, err := r.data.cc.GetArticleStatistic(ctx, &creationV1.GetArticleStatisticReq{
		Id: id,
	})
	if err != nil {
		return "", err
	}
	return result.Uuid, nil
}

func (r *commentRepo) GetTalkAuthor(ctx context.Context, id int32) (string, error) {
	result, err := r.data.cc.GetTalkStatistic(ctx, &creationV1.GetTalkStatisticReq{
		Id: id,
	})
	if err != nil {
		return "", err
	}
	return result.Uuid, nil
}

func (r *commentRepo) commentStatisticExist(ctx context.Context, ids []int32) ([]int32, []int32, error) {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range ids {
			pipe.Exists(ctx, "comment_"+strconv.Itoa(int(item)))
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if comment statistic exist from cache: ids(%v)", ids))
	}

	exists := make([]int32, 0, len(cmd))
	unExists := make([]int32, 0, len(cmd))
	for index, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		if exist == 1 {
			exists = append(exists, ids[index])
		} else {
			unExists = append(unExists, ids[index])
		}
	}
	return exists, unExists, nil
}

func (r *commentRepo) getCommentStatisticFromCache(ctx context.Context, exists []int32, commentListStatistic *[]*biz.CommentStatistic) error {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, id := range exists {
			pipe.HMGet(ctx, "comment_"+strconv.Itoa(int(id)), "agree", "comment")
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get comment list statistic from cache: ids(%v)", exists))
	}

	for index, item := range cmd {
		val := []int32{0, 0}
		for _index, count := range item.(*redis.SliceCmd).Val() {
			if count == nil {
				break
			}
			num, err := strconv.ParseInt(count.(string), 10, 32)
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
			}
			val[_index] = int32(num)
		}
		*commentListStatistic = append(*commentListStatistic, &biz.CommentStatistic{
			CommentId: exists[index],
			Agree:     val[0],
			Comment:   val[1],
		})
	}
	return nil
}

func (r *commentRepo) getSubCommentStatisticFromCache(ctx context.Context, exists []int32, commentListStatistic *[]*biz.CommentStatistic) error {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, id := range exists {
			pipe.HMGet(ctx, "comment_"+strconv.Itoa(int(id)), "agree")
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get sub comment list statistic from cache: ids(%v)", exists))
	}

	for index, item := range cmd {
		count := item.(*redis.SliceCmd).Val()[0]
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val := int32(num)
		*commentListStatistic = append(*commentListStatistic, &biz.CommentStatistic{
			CommentId: exists[index],
			Agree:     val,
		})
	}
	return nil
}

func (r *commentRepo) getCommentStatisticFromDb(ctx context.Context, unExists []int32, commentListStatistic *[]*biz.CommentStatistic) error {
	list := make([]*Comment, 0, cap(unExists))
	err := r.data.db.WithContext(ctx).Select("comment_id", "agree", "comment").Where("comment_id IN ?", unExists).Find(&list).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get comment statistic list from db: ids(%v)", unExists))
	}

	for _, item := range list {
		*commentListStatistic = append(*commentListStatistic, &biz.CommentStatistic{
			CommentId: item.CommentId,
			Agree:     item.Agree,
			Comment:   item.Comment,
		})
	}

	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setCommentStatisticToCache(list)
		})()
	}

	return nil
}

func (r *commentRepo) getSubCommentStatisticFromDb(ctx context.Context, unExists []int32, commentListStatistic *[]*biz.CommentStatistic) error {
	list := make([]*SubComment, 0, cap(unExists))
	err := r.data.db.WithContext(ctx).Select("comment_id", "agree").Where("comment_id IN ?", unExists).Find(&list).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get sub comment statistic list from db: ids(%v)", unExists))
	}

	for _, item := range list {
		*commentListStatistic = append(*commentListStatistic, &biz.CommentStatistic{
			CommentId: item.CommentId,
			Agree:     item.Agree,
		})
	}

	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setSubCommentStatisticToCache(list)
		})()
	}

	return nil
}

func (r *commentRepo) setCommentStatisticToCache(commentList []*Comment) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, item := range commentList {
			key := "comment_" + strconv.Itoa(int(item.CommentId))
			pipe.HSetNX(context.Background(), key, "agree", item.Agree)
			pipe.HSetNX(context.Background(), key, "comment", item.Comment)
			pipe.Expire(context.Background(), key, time.Minute*30)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set comment statistic to cache, err(%s)", err.Error())
	}
}

func (r *commentRepo) setSubCommentStatisticToCache(commentList []*SubComment) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, item := range commentList {
			key := "comment_" + strconv.Itoa(int(item.CommentId))
			pipe.HSetNX(context.Background(), key, "agree", item.Agree)
			pipe.Expire(context.Background(), key, time.Minute*30)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set sub comment statistic to cache, err(%s)", err.Error())
	}
}

func (r *commentRepo) CreateCommentDraft(ctx context.Context, uuid string) (int32, error) {
	draft := &CommentDraft{
		Uuid: uuid,
	}
	err := r.data.DB(ctx).Select("Uuid").Create(draft).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to create a comment draft: uuid(%s)", uuid))
	}
	return int32(draft.ID), nil
}

func (r *commentRepo) CreateCommentFolder(ctx context.Context, id int32, uuid string) error {
	name := "comment/" + uuid + "/" + strconv.Itoa(int(id)) + "/"
	_, err := r.data.cosCli.Object.Put(ctx, name, strings.NewReader(""), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a comment folder: id(%v)", id))
	}
	return nil
}

func (r *commentRepo) SendComment(ctx context.Context, id int32, uuid string) (*biz.CommentDraft, error) {
	cd := &CommentDraft{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&CommentDraft{}).Where("id = ? and uuid = ? and status = ?", id, uuid, 1).Updates(cd).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 2: uuid(%s), id(%v)", uuid, id))
	}
	return &biz.CommentDraft{
		Uuid: uuid,
		Id:   id,
	}, nil
}

func (r *commentRepo) SetRecord(ctx context.Context, id int32, uuid, ip string) error {
	record := &Record{
		CommonId: id,
		Uuid:     uuid,
		Ip:       ip,
	}
	err := r.data.DB(ctx).Create(record).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add an record to db: id(%v), uuid(%s), ip(%s)", id, uuid, ip))
	}
	return nil
}

func (r *commentRepo) SetUserCommentAgreeToCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_comment_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "16415c44ae0544f8c7f85d841521813d41c35994", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user comment agree to cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *commentRepo) CancelUserCommentAgreeFromCache(ctx context.Context, id int32, userUuid string) error {
	_, err := r.data.redisCli.SRem(ctx, "user_comment_agree_"+userUuid, id).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user comment agree from cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *commentRepo) SetCommentAgree(ctx context.Context, id int32, uuid string) error {
	c := Comment{}
	err := r.data.DB(ctx).Model(&c).Where("comment_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add comment agree: id(%v)", id))
	}
	return nil
}

func (r *commentRepo) SetSubCommentAgree(ctx context.Context, id int32, uuid string) error {
	sc := SubComment{}
	err := r.data.DB(ctx).Model(&sc).Where("comment_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add sub comment agree: id(%v)", id))
	}
	return nil
}

func (r *commentRepo) SetCommentComment(ctx context.Context, id int32) error {
	c := Comment{}
	err := r.data.DB(ctx).Model(&c).Where("comment_id = ?", id).Update("comment", gorm.Expr("comment + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add comment comment: id(%v)", id))
	}
	return nil
}

func (r *commentRepo) SetUserCommentAgree(ctx context.Context, id int32, userUuid string) error {
	ca := &CommentAgree{
		CommentId: id,
		Uuid:      userUuid,
		Status:    1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(ca).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set a comment agree to db: comment_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *commentRepo) SetCommentAgreeToCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	hotKey := fmt.Sprintf("comment_%v_%v_hot", creationId, creationType)
	statisticKey := fmt.Sprintf("comment_%v", id)
	userKey := fmt.Sprintf("user_comment_agree_%s", userUuid)

	keys := []string{hotKey, statisticKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), id}
	_, err := r.data.redisCli.EvalSha(ctx, "d9069cfea8db1d598aca74d388d8cf311099c769", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to update(add) comment cache: id(%v), creationId(%v), creationType(%v), uuid(%s), userUuid(%s) ", id, creationId, creationType, uuid, userUuid))
	}
	return nil
}

func (r *commentRepo) SetSubCommentAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("comment_%v", id)
	userKey := fmt.Sprintf("user_comment_agree_%s", userUuid)

	keys := []string{statisticKey, userKey}
	values := []interface{}{id}
	_, err := r.data.redisCli.EvalSha(ctx, "ca1f3e33b82e540bf8a92551c78ba411eb6944be", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to update(add) sub comment cache: id(%v), uuid(%s), userUuid(%s) ", id, uuid, userUuid))
	}
	return nil
}

func (r *commentRepo) SetCommentContentIrregular(ctx context.Context, review *biz.TextReview) (*biz.TextReview, error) {
	ar := &CommentContentReview{
		CommentId: review.CommentId,
		Comment:   review.Comment,
		Kind:      review.Kind,
		Uuid:      review.Uuid,
		JobId:     review.JobId,
		Label:     review.Label,
		Result:    review.Result,
		Section:   review.Section,
	}
	err := r.data.DB(ctx).Select("CommentId", "Comment", "Kind", "Uuid", "JobId", "Label", "Result", "Section").Create(ar).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to add comment content review record: review(%v)", review))
	}
	review.Id = int32(ar.ID)
	review.CreateAt = time.Now().Format("2006-01-02")
	return review, nil
}

func (r *commentRepo) SetCommentContentIrregularToCache(ctx context.Context, review *biz.TextReview) error {
	marshal, err := review.MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set comment content irregular to json: json.Marshal(%v)", review))
	}

	keys := []string{"comment_content_irregular_" + review.Uuid}
	values := []interface{}{marshal}
	_, err = r.data.redisCli.EvalSha(ctx, "8f6205011a2b264278a7c5bc0a2bcd1006ac6e5d", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set comment content irregular to cache: review(%v)", review))
	}
	return nil
}

func (r *commentRepo) CancelCommentAgree(ctx context.Context, id int32, uuid string) error {
	c := Comment{}
	err := r.data.DB(ctx).Model(&c).Where("comment_id = ? and uuid = ? and agree > 0", id, uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel comment agree: id(%v)", id))
	}
	return nil
}

func (r *commentRepo) CancelSubCommentAgree(ctx context.Context, id int32, uuid string) error {
	sc := SubComment{}
	err := r.data.DB(ctx).Model(&sc).Where("comment_id = ? and uuid = ? and agree > 0", id, uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel sub comment agree: id(%v)", id))
	}
	return nil
}

func (r *commentRepo) CancelUserCommentAgree(ctx context.Context, id int32, userUuid string) error {
	ca := CommentAgree{}
	err := r.data.DB(ctx).Model(&ca).Where("comment_id = ? and uuid = ?", id, userUuid).Update("status", 2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel a comment agree to db: comment_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *commentRepo) CancelCommentAgreeFromCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	hotKey := fmt.Sprintf("comment_%v_%v_hot", creationId, creationType)
	statisticKey := fmt.Sprintf("comment_%v", id)
	userKey := fmt.Sprintf("user_comment_agree_%s", userUuid)

	keys := []string{hotKey, statisticKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), id}
	_, err := r.data.redisCli.EvalSha(ctx, "9dc4d48a178101039febb85f383bdf1ce03084e8", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to update(cancel) comment cache: id(%v), creationId(%v), creationType(%v), uuid(%s), userUuid(%s) ", id, creationId, creationType, uuid, userUuid))
	}
	return nil
}

func (r *commentRepo) CancelSubCommentAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("comment_%v", id)
	userKey := fmt.Sprintf("user_comment_agree_%s", userUuid)

	keys := []string{statisticKey, userKey}
	values := []interface{}{id}
	_, err := r.data.redisCli.EvalSha(ctx, "1167c6dce1ed0d4e0a11e9fed90cae399f49f259", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to update(cancel) sub comment cache: id(%v), uuid(%s), userUuid(%s) ", id, uuid, userUuid))
	}
	return nil
}

func (r *commentRepo) SendReviewToMq(ctx context.Context, review *biz.CommentReview) error {
	data, err := review.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}

func (r *commentRepo) SendCommentToMq(ctx context.Context, comment *biz.Comment, mode string) error {
	commentMap := &biz.SendCommentMap{
		Uuid:         comment.Uuid,
		Id:           comment.CommentId,
		CreationId:   comment.CreationId,
		CreationType: comment.CreationType,
		Mode:         mode,
	}
	data, err := commentMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{comment.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send comment to mq: %v", comment))
	}
	return nil
}

func (r *commentRepo) SendCommentContentIrregularToMq(ctx context.Context, review *biz.TextReview) error {
	data, err := review.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send comment content review to mq: %v", err))
	}
	return nil
}

func (r *commentRepo) SendSubCommentToMq(ctx context.Context, comment *biz.SubComment, mode string) error {
	commentMap := &biz.SendSubCommentMap{
		Uuid:     comment.Uuid,
		Id:       comment.CommentId,
		RootId:   comment.RootId,
		ParentId: comment.ParentId,
		Mode:     mode,
	}
	data, err := commentMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{comment.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send comment to mq: %v", comment))
	}
	return nil
}

func (r *commentRepo) SendCommentAgreeToMq(ctx context.Context, id, creationId, creationType int32, uuid, userUuid, mode string) error {
	commentMap := &biz.SendCommentAgreeMap{
		Uuid:         uuid,
		Id:           id,
		CreationId:   creationId,
		CreationType: creationType,
		UserUuid:     userUuid,
		Mode:         mode,
	}
	data, err := commentMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{userUuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set comment agree to mq: id(%v), uuid(%s), userUuid(%s), mode(%s)", id, uuid, userUuid, mode))
	}
	return nil
}

func (r *commentRepo) SendSubCommentAgreeToMq(ctx context.Context, id int32, uuid, userUuid, mode string) error {
	commentMap := &biz.SendSubCommentAgreeMap{
		Uuid:     uuid,
		Id:       id,
		UserUuid: userUuid,
		Mode:     mode,
	}
	data, err := commentMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{userUuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set sub comment agree to mq: id(%v), uuid(%s), userUuid(%s), mode(%s)", id, uuid, userUuid, mode))
	}
	return nil
}

func (r *commentRepo) SendCommentStatisticToMq(ctx context.Context, uuid, userUuid, mode string) error {
	achievement := &biz.SendCommentStatisticMap{
		Uuid:     uuid,
		UserUuid: userUuid,
		Mode:     mode,
	}
	data, err := achievement.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk statistic to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *commentRepo) SendScoreToMq(ctx context.Context, score int32, uuid, mode string) error {
	scoreMap := &biz.SendScoreMap{
		Uuid:  uuid,
		Score: score,
		Mode:  mode,
	}
	data, err := scoreMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send score to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *commentRepo) CreateComment(ctx context.Context, id, creationId, creationType int32, creationAuthor, uuid string) error {
	comment := &Comment{
		CommentId:      id,
		Uuid:           uuid,
		CreationId:     creationId,
		CreationType:   creationType,
		CreationAuthor: creationAuthor,
	}
	err := r.data.DB(ctx).Select("CommentId", "Uuid", "CreationId", "CreationType", "CreationAuthor").Create(comment).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a comment: id(%v), uuid(%s), CreationAuthor(%s), creationId(%v), creationType(%v)", id, uuid, creationAuthor, creationId, creationType))
	}
	return nil
}

func (r *commentRepo) CreateSubComment(ctx context.Context, id, rootId, parentId, creationId, creationType int32, creationAuthor, rootUser, uuid, reply string) error {
	sc := &SubComment{}
	sc.CommentId = id
	sc.RootId = rootId
	sc.ParentId = parentId
	sc.CreationId = creationId
	sc.CreationType = creationType
	sc.CreationAuthor = creationAuthor
	sc.RootUser = rootUser
	sc.Uuid = uuid
	sc.Reply = reply
	err := r.data.DB(ctx).Model(&SubComment{}).Create(sc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a sub comment: id(%v), rootId(%v), parentId(%v),creationId(%v), creationType(%v), creationAuthor(%s), rootUser(%s), uuid(%s), reply(%s)", id, rootId, parentId, creationId, creationType, creationAuthor, rootUser, uuid, reply))
	}
	return nil
}

func (r *commentRepo) CreateCommentCache(ctx context.Context, id, creationId, creationType int32, creationAuthor, uuid string) error {
	ids := strconv.Itoa(int(id))
	creationIds := strconv.Itoa(int(creationId))
	creationTypes := strconv.Itoa(int(creationType))
	commentUserReplyField := ""
	commentUserRepliedField := ""
	userCommentCreationReplyList := ""
	userCommentCreationRepliedList := ""
	userCommentCreationMessageRepliedList := "user_comment_replied_list_" + creationAuthor
	switch creationType {
	case 1:
		commentUserReplyField = "article_reply"
		commentUserRepliedField = "article_replied"
		userCommentCreationReplyList = "user_comment_article_reply_list_" + uuid
		userCommentCreationRepliedList = "user_comment_article_replied_list_" + creationAuthor
	case 3:
		commentUserReplyField = "talk_reply"
		commentUserRepliedField = "talk_replied"
		userCommentCreationReplyList = "user_comment_talk_reply_list_" + uuid
		userCommentCreationRepliedList = "user_comment_talk_replied_list_" + creationAuthor
	}

	keys := []string{
		"comment_" + creationIds + "_" + creationTypes,
		"comment_" + creationIds + "_" + creationTypes + "_hot",
		"comment_" + ids,
		"comment_user_" + uuid,
		"comment_user_" + creationAuthor,
		commentUserReplyField,
		commentUserRepliedField,
		userCommentCreationReplyList,
		userCommentCreationRepliedList,
		userCommentCreationMessageRepliedList,
	}
	values := []interface{}{
		int32(time.Now().Unix()),
		ids + "%" + uuid,
		ids + "%" + creationIds + "%" + creationAuthor,
		ids + "%" + creationIds + "%" + uuid,
		ids + "%" + creationIds + "%" + creationTypes + "%" + uuid,
	}
	_, err := r.data.redisCli.EvalSha(ctx, "8651c46058c26d194cb9ac3c575b0ef3ca05dba5", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create comment cache: id(%v), uuid(%s), creationAuthor(%v), creationId(%v), creationType(%v)", id, uuid, creationAuthor, creationId, creationType))
	}
	return nil
}

func (r *commentRepo) CreateSubCommentCache(ctx context.Context, id, rootId, parentId, creationId, creationType int32, creationAuthor, rootUser, uuid, reply string) error {
	ids := strconv.Itoa(int(id))
	rootIds := strconv.Itoa(int(rootId))
	parentIds := strconv.Itoa(int(parentId))
	creationIds := strconv.Itoa(int(creationId))
	creationTypes := strconv.Itoa(int(creationType))
	commentUserReplyField := ""
	commentUserRepliedField := ""
	userSubCommentCreationReplyList := ""
	userSubCommentCreationRepliedListForRoot := ""
	userSubCommentCreationRepliedListForParent := ""
	userSubCommentCreationMessageRepliedListForRoot := "user_sub_comment_replied_list_" + rootUser
	userSubCommentCreationMessageRepliedListForParent := "user_sub_comment_replied_list_" + reply
	switch creationType {
	case 1:
		commentUserReplyField = "article_reply_sub"
		commentUserRepliedField = "article_replied_sub"
		userSubCommentCreationReplyList = "user_sub_comment_article_reply_list_" + uuid
		userSubCommentCreationRepliedListForRoot = "user_sub_comment_article_replied_list_" + rootUser
		userSubCommentCreationRepliedListForParent = "user_sub_comment_article_replied_list_" + reply
	case 3:
		commentUserReplyField = "talk_reply_sub"
		commentUserRepliedField = "talk_replied_sub"
		userSubCommentCreationReplyList = "user_sub_comment_talk_reply_list_" + uuid
		userSubCommentCreationRepliedListForRoot = "user_sub_comment_talk_replied_list_" + rootUser
		userSubCommentCreationRepliedListForParent = "user_sub_comment_talk_replied_list_" + reply
	}

	keys := []string{
		"comment_" + ids,
		"comment_" + rootIds,
		"sub_comment_" + rootIds,
		"comment_user_" + uuid,
		"comment_user_" + rootUser,
		"comment_user_" + reply,
		commentUserReplyField,
		commentUserRepliedField,
		userSubCommentCreationReplyList,
		userSubCommentCreationRepliedListForRoot,
		userSubCommentCreationRepliedListForParent,
		userSubCommentCreationMessageRepliedListForRoot,
		userSubCommentCreationMessageRepliedListForParent,
	}
	values := []interface{}{
		int32(time.Now().Unix()),
		ids + "%" + uuid + "%" + reply,
		ids + "%" + creationIds + "%" + rootIds + "%" + parentIds + "%" + creationAuthor + "%" + rootUser + "%" + reply,
		ids + "%" + creationIds + "%" + rootIds + "%" + parentIds + "%" + uuid + "%" + creationAuthor + "%" + rootUser + "%" + reply,
		ids + "%" + creationIds + "%" + creationTypes + "%" + rootIds + "%" + parentIds + "%" + uuid + "%" + creationAuthor + "%" + rootUser + "%" + reply,
		parentId,
		rootUser,
		reply,
	}
	_, err := r.data.redisCli.EvalSha(ctx, "6cf899977cbe0951aba44112f0a98130cef04e1e", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create sub comment cache: id(%v), rootId(%v), parentId(%v), creationId(%v), creationType(%v), creationAuthor(%s), rootUser(%s),uuid(%s), reply(%s)", id, rootId, parentId, creationId, creationType, creationAuthor, rootUser, uuid, reply))
	}
	return nil
}

func (r *commentRepo) AddCreationComment(ctx context.Context, createId, createType int32, uuid string) error {
	_, err := r.data.cc.AddCreationComment(ctx, &creationV1.AddCreationCommentReq{
		Uuid:         uuid,
		CreationId:   createId,
		CreationType: createType,
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation comment: createId(%v), createType(%v),uuid(%s)", createId, createType, uuid))
	}
	return nil
}

func (r *commentRepo) AddArticleCommentUser(ctx context.Context, creationAuthor, uuid string) error {
	cu1 := &CommentUser{
		Uuid:         uuid,
		Comment:      1,
		ArticleReply: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"comment": gorm.Expr("comment + ?", 1), "article_reply": gorm.Expr("article_reply + ?", 1)}),
	}).Create(cu1).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article comment user reply: uuid(%v)", uuid))
	}

	cu2 := &CommentUser{
		Uuid:           creationAuthor,
		ArticleReplied: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"article_replied": gorm.Expr("article_replied + ?", 1)}),
	}).Create(cu2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article comment user replied: uuid(%v)", creationAuthor))
	}
	return nil
}

func (r *commentRepo) ReduceArticleCommentUser(ctx context.Context, creationAuthor, uuid string) error {
	cu1 := &CommentUser{}
	err := r.data.DB(ctx).Model(cu1).Where("uuid = ? and comment > 0", uuid).Update("comment", gorm.Expr("comment - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article comment user reply: uuid(%v)", uuid))
	}

	err = r.data.DB(ctx).Model(cu1).Where("uuid = ? and article_reply > 0", uuid).Update("article_reply", gorm.Expr("article_reply - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article comment user reply: uuid(%v)", uuid))
	}

	err = r.data.DB(ctx).Model(cu1).Where("uuid = ? and article_replied > 0", creationAuthor).Update("article_replied", gorm.Expr("article_replied - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article comment user replied: uuid(%v)", uuid))
	}
	return nil
}

func (r *commentRepo) ReduceTalkCommentUser(ctx context.Context, creationAuthor, uuid string) error {
	cu1 := &CommentUser{}
	err := r.data.DB(ctx).Model(cu1).Where("uuid = ? and comment > 0", uuid).Update("comment", gorm.Expr("comment - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk comment user reply: uuid(%v)", uuid))
	}

	err = r.data.DB(ctx).Model(cu1).Where("uuid = ? and talk_reply > 0", uuid).Update("talk_reply", gorm.Expr("talk_reply - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk comment user reply: uuid(%v)", uuid))
	}

	err = r.data.DB(ctx).Model(cu1).Where("uuid = ? and talk_replied > 0", creationAuthor).Update("talk_replied", gorm.Expr("talk_replied - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk comment user replied: uuid(%v)", uuid))
	}
	return nil
}

func (r *commentRepo) ReduceCreationComment(ctx context.Context, createId, createType int32, uuid string) error {
	_, err := r.data.cc.ReduceCreationComment(ctx, &creationV1.ReduceCreationCommentReq{
		Uuid:         uuid,
		CreationId:   createId,
		CreationType: createType,
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation comment: createId(%v), createType(%v), uuid(%s)", createId, createType, uuid))
	}
	return nil
}

func (r *commentRepo) ReduceArticleSubCommentUser(ctx context.Context, parentId int32, rootUser, parentUser, uuid string) error {
	cu1 := &CommentUser{}
	err := r.data.DB(ctx).Model(cu1).Where("uuid = ? and comment > 0", uuid).Update("comment", gorm.Expr("comment - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article sub comment user comment: uuid(%v)", uuid))
	}

	cu2 := &CommentUser{}
	err = r.data.DB(ctx).Model(cu2).Where("uuid = ? and article_reply_sub > 0", uuid).Update("article_reply_sub", gorm.Expr("article_reply_sub - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article sun comment user reply: uuid(%v)", uuid))
	}

	cu3 := &CommentUser{}
	err = r.data.DB(ctx).Model(cu3).Where("uuid = ? and article_replied_sub > 0", rootUser).Update("article_replied_sub", gorm.Expr("article_replied_sub - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article sub comment user replied: uuid(%v)", rootUser))
	}

	if parentId == 0 || rootUser == parentUser {
		return nil
	}

	cu4 := &CommentUser{}
	err = r.data.DB(ctx).Model(cu4).Where("uuid = ? and article_replied_sub > 0", parentUser).Update("article_replied_sub", gorm.Expr("article_replied_sub - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article sub comment user replied: uuid(%v)", parentUser))
	}
	return nil
}

func (r *commentRepo) ReduceTalkSubCommentUser(ctx context.Context, parentId int32, rootUser, parentUser, uuid string) error {
	cu1 := &CommentUser{}
	err := r.data.DB(ctx).Model(cu1).Where("uuid = ? and comment > 0", uuid).Update("comment", gorm.Expr("comment - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk sub comment user comment: uuid(%v)", uuid))
	}

	cu2 := &CommentUser{}
	err = r.data.DB(ctx).Model(cu2).Where("uuid = ? and talk_reply_sub > 0", uuid).Update("talk_reply_sub", gorm.Expr("talk_reply_sub - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk sun comment user reply: uuid(%v)", uuid))
	}

	cu3 := &CommentUser{}
	err = r.data.DB(ctx).Model(cu3).Where("uuid = ? and talk_replied_sub > 0", rootUser).Update("talk_replied_sub", gorm.Expr("talk_replied_sub - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk sub comment user replied: uuid(%v)", rootUser))
	}

	if parentId == 0 || rootUser == parentUser {
		return nil
	}

	cu4 := &CommentUser{}
	err = r.data.DB(ctx).Model(cu4).Where("uuid = ? and talk_replied_sub > 0", parentUser).Update("talk_replied_sub", gorm.Expr("talk_replied_sub - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce talk sub comment user replied: uuid(%v)", parentUser))
	}
	return nil
}

func (r *commentRepo) AddArticleSubCommentUser(ctx context.Context, parentId int32, rootUser, parentUser, uuid string) error {
	ach1 := &CommentUser{
		Uuid:            uuid,
		Comment:         1,
		ArticleReplySub: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"comment": gorm.Expr("comment + ?", 1), "article_reply_sub": gorm.Expr("article_reply_sub + ?", 1)}),
	}).Create(ach1).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article sub comment user reply: uuid(%v)", uuid))
	}

	ach2 := &CommentUser{
		Uuid:              rootUser,
		ArticleRepliedSub: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"article_replied_sub": gorm.Expr("article_replied_sub + ?", 1)}),
	}).Create(ach2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article sub comment user replied: rootUser(%v)", rootUser))
	}

	if parentId == 0 || rootUser == parentUser {
		return nil
	}

	ach3 := &CommentUser{
		Uuid:              parentUser,
		ArticleRepliedSub: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"article_replied_sub": gorm.Expr("article_replied_sub + ?", 1)}),
	}).Create(ach3).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article sub comment user replied: parentUser(%v)", parentUser))
	}
	return nil
}

func (r *commentRepo) AddTalkCommentUser(ctx context.Context, creationAuthor, uuid string) error {
	ach1 := &CommentUser{
		Uuid:      uuid,
		Comment:   1,
		TalkReply: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"comment": gorm.Expr("comment + ?", 1), "talk_reply": gorm.Expr("talk_reply + ?", 1)}),
	}).Create(ach1).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk comment user reply: uuid(%v)", uuid))
	}

	ach2 := &CommentUser{
		Uuid:        creationAuthor,
		TalkReplied: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"talk_replied": gorm.Expr("talk_replied + ?", 1)}),
	}).Create(ach2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk comment user replied: uuid(%v)", creationAuthor))
	}
	return nil
}

func (r *commentRepo) AddTalkSubCommentUser(ctx context.Context, parentId int32, rootUser, parentUser, uuid string) error {
	ach1 := &CommentUser{
		Uuid:         uuid,
		Comment:      1,
		TalkReplySub: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"comment": gorm.Expr("comment + ?", 1), "talk_reply_sub": gorm.Expr("talk_reply_sub + ?", 1)}),
	}).Create(ach1).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk sub comment user reply: uuid(%v)", uuid))
	}

	ach2 := &CommentUser{
		Uuid:           rootUser,
		TalkRepliedSub: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"talk_replied_sub": gorm.Expr("talk_replied_sub + ?", 1)}),
	}).Create(ach2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk sub comment user replied: rootUser(%v)", rootUser))
	}

	if parentId == 0 || rootUser == parentUser {
		return nil
	}

	ach3 := &CommentUser{
		Uuid:           parentUser,
		TalkRepliedSub: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"talk_replied_sub": gorm.Expr("talk_replied_sub + ?", 1)}),
	}).Create(ach3).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add talk sub comment user replied: parentUser(%v)", parentUser))
	}
	return nil
}

func (r *commentRepo) subCommentCacheExist(ctx context.Context, ids, rootIds string) (int32, int32, error) {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, "sub_comment_"+rootIds)
		pipe.Exists(ctx, "comment_"+rootIds)
		return nil
	})
	if err != nil {
		return 0, 0, errors.Wrapf(err, fmt.Sprintf("fail to check if sub comment exist from cache: commentId(%s),rootId(%s)", ids, rootIds))
	}

	exists := make([]int32, 0, len(cmd))
	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, int32(exist))
	}
	return exists[0], exists[1], nil
}

func (r *commentRepo) DeleteCommentDraft(ctx context.Context, id int32, uuid string) error {
	cd := &CommentDraft{}
	cd.ID = uint(id)
	err := r.data.DB(ctx).Where("uuid = ?", uuid).Delete(cd).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete comment draft: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *commentRepo) RemoveComment(ctx context.Context, id int32, uuid string) error {
	c := &Comment{}
	err := r.data.DB(ctx).Where("comment_id = ? and (uuid = ? or creation_author = ?)", id, uuid, uuid).Delete(c).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete a comment: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *commentRepo) RemoveSubComment(ctx context.Context, id, rootId int32, uuid string) error {
	c := &Comment{}
	err := r.data.DB(ctx).Model(c).Where("comment_id = ? and comment > 0", rootId).Update("comment", gorm.Expr("comment - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel comment comment: id(%v)", rootId))
	}

	sc := &SubComment{}
	err = r.data.DB(ctx).Where("comment_id = ? and (uuid = ? or root_user = ?)", id, uuid, uuid).Delete(sc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete a comment: id(%v), uuid(%s)", id, uuid))
	}

	return nil
}

func (r *commentRepo) RemoveCommentAgree(ctx context.Context, id int32, uuid string) error {
	ca := &CommentAgree{}
	err := r.data.DB(ctx).Where("comment_id = ? and uuid = ?", id, uuid).Delete(ca).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete a comment agree record: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *commentRepo) RemoveCommentCache(ctx context.Context, id, creationId, creationType int32, creationAuthor, uuid string) error {
	ids := strconv.Itoa(int(id))
	creationIds := strconv.Itoa(int(creationId))
	creationTypes := strconv.Itoa(int(creationType))
	commentUserReplyField := ""
	commentUserRepliedField := ""
	userCommentCreationReplyList := ""
	userCommentCreationRepliedList := ""
	userCommentCreationMessageRepliedList := "user_comment_replied_list_" + creationAuthor
	switch creationType {
	case 1:
		commentUserReplyField = "article_reply"
		commentUserRepliedField = "article_replied"
		userCommentCreationReplyList = "user_comment_article_reply_list_" + uuid
		userCommentCreationRepliedList = "user_comment_article_replied_list_" + creationAuthor
	case 3:
		commentUserReplyField = "talk_reply"
		commentUserRepliedField = "talk_replied"
		userCommentCreationReplyList = "user_comment_talk_reply_list_" + uuid
		userCommentCreationRepliedList = "user_comment_talk_replied_list_" + creationAuthor
	}

	keys := []string{
		"comment_" + creationIds + "_" + creationTypes,
		"comment_" + creationIds + "_" + creationTypes + "_hot",
		"comment_" + ids,
		"comment_user_" + uuid,
		"comment_user_" + creationAuthor,
		commentUserReplyField,
		commentUserRepliedField,
		userCommentCreationReplyList,
		userCommentCreationRepliedList,
		userCommentCreationMessageRepliedList,
	}
	values := []interface{}{
		id,
		ids + "%" + uuid,
		ids + "%" + creationIds + "%" + creationAuthor,
		ids + "%" + creationIds + "%" + uuid,
		ids + "%" + creationIds + "%" + creationTypes + "%" + uuid,
	}

	_, err := r.data.redisCli.EvalSha(ctx, "ed14acc35b617d746219bc729f8930557bc48fa5", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to remove comment cache: id(%v), uuid(%s), creationId(%v), creationType(%v), creationAuthor(%v)", id, uuid, creationId, creationType, creationAuthor))
	}
	return nil
}

func (r *commentRepo) RemoveSubCommentCache(ctx context.Context, id, rootId, parentId, creationId, creationType int32, creationAuthor, rootUser, uuid, reply string) error {
	ids := strconv.Itoa(int(id))
	rootIds := strconv.Itoa(int(rootId))
	parentIds := strconv.Itoa(int(parentId))
	creationIds := strconv.Itoa(int(creationId))
	creationTypes := strconv.Itoa(int(creationType))
	commentUserReplyField := ""
	commentUserRepliedField := ""
	userSubCommentCreationReplyList := ""
	userSubCommentCreationRepliedListForRoot := ""
	userSubCommentCreationRepliedListForParent := ""
	userSubCommentCreationMessageRepliedListForRoot := "user_sub_comment_replied_list_" + rootUser
	userSubCommentCreationMessageRepliedListForParent := "user_sub_comment_replied_list_" + reply
	switch creationType {
	case 1:
		commentUserReplyField = "article_reply_sub"
		commentUserRepliedField = "article_replied_sub"
		userSubCommentCreationReplyList = "user_sub_comment_article_reply_list_" + uuid
		userSubCommentCreationRepliedListForRoot = "user_sub_comment_article_replied_list_" + rootUser
		userSubCommentCreationRepliedListForParent = "user_sub_comment_article_replied_list_" + reply
	case 3:
		commentUserReplyField = "talk_reply_sub"
		commentUserRepliedField = "talk_replied_sub"
		userSubCommentCreationReplyList = "user_sub_comment_talk_reply_list_" + uuid
		userSubCommentCreationRepliedListForRoot = "user_sub_comment_talk_replied_list_" + rootUser
		userSubCommentCreationRepliedListForParent = "user_sub_comment_talk_replied_list_" + reply
	}

	keys := []string{
		"comment_" + ids,
		"comment_" + rootIds,
		"sub_comment_" + rootIds,
		"comment_user_" + uuid,
		"comment_user_" + rootUser,
		"comment_user_" + reply,
		commentUserReplyField,
		commentUserRepliedField,
		userSubCommentCreationReplyList,
		userSubCommentCreationRepliedListForRoot,
		userSubCommentCreationRepliedListForParent,
		userSubCommentCreationMessageRepliedListForRoot,
		userSubCommentCreationMessageRepliedListForParent,
	}
	values := []interface{}{
		id,
		ids + "%" + uuid + "%" + reply,
		ids + "%" + creationIds + "%" + rootIds + "%" + parentIds + "%" + creationAuthor + "%" + rootUser + "%" + reply,
		ids + "%" + creationIds + "%" + rootIds + "%" + parentIds + "%" + uuid + "%" + creationAuthor + "%" + rootUser + "%" + reply,
		ids + "%" + creationIds + "%" + creationTypes + "%" + rootIds + "%" + parentIds + "%" + uuid + "%" + creationAuthor + "%" + rootUser + "%" + reply,
		parentId,
		rootUser,
		reply,
	}
	_, err := r.data.redisCli.EvalSha(ctx, "990385b0f935e627854bcf10a4fdf6a0dce74d47", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to remove sub comment cache: id(%v), rootId(%v), parentId(%v), creationId(%v), creationType(%v), creationAuthor(%s), rootUser(%s),uuid(%s), reply(%s)", id, rootId, parentId, creationId, creationType, creationAuthor, rootUser, uuid, reply))
	}
	return nil
}
