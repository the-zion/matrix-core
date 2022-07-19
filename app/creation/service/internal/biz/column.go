package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type ColumnRepo interface {
	GetLastColumnDraft(ctx context.Context, uuid string) (*ColumnDraft, error)
	CreateColumnDraft(ctx context.Context, uuid string) (int32, error)
	CreateColumnFolder(ctx context.Context, id int32, uuid string) error
	CreateColumn(ctx context.Context, uuid, name, introduce, cover string, auth int32) (int32, error)
	CreateColumnStatistic(ctx context.Context, id, auth int32, uuid string) error
	SendColumnToMq(ctx context.Context, column *Column, mode string) error
	CreateColumnCache(ctx context.Context, id, auth int32, uuid string) error
	CreateColumnSearch(ctx context.Context, id int32, uuid string) error
	GetColumnList(ctx context.Context, page int32) ([]*Column, error)
	GetColumnListHot(ctx context.Context, page int32) ([]*ColumnStatistic, error)
	GetColumnListStatistic(ctx context.Context, ids []int32) ([]*ColumnStatistic, error)
}

type ColumnUseCase struct {
	repo ColumnRepo
	tm   Transaction
	log  *log.Helper
}

func NewColumnUseCase(repo ColumnRepo, tm Transaction, logger log.Logger) *ColumnUseCase {
	return &ColumnUseCase{
		repo: repo,
		tm:   tm,
		log:  log.NewHelper(log.With(logger, "module", "creation/biz/columnUseCase")),
	}
}

func (r *ColumnUseCase) GetLastColumnDraft(ctx context.Context, uuid string) (*ColumnDraft, error) {
	draft, err := r.repo.GetLastColumnDraft(ctx, uuid)
	if kerrors.IsNotFound(err) {
		return nil, v1.ErrorRecordNotFound("last draft not found: %s", err.Error())
	}
	if err != nil {
		return nil, v1.ErrorGetColumnDraftFailed("get last draft failed: %s", err.Error())
	}
	return draft, nil
}

func (r *ColumnUseCase) CreateColumnDraft(ctx context.Context, uuid string) (int32, error) {
	var id int32
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = r.repo.CreateColumnDraft(ctx, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create column draft failed: %s", err.Error())
		}

		err = r.repo.CreateColumnFolder(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create column folder failed: %s", err.Error())
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ColumnUseCase) CreateColumn(ctx context.Context, uuid, name, introduce, cover string, auth int32) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		id, err := r.repo.CreateColumn(ctx, uuid, name, introduce, cover, auth)
		if err != nil {
			return v1.ErrorCreateColumnFailed("create column failed: %s", err.Error())
		}

		err = r.repo.CreateColumnStatistic(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorCreateColumnFailed("create column statistic failed: %s", err.Error())
		}

		err = r.repo.SendColumnToMq(ctx, &Column{
			Id:   id,
			Uuid: uuid,
			Auth: auth,
		}, "create_column_cache_and_search")
		if err != nil {
			return v1.ErrorCreateColumnFailed("create column to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) CreateColumnCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.CreateColumnCache(ctx, id, auth, uuid)
	if err != nil {
		return v1.ErrorCreateColumnFailed("create column cache failed: %s", err.Error())
	}

	if auth == 2 {
		return nil
	}

	err = r.repo.CreateColumnSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorCreateColumnFailed("create column search failed: %s", err.Error())
	}
	return nil
}

func (r *ColumnUseCase) GetColumnList(ctx context.Context, page int32) ([]*Column, error) {
	columnList, err := r.repo.GetColumnList(ctx, page)
	if err != nil {
		return nil, v1.ErrorGetColumnListFailed("get column list failed: %s", err.Error())
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetColumnListHot(ctx context.Context, page int32) ([]*ColumnStatistic, error) {
	columnList, err := r.repo.GetColumnListHot(ctx, page)
	if err != nil {
		return nil, v1.ErrorGetColumnListFailed("get column hot list failed: %s", err.Error())
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetColumnListStatistic(ctx context.Context, ids []int32) ([]*ColumnStatistic, error) {
	statisticList, err := r.repo.GetColumnListStatistic(ctx, ids)
	if err != nil {
		return nil, v1.ErrorGetStatisticFailed("get column list statistic failed: %s", err.Error())
	}
	return statisticList, nil
}
