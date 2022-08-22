package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
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
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("comment draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%s)", uuid))
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
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user comment agree from db: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range list {
		agreeMap[item.CommentId] = true
	}
	if len(list) != 0 {
		go r.setUserCommentAgreeToCache(uuid, list)
	}
	return agreeMap, nil
}

func (r *commentRepo) setUserCommentAgreeToCache(uuid string, agreeList []*CommentAgree) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0)
		key := "user_comment_agree_" + uuid
		for _, item := range agreeList {
			set = append(set, item.CommentId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Hour*8)
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

	go r.setCommentUserToCache(key, commentUser)

	return commentUser, nil
}

func (r *commentRepo) getCommentUserFromCache(ctx context.Context, key string) (*biz.CommentUser, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "comment").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment user form cache: key(%s)", key))
	}
	val := []int32{0}
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
		Comment: val[0],
	}, nil
}

func (r *commentRepo) getCommentUserFromDB(ctx context.Context, uuid string) (*biz.CommentUser, error) {
	cu := &CommentUser{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(cu).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get comment user from db: uuid(%v)", uuid))
	}
	return &biz.CommentUser{
		Comment: cu.Comment,
	}, nil
}

func (r *commentRepo) setCommentUserToCache(key string, commentUser *biz.CommentUser) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(context.Background(), key, "comment", commentUser.Comment)
		pipe.Expire(ctx, key, time.Hour*8)
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
		go r.setCommentToCache(creationId, creationType, comment)
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
		go r.setSubCommentToCache(id, subComment)
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

	comment := make([]*biz.Comment, 0)
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

	subComment := make([]*biz.SubComment, 0)
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
	err := r.data.db.WithContext(ctx).Where("creation_id = ? and creation_type = ?", creationId, creationType).Order("comment_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0)
	for _, item := range list {
		comment = append(comment, &biz.Comment{
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
	err := r.data.db.WithContext(ctx).Where("root_id = ?", id).Order("comment_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get sub comment from db: page(%v)", page))
	}

	subComment := make([]*biz.SubComment, 0)
	for _, item := range list {
		subComment = append(subComment, &biz.SubComment{
			CommentId: item.CommentId,
			Uuid:      item.Uuid,
			Reply:     item.Reply,
		})
	}
	return subComment, nil
}

func (r *commentRepo) setCommentToCache(creationId, creationType int32, comment []*biz.Comment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		key := fmt.Sprintf("comment_%v_%v", creationId, creationType)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.CommentId),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set comment to cache: comment(%v)", comment)
	}
}

func (r *commentRepo) setSubCommentToCache(id int32, subComment []*biz.SubComment) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		key := fmt.Sprintf("sub_comment_%v", id)
		for _, item := range subComment {
			z = append(z, &redis.Z{
				Score:  float64(item.CommentId),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + item.Uuid + "%" + item.Reply,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set sub comment to cache: comment(%v)", subComment)
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
		go r.setCommentHotToCache(creationId, creationType, comment)
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

	comment := make([]*biz.Comment, 0)
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
	err := r.data.db.WithContext(ctx).Where("creation_id = ? and creation_type = ?", creationId, creationType).Order("agree desc, comment_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get comment hot from db: page(%v)", page))
	}

	comment := make([]*biz.Comment, 0)
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
		z := make([]*redis.Z, 0)
		key := fmt.Sprintf("comment_%v_%v_hot", creationId, creationType)
		for _, item := range comment {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.CommentId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		pipe.Expire(context.Background(), key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set comment hot to cache: comment(%v)", comment)
	}
}

func (r *commentRepo) GetCommentListStatistic(ctx context.Context, ids []int32) ([]*biz.CommentStatistic, error) {
	commentListStatistic := make([]*biz.CommentStatistic, 0)
	exists, unExists, err := r.commentStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

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
	commentListStatistic := make([]*biz.CommentStatistic, 0)
	exists, unExists, err := r.commentStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

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

func (r *commentRepo) GetRootCommentUserId(ctx context.Context, id int32) (string, error) {
	sc := &Comment{}
	err := r.data.DB(ctx).Where("comment_id = ?", id).First(sc).Error
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("db query system error: rootId(%v)", id))
	}
	return sc.Uuid, nil
}

func (r *commentRepo) GetParentCommentUserId(ctx context.Context, id int32) (string, error) {
	sc := &SubComment{}
	err := r.data.DB(ctx).Where("comment_id = ?", id).First(sc).Error
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("db query system error: parentId(%v)", id))
	}
	return sc.Uuid, nil
}

func (r *commentRepo) GetSubCommentReply(ctx context.Context, id int32, uuid string) (string, error) {
	sc := &SubComment{}
	err := r.data.DB(ctx).Where("comment_id = ? and uuid = ?", id, uuid).First(sc).Error
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("db query system error: parentId(%v)", id))
	}
	return sc.Reply, nil
}

func (r *commentRepo) commentStatisticExist(ctx context.Context, ids []int32) ([]int32, []int32, error) {
	exists := make([]int32, 0)
	unExists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range ids {
			pipe.Exists(ctx, "comment_"+strconv.Itoa(int(item)))
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if comment statistic exist from cache: ids(%v)", ids))
	}

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
	list := make([]*Comment, 0)
	err := r.data.db.WithContext(ctx).Where("comment_id IN ?", unExists).Find(&list).Error
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
		go r.setCommentStatisticToCache(list)
	}

	return nil
}

func (r *commentRepo) getSubCommentStatisticFromDb(ctx context.Context, unExists []int32, commentListStatistic *[]*biz.CommentStatistic) error {
	list := make([]*SubComment, 0)
	err := r.data.db.WithContext(ctx).Where("comment_id IN ?", unExists).Find(&list).Error
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
		go r.setSubCommentStatisticToCache(list)
	}

	return nil
}

func (r *commentRepo) setCommentStatisticToCache(commentList []*Comment) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, item := range commentList {
			key := "comment_" + strconv.Itoa(int(item.CommentId))
			pipe.HSetNX(context.Background(), key, "agree", item.Agree)
			pipe.HSetNX(context.Background(), key, "comment", item.Comment)
			pipe.Expire(context.Background(), key, time.Hour*8)
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
			pipe.Expire(context.Background(), key, time.Hour*8)
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
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						local result = redis.call("SADD", key, change)
						return result
					end
					return 0
	`)
	keys := []string{"user_comment_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	exists := make([]int64, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, hotKey)
		pipe.Exists(ctx, statisticKey)
		pipe.Exists(ctx, userKey)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to check if cache about comment exist: id(%v), creationId(%v), creationType(%v), uuid(%s), userUuid(%s) ", id, creationId, creationType, uuid, userUuid))
	}

	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, exist)
	}

	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if exists[0] == 1 {
			pipe.ZIncrBy(ctx, hotKey, 1, fmt.Sprintf("%v%s%s", id, "%", uuid))
		}
		if exists[1] == 1 {
			pipe.HIncrBy(ctx, statisticKey, "agree", 1)
		}
		if exists[2] == 1 {
			pipe.SAdd(ctx, userKey, id)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to update(add) comment cache: id(%v), creationId(%v), creationType(%v), uuid(%s), userUuid(%s) ", id, creationId, creationType, uuid, userUuid))
	}
	return nil
}

func (r *commentRepo) SetSubCommentAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("comment_%v", id)
	userKey := fmt.Sprintf("user_comment_agree_%s", userUuid)
	exists := make([]int64, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, statisticKey)
		pipe.Exists(ctx, userKey)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to check if cache about sub comment exist: id(%v), uuid(%s), userUuid(%s) ", id, uuid, userUuid))
	}

	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, exist)
	}

	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if exists[0] == 1 {
			pipe.HIncrBy(ctx, statisticKey, "agree", 1)
		}
		if exists[1] == 1 {
			pipe.SAdd(ctx, userKey, id)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to update(add) sub comment cache: id(%v), uuid(%s), userUuid(%s) ", id, uuid, userUuid))
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

	var script = redis.NewScript(`
					local hotKey = KEYS[1]
                    local member = ARGV[1]
					local hotKeyExist = redis.call("EXISTS", hotKey)
					if hotKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", hotKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", hotKey, -1, member)
						end
					end

					local statisticKey = KEYS[2]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[3]
					local commentId = ARGV[2]
					redis.call("SREM", userKey, commentId)
					return 0
	`)
	keys := []string{hotKey, statisticKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), id}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to update(cancel) comment cache: id(%v), creationId(%v), creationType(%v), uuid(%s), userUuid(%s) ", id, creationId, creationType, uuid, userUuid))
	}
	return nil
}

func (r *commentRepo) CancelSubCommentAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("comment_%v", id)
	userKey := fmt.Sprintf("user_comment_agree_%s", userUuid)

	var script = redis.NewScript(`
					local statisticKey = KEYS[1]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[2]
					local commentId = ARGV[1]
					redis.call("SREM", userKey, commentId)
					return 0
	`)
	keys := []string{statisticKey, userKey}
	values := []interface{}{id}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to update(cancel) sub comment cache: id(%v), uuid(%s), userUuid(%s) ", id, uuid, userUuid))
	}
	return nil
}

func (r *commentRepo) SendReviewToMq(ctx context.Context, review *biz.CommentReview) error {
	reviewMap := map[string]interface{}{}
	reviewMap["uuid"] = review.Uuid
	reviewMap["id"] = review.Id
	reviewMap["mode"] = review.Mode
	data, err := json.Marshal(reviewMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "review",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.reviewMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}

func (r *commentRepo) SendCommentToMq(ctx context.Context, comment *biz.Comment, mode string) error {
	commentMap := map[string]interface{}{}
	commentMap["uuid"] = comment.Uuid
	commentMap["id"] = comment.CommentId
	commentMap["creationId"] = comment.CreationId
	commentMap["creationType"] = comment.CreationType
	commentMap["mode"] = mode

	data, err := json.Marshal(commentMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "comment",
		Body:  data,
	}
	msg.WithKeys([]string{comment.Uuid})
	_, err = r.data.commonMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send comment to mq: %v", comment))
	}
	return nil
}

func (r *commentRepo) SendSubCommentToMq(ctx context.Context, comment *biz.SubComment, mode string) error {
	commentMap := map[string]interface{}{}
	commentMap["uuid"] = comment.Uuid
	commentMap["id"] = comment.CommentId
	commentMap["rootId"] = comment.RootId
	commentMap["parentId"] = comment.ParentId
	commentMap["mode"] = mode

	data, err := json.Marshal(commentMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "comment",
		Body:  data,
	}
	msg.WithKeys([]string{comment.Uuid})
	_, err = r.data.commonMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send comment to mq: %v", comment))
	}
	return nil
}

func (r *commentRepo) SendCommentAgreeToMq(ctx context.Context, id, creationId, creationType int32, uuid, userUuid, mode string) error {
	commentMap := map[string]interface{}{}
	commentMap["uuid"] = uuid
	commentMap["id"] = id
	commentMap["creationId"] = creationId
	commentMap["creationType"] = creationType
	commentMap["userUuid"] = userUuid
	commentMap["mode"] = mode

	data, err := json.Marshal(commentMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "comment",
		Body:  data,
	}
	msg.WithKeys([]string{userUuid})
	_, err = r.data.commonMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set comment agree to mq: id(%v), uuid(%s), userUuid(%s), mode(%s)", id, uuid, userUuid, mode))
	}
	return nil
}

func (r *commentRepo) SendSubCommentAgreeToMq(ctx context.Context, id int32, uuid, userUuid, mode string) error {
	commentMap := map[string]interface{}{}
	commentMap["uuid"] = uuid
	commentMap["id"] = id
	commentMap["userUuid"] = userUuid
	commentMap["mode"] = mode

	data, err := json.Marshal(commentMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "comment",
		Body:  data,
	}
	msg.WithKeys([]string{userUuid})
	_, err = r.data.commonMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set sub comment agree to mq: id(%v), uuid(%s), userUuid(%s), mode(%s)", id, uuid, userUuid, mode))
	}
	return nil
}

func (r *commentRepo) SendCommentStatisticToMq(ctx context.Context, uuid, userUuid, mode string) error {
	achievement := map[string]interface{}{}
	achievement["uuid"] = uuid
	achievement["userUuid"] = userUuid
	achievement["mode"] = mode

	data, err := json.Marshal(achievement)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "achievement",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.achievementMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk statistic to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *commentRepo) SendScoreToMq(ctx context.Context, score int32, uuid, mode string) error {
	scoreMap := map[string]interface{}{}
	scoreMap["uuid"] = uuid
	scoreMap["score"] = score
	scoreMap["mode"] = mode

	data, err := json.Marshal(scoreMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "achievement",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.achievementMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send score to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *commentRepo) CreateComment(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	comment := &Comment{
		CommentId:    id,
		Uuid:         uuid,
		CreationId:   creationId,
		CreationType: creationType,
	}
	err := r.data.DB(ctx).Select("CommentId", "Uuid", "CreationId", "CreationType").Create(comment).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a comment: id(%v), uuid(%s), creationId(%v), creationType(%v)", id, uuid, creationId, creationType))
	}
	return nil
}

func (r *commentRepo) CreateSubComment(ctx context.Context, id, rootId, parentId int32, uuid, reply string) error {
	sc := &SubComment{}
	sc.CommentId = id
	sc.RootId = rootId
	sc.ParentId = parentId
	sc.Uuid = uuid
	sc.Reply = reply
	err := r.data.DB(ctx).Model(&SubComment{}).Create(sc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a sub comment: id(%v), rootId(%v), parentId(%v),uuid(%s), reply(%s)", id, rootId, parentId, uuid, reply))
	}
	return nil
}

func (r *commentRepo) CreateCommentCache(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	creationIds := strconv.Itoa(int(creationId))
	creationTypes := strconv.Itoa(int(creationType))

	newExist, hotExist, err := r.commentCacheExist(ctx, creationIds, creationTypes)
	if err != nil {
		return err
	}

	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, "comment_"+ids, "agree", 0)
		pipe.HSetNX(ctx, "comment_"+ids, "comment", 0)
		pipe.Expire(ctx, "comment_"+ids, time.Hour*8)

		if newExist == 1 {
			pipe.ZAddNX(ctx, "comment_"+creationIds+"_"+creationTypes, &redis.Z{
				Score:  float64(id),
				Member: ids + "%" + uuid,
			})
			pipe.Expire(ctx, "comment_"+creationIds+"_"+creationTypes, time.Hour*8)
		}

		if hotExist == 1 {
			pipe.ZAddNX(ctx, "comment_"+creationIds+"_"+creationTypes+"_hot", &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid,
			})
			pipe.Expire(ctx, "comment_"+creationIds+"_"+creationTypes+"_hot", time.Hour*8)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create comment cache: id(%v), uuid(%s), creationId(%v), creationType(%v)", id, uuid, creationId, creationType))
	}
	return nil
}

func (r *commentRepo) commentCacheExist(ctx context.Context, creationIds, creationTypes string) (int32, int32, error) {
	exists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, "comment_"+creationIds+"_"+creationTypes)
		pipe.Exists(ctx, "comment_"+creationIds+"_"+creationTypes+"_hot")
		return nil
	})
	if err != nil {
		return 0, 0, errors.Wrapf(err, fmt.Sprintf("fail to check if comment exist from cache: creationIds(%s),creationTypes(%s)", creationIds, creationTypes))
	}

	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, int32(exist))
	}
	return exists[0], exists[1], nil
}

func (r *commentRepo) CreateSubCommentCache(ctx context.Context, id, rootId int32, uuid, reply string) error {
	ids := strconv.Itoa(int(id))
	rootIds := strconv.Itoa(int(rootId))

	subCommentExist, commentExist, err := r.subCommentCacheExist(ctx, ids, rootIds)
	if err != nil {
		return err
	}

	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, "comment_"+ids, "agree", 0)
		pipe.Expire(ctx, "comment_"+ids, time.Hour*8)

		if subCommentExist == 1 {
			pipe.ZAddNX(ctx, "sub_comment_"+rootIds, &redis.Z{
				Score:  float64(id),
				Member: ids + "%" + uuid + "%" + reply,
			})
			pipe.Expire(ctx, "sub_comment_"+rootIds, time.Hour*8)
		}

		if commentExist == 1 {
			pipe.HIncrBy(ctx, "comment_"+rootIds, "comment", 1)
			pipe.Expire(ctx, "comment_"+rootIds, time.Hour*8)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create sub comment cache: id(%v), rootId(%v), uuid(%s), reply(%s)", id, rootId, uuid, reply))
	}
	return nil
}

func (r *commentRepo) subCommentCacheExist(ctx context.Context, ids, rootIds string) (int32, int32, error) {
	exists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, "sub_comment_"+rootIds)
		pipe.Exists(ctx, "comment_"+rootIds)
		return nil
	})
	if err != nil {
		return 0, 0, errors.Wrapf(err, fmt.Sprintf("fail to check if sub comment exist from cache: commentId(%s),rootId(%s)", ids, rootIds))
	}

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
	err := r.data.DB(ctx).Where("comment_id = ? and uuid = ?", id, uuid).Delete(c).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete a comment: id(%v), uuid(%s)", id, uuid))
	}

	sc := &SubComment{}
	err = r.data.DB(ctx).Where("root_id = ?", id).Delete(sc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete subcomment: id(%v), uuid(%s)", id, uuid))
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
	err = r.data.DB(ctx).Where("comment_id = ? and uuid = ?", id, uuid).Delete(sc).Error
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

func (r *commentRepo) RemoveCommentCache(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	creationIds := strconv.Itoa(int(creationId))
	creationTypes := strconv.Itoa(int(creationType))

	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, "comment_"+ids)
		pipe.ZRem(ctx, "comment_"+creationIds+"_"+creationTypes, ids+"%"+uuid)
		pipe.ZRem(ctx, "comment_"+creationIds+"_"+creationTypes+"_hot", ids+"%"+uuid)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to remove comment cache: id(%v), uuid(%s), creationId(%v), creationType(%v)", id, uuid, creationId, creationType))
	}
	return nil
}

func (r *commentRepo) RemoveSubCommentCache(ctx context.Context, id, rootId int32, uuid, reply, mode string) error {
	ids := strconv.Itoa(int(id))
	rootIds := strconv.Itoa(int(rootId))

	commetkey := "comment_" + ids
	subCommentKey := "sub_comment_" + rootIds
	rootKey := "comment_" + rootIds
	var script = redis.NewScript(`
					local commetkey = KEYS[1]
					redis.call("DEL", commetkey)

					local subCommentKey = KEYS[2]
                    local member = ARGV[1]
					redis.call("ZREM", subCommentKey, member)

					local rootKey = KEYS[3]
					local mode = ARGV[2]
					local value = redis.call("EXISTS", rootKey)
					if (mode == "final" and value == 1) then
						local number = tonumber(redis.call("HGET", rootKey, "comment"))
						if number > 0 then
  							redis.call("HINCRBY", rootKey, "comment", -1)
						end
					end
					return 0
	`)
	keys := []string{commetkey, subCommentKey, rootKey, mode}
	values := []interface{}{ids + "%" + uuid + "%" + reply, mode}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to remove comment cache: id(%v), uuid(%s), rootId(%v), reply(%s)", id, uuid, rootId, reply))
	}
	return nil
}
