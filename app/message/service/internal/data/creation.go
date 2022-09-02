package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	creationV1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"strconv"
)

type creationRepo struct {
	data *Data
	log  *log.Helper
}

func NewCreationRepo(data *Data, logger log.Logger) biz.CreationRepo {
	return &creationRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "message/data/creation")),
	}
}

func (r *creationRepo) ToReviewCreateArticle(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "article/" + uuid + "/" + strconv.Itoa(int(id)) + "/content",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCreationCli.callback["article_create"],
		},
	}

	_, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article create review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) ToReviewEditArticle(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "article/" + uuid + "/" + strconv.Itoa(int(id)) + "/content-edit",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCreationCli.callback["article_edit"],
		},
	}

	_, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article edit review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) ArticleCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateArticle(ctx, &creationV1.CreateArticleReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) ArticleEditReviewPass(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.EditArticle(ctx, &creationV1.EditArticleReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CreateArticleDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateArticleDbCacheAndSearch(ctx, &creationV1.CreateArticleDbCacheAndSearchReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) EditArticleCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.EditArticleCosAndSearch(ctx, &creationV1.EditArticleCosAndSearchReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) DeleteArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.DeleteArticleCacheAndSearch(ctx, &creationV1.DeleteArticleCacheAndSearchReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetArticleViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.SetArticleViewDbAndCache(ctx, &creationV1.SetArticleViewDbAndCacheReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetArticleAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetArticleAgreeDbAndCache(ctx, &creationV1.SetArticleAgreeDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetArticleCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetArticleCollectDbAndCache(ctx, &creationV1.SetArticleCollectDbAndCacheReq{
		Id:            id,
		CollectionsId: collectionsId,
		Uuid:          uuid,
		UserUuid:      userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CancelArticleAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelArticleAgreeDbAndCache(ctx, &creationV1.CancelArticleAgreeDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CancelArticleCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelArticleCollectDbAndCache(ctx, &creationV1.CancelArticleCollectDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) ToReviewCreateTalk(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "talk/" + uuid + "/" + strconv.Itoa(int(id)) + "/content",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCreationCli.callback["talk_create"],
		},
	}

	_, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk create review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) ToReviewEditTalk(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "talk/" + uuid + "/" + strconv.Itoa(int(id)) + "/content-edit",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCreationCli.callback["talk_edit"],
		},
	}

	_, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk edit review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) TalkCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateTalk(ctx, &creationV1.CreateTalkReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) TalkEditReviewPass(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.EditTalk(ctx, &creationV1.EditTalkReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CreateTalkDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateTalkDbCacheAndSearch(ctx, &creationV1.CreateTalkDbCacheAndSearchReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) EditTalkCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.EditTalkCosAndSearch(ctx, &creationV1.EditTalkCosAndSearchReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetTalkViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.SetTalkViewDbAndCache(ctx, &creationV1.SetTalkViewDbAndCacheReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetTalkAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetTalkAgreeDbAndCache(ctx, &creationV1.SetTalkAgreeDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CancelTalkAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelTalkAgreeDbAndCache(ctx, &creationV1.CancelTalkAgreeDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetTalkCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetTalkCollectDbAndCache(ctx, &creationV1.SetTalkCollectDbAndCacheReq{
		Id:            id,
		CollectionsId: collectionsId,
		Uuid:          uuid,
		UserUuid:      userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CancelTalkCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelTalkCollectDbAndCache(ctx, &creationV1.CancelTalkCollectDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) DeleteTalkCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.DeleteTalkCacheAndSearch(ctx, &creationV1.DeleteTalkCacheAndSearchReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) ToReviewCreateColumn(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "column/" + uuid + "/" + strconv.Itoa(int(id)) + "/content",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCreationCli.callback["column_create"],
		},
	}

	_, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send column create review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) ToReviewEditColumn(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "column/" + uuid + "/" + strconv.Itoa(int(id)) + "/content-edit",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCreationCli.callback["column_edit"],
		},
	}

	_, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send column edit review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) ColumnCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateColumn(ctx, &creationV1.CreateColumnReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) ColumnEditReviewPass(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.EditColumn(ctx, &creationV1.EditColumnReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CreateColumnDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateColumnDbCacheAndSearch(ctx, &creationV1.CreateColumnDbCacheAndSearchReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) EditColumnCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.EditColumnCosAndSearch(ctx, &creationV1.EditColumnCosAndSearchReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) DeleteColumnCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.DeleteColumnCacheAndSearch(ctx, &creationV1.DeleteColumnCacheAndSearchReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetColumnAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetColumnAgreeDbAndCache(ctx, &creationV1.SetColumnAgreeDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetColumnViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.SetColumnViewDbAndCache(ctx, &creationV1.SetColumnViewDbAndCacheReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CancelColumnAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelColumnAgreeDbAndCache(ctx, &creationV1.CancelColumnAgreeDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetColumnCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	_, err := r.data.cc.SetColumnCollectDbAndCache(ctx, &creationV1.SetColumnCollectDbAndCacheReq{
		Id:            id,
		CollectionsId: collectionsId,
		Uuid:          uuid,
		UserUuid:      userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CancelColumnCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.cc.CancelColumnCollectDbAndCache(ctx, &creationV1.CancelColumnCollectDbAndCacheReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) AddColumnIncludesDbAndCache(ctx context.Context, id, articleId int32, uuid string) error {
	_, err := r.data.cc.AddColumnIncludesDbAndCache(ctx, &creationV1.AddColumnIncludesDbAndCacheReq{
		Id:        id,
		ArticleId: articleId,
		Uuid:      uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) DeleteColumnIncludesDbAndCache(ctx context.Context, id, articleId int32, uuid string) error {
	_, err := r.data.cc.DeleteColumnIncludesDbAndCache(ctx, &creationV1.DeleteColumnIncludesDbAndCacheReq{
		Id:        id,
		ArticleId: articleId,
		Uuid:      uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) SetColumnSubscribeDbAndCache(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.SetColumnSubscribeDbAndCache(ctx, &creationV1.SetColumnSubscribeDbAndCacheReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CancelColumnSubscribeDbAndCache(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.CancelColumnSubscribeDbAndCache(ctx, &creationV1.CancelColumnSubscribeDbAndCacheReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) ToReviewCreateCollections(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "collections/" + uuid + "/" + strconv.Itoa(int(id)) + "/content",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCreationCli.callback["collections_create"],
		},
	}

	_, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send collections create review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) ToReviewEditCollections(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "collections/" + uuid + "/" + strconv.Itoa(int(id)) + "/content-edit",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCreationCli.callback["collections_edit"],
		},
	}

	_, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send collections edit review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) CollectionsCreateReviewPass(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateCollections(ctx, &creationV1.CreateCollectionsReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CreateCollectionsDbAndCache(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateCollectionsDbAndCache(ctx, &creationV1.CreateCollectionsDbAndCacheReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) EditCollectionsCos(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.EditCollectionsCos(ctx, &creationV1.EditCollectionsCosReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) DeleteCollectionsCache(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.cc.DeleteCollectionsCache(ctx, &creationV1.DeleteCollectionsCacheReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) CollectionsEditReviewPass(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.EditCollections(ctx, &creationV1.EditCollectionsReq{
		Id:   id,
		Auth: auth,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) AddCreationComment(ctx context.Context, createId, createType int32, uuid string) {
	_, _ = r.data.cc.AddCreationComment(ctx, &creationV1.AddCreationCommentReq{
		Uuid:         uuid,
		CreationId:   createId,
		CreationType: createType,
	})
}

func (r *creationRepo) GetCreationUser(ctx context.Context, uuid string) (int32, int32, int32, error) {
	reply, err := r.data.cc.GetCreationUser(ctx, &creationV1.GetCreationUserReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, 0, 0, err
	}
	return reply.Article, reply.Talk, reply.Collect, nil
}
