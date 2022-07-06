package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"strconv"
	"strings"
)

var _ biz.CreationRepo = (*creationRepo)(nil)

type creationRepo struct {
	data *Data
	log  *log.Helper
}

func NewCreationRepo(data *Data, logger log.Logger) biz.CreationRepo {
	return &creationRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/creation")),
	}
}

func (r *creationRepo) GetLeaderBoard(ctx context.Context) ([]*biz.LeaderBoard, error) {
	return r.getLeaderBoardFromCache(ctx)
}

func (r *creationRepo) getLeaderBoardFromCache(ctx context.Context) ([]*biz.LeaderBoard, error) {
	list, err := r.data.redisCli.ZRevRange(ctx, "leaderboard", 0, 9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get leader board from cache")
	}

	board := make([]*biz.LeaderBoard, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		board = append(board, &biz.LeaderBoard{
			Id:   int32(id),
			Uuid: member[1],
			Mode: member[2],
		})
	}
	return board, nil
}
