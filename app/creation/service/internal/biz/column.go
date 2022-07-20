package biz

import "C"
import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type ColumnRepo interface {
	CreateColumnDraft(ctx context.Context, uuid string) (int32, error)
	CreateColumnFolder(ctx context.Context, id int32, uuid string) error
	CreateColumn(ctx context.Context, id, auth int32, uuid string) error
	CreateColumnStatistic(ctx context.Context, id, auth int32, uuid string) error
	CreateColumnCache(ctx context.Context, id, auth int32, uuid string) error
	CreateColumnSearch(ctx context.Context, id int32, uuid string) error

	UpdateColumnCache(ctx context.Context, id, auth int32, uuid string) error
	EditColumnCos(ctx context.Context, id int32, uuid string) error
	EditColumnSearch(ctx context.Context, id int32, uuid string) error

	DeleteColumnDraft(ctx context.Context, id int32, uuid string) error
	DeleteColumnStatistic(ctx context.Context, id int32, uuid string) error
	DeleteColumn(ctx context.Context, id int32, uuid string) error
	DeleteColumnCache(ctx context.Context, id int32, uuid string) error
	DeleteColumnSearch(ctx context.Context, id int32, uuid string) error

	GetColumn(ctx context.Context, id int32) (*Column, error)
	GetLastColumnDraft(ctx context.Context, uuid string) (*ColumnDraft, error)
	GetColumnList(ctx context.Context, page int32) ([]*Column, error)
	GetColumnListHot(ctx context.Context, page int32) ([]*ColumnStatistic, error)
	GetUserColumnList(ctx context.Context, page int32, uuid string) ([]*Column, error)
	GetUserColumnListVisitor(ctx context.Context, page int32, uuid string) ([]*Column, error)
	GetColumnCount(ctx context.Context, uuid string) (int32, error)
	GetColumnCountVisitor(ctx context.Context, uuid string) (int32, error)
	GetColumnListStatistic(ctx context.Context, ids []int32) ([]*ColumnStatistic, error)

	SendColumn(ctx context.Context, id int32, uuid string) (*ColumnDraft, error)
	SendColumnToMq(ctx context.Context, column *Column, mode string) error
	SendReviewToMq(ctx context.Context, review *ColumnReview) error

	FreezeColumnCos(ctx context.Context, id int32, uuid string) error
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

func (r *ColumnUseCase) SendColumn(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		draft, err := r.repo.SendColumn(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateColumnFailed("send column failed: %s", err.Error())
		}
		err = r.repo.SendReviewToMq(ctx, &ColumnReview{
			Uuid: draft.Uuid,
			Id:   draft.Id,
			Mode: "create",
		})
		if err != nil {
			return v1.ErrorCreateColumnFailed("send create review to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) SendColumnEdit(ctx context.Context, id int32, uuid string) error {
	column, err := r.repo.GetColumn(ctx, id)
	if err != nil {
		return v1.ErrorGetColumnFailed("get column failed: %s", err.Error())
	}

	if column.Uuid != uuid {
		return v1.ErrorGetColumnFailed("send column edit failed: no auth")
	}

	err = r.repo.SendReviewToMq(ctx, &ColumnReview{
		Uuid: column.Uuid,
		Id:   column.ColumnId,
		Mode: "edit",
	})
	if err != nil {
		return v1.ErrorGetColumnFailed("send edit review to mq failed: %s", err.Error())
	}
	return nil
}

func (r *ColumnUseCase) CreateColumn(ctx context.Context, id, auth int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteColumnDraft(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateColumnFailed("delete column draft failed: %s", err.Error())
		}

		err = r.repo.CreateColumn(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorCreateColumnFailed("create column failed: %s", err.Error())
		}

		err = r.repo.CreateColumnStatistic(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorCreateColumnFailed("create column statistic failed: %s", err.Error())
		}

		err = r.repo.SendColumnToMq(ctx, &Column{
			ColumnId: id,
			Uuid:     uuid,
			Auth:     auth,
		}, "create_column_cache_and_search")
		if err != nil {
			return v1.ErrorCreateColumnFailed("create column to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) EditColumn(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.SendColumnToMq(ctx, &Column{
		ColumnId: id,
		Auth:     auth,
		Uuid:     uuid,
	}, "edit_column_cos_and_search")
	if err != nil {
		return v1.ErrorEditColumnFailed("edit column to mq failed: %s", err.Error())
	}
	return nil
}

func (r *ColumnUseCase) DeleteColumn(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteColumn(ctx, id, uuid)
		if err != nil {
			return v1.ErrorDeleteColumnFailed("delete column failed: %s", err.Error())
		}

		err = r.repo.DeleteColumnStatistic(ctx, id, uuid)
		if err != nil {
			return v1.ErrorDeleteColumnFailed("delete column statistic failed: %s", err.Error())
		}

		err = r.repo.SendColumnToMq(ctx, &Column{
			ColumnId: id,
			Uuid:     uuid,
		}, "delete_column_cache_and_search")
		if err != nil {
			return v1.ErrorDeleteColumnFailed("delete column to mq failed: %s", err.Error())
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

func (r *ColumnUseCase) EditColumnCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.UpdateColumnCache(ctx, id, auth, uuid)
	if err != nil {
		return v1.ErrorEditColumnFailed("edit column cache failed: %s", err.Error())
	}

	err = r.repo.EditColumnCos(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditColumnFailed("edit column cache failed: %s", err.Error())
	}

	err = r.repo.EditColumnSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditColumnFailed("edit column search failed: %s", err.Error())
	}
	return nil
}

func (r *ColumnUseCase) DeleteColumnCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	err := r.repo.DeleteColumnCache(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteColumnFailed("delete column cache failed: %s", err.Error())
	}

	err = r.repo.FreezeColumnCos(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteColumnFailed("freeze column cos failed: %s", err.Error())
	}

	err = r.repo.DeleteColumnSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteColumnFailed("delete column search failed: %s", err.Error())
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

func (r *ColumnUseCase) GetUserColumnList(ctx context.Context, page int32, uuid string) ([]*Column, error) {
	columnList, err := r.repo.GetUserColumnList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetColumnListFailed("get user column list failed: %s", err.Error())
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetUserColumnListVisitor(ctx context.Context, page int32, uuid string) ([]*Column, error) {
	columnList, err := r.repo.GetUserColumnListVisitor(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetColumnListFailed("get user column list visitor failed: %s", err.Error())
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetColumnCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetColumnCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get column count failed: %s", err.Error())
	}
	return count, nil
}

func (r *ColumnUseCase) GetColumnCountVisitor(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetColumnCountVisitor(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get column count visitor failed: %s", err.Error())
	}
	return count, nil
}

func (r *ColumnUseCase) GetColumnListStatistic(ctx context.Context, ids []int32) ([]*ColumnStatistic, error) {
	statisticList, err := r.repo.GetColumnListStatistic(ctx, ids)
	if err != nil {
		return nil, v1.ErrorGetStatisticFailed("get column list statistic failed: %s", err.Error())
	}
	return statisticList, nil
}
