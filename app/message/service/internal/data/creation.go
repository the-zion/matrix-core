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

func (r *creationRepo) CreateArticleCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateArticleCacheAndSearch(ctx, &creationV1.CreateArticleCacheAndSearchReq{
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

func (r *creationRepo) CreateTalkCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateTalkCacheAndSearch(ctx, &creationV1.CreateTalkCacheAndSearchReq{
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

func (r *creationRepo) CreateColumnCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	_, err := r.data.cc.CreateColumnCacheAndSearch(ctx, &creationV1.CreateColumnCacheAndSearchReq{
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
