package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
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

func (r *creationRepo) ToReviewArticle(mode string, msgs ...*primitive.MessageExt) {
	for _, i := range msgs {
		m := map[string]interface{}{"Uuid": "", "Id": ""}
		err := json.Unmarshal(i.Body, &m)
		if err != nil {
			log.Errorf("fail to unmarshal profile: err(%v)", err)
		}

		opt := &cos.PutTextAuditingJobOptions{
			InputObject: "article/" + strconv.Itoa(int(m["Id"].(float64))) + "/" + m["Uuid"].(string),
			Conf: &cos.TextAuditingJobConf{
				CallbackVersion: "Detail",
				Callback:        r.data.cosCreationCli.callback[mode],
			},
		}

		res, _, err := r.data.cosCreationCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
		fmt.Println(res)
		if err != nil {
			log.Errorf("fail to send article review request to cos: err(%v)", err)
		}
	}
}

func (r *creationRepo) ArticleDraftReviewPass(ctx context.Context, uuid string, id int32) error {
	_, err := r.data.cc.CreateArticle(ctx, &creationV1.CreateArticleReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}
