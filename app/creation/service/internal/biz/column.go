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
	AddColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error
	AddColumnIncludesToCache(ctx context.Context, id, articleId int32, uuid string) error

	UpdateColumnCache(ctx context.Context, id, auth int32, uuid string) error
	EditColumnCos(ctx context.Context, id int32, uuid string) error
	EditColumnSearch(ctx context.Context, id int32, uuid string) error

	DeleteColumnDraft(ctx context.Context, id int32, uuid string) error
	DeleteColumnStatistic(ctx context.Context, id int32, uuid string) error
	DeleteColumn(ctx context.Context, id int32, uuid string) error
	DeleteColumnCache(ctx context.Context, id int32, uuid string) error
	DeleteColumnSearch(ctx context.Context, id int32, uuid string) error
	DeleteColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error
	DeleteColumnIncludesFromCache(ctx context.Context, id, articleId int32, uuid string) error

	GetColumn(ctx context.Context, id int32) (*Column, error)
	GetLastColumnDraft(ctx context.Context, uuid string) (*ColumnDraft, error)
	GetColumnList(ctx context.Context, page int32) ([]*Column, error)
	GetColumnListHot(ctx context.Context, page int32) ([]*ColumnStatistic, error)
	GetColumnHotFromDB(ctx context.Context, page int32) ([]*ColumnStatistic, error)
	GetUserColumnList(ctx context.Context, page int32, uuid string) ([]*Column, error)
	GetUserColumnListVisitor(ctx context.Context, page int32, uuid string) ([]*Column, error)
	GetColumnCount(ctx context.Context, uuid string) (int32, error)
	GetColumnCountVisitor(ctx context.Context, uuid string) (int32, error)
	GetColumnStatistic(ctx context.Context, id int32) (*ColumnStatistic, error)
	GetColumnListStatistic(ctx context.Context, ids []int32) ([]*ColumnStatistic, error)
	GetColumnAgreeJudge(ctx context.Context, id int32, uuid string) (bool, error)
	GetColumnCollectJudge(ctx context.Context, id int32, uuid string) (bool, error)
	GetSubscribeList(ctx context.Context, page int32, uuid string) ([]*Subscribe, error)
	GetSubscribeListCount(ctx context.Context, uuid string) (int32, error)
	GetColumnSubscribes(ctx context.Context, uuid string, ids []int32) ([]*Subscribe, error)
	GetColumnSearch(ctx context.Context, page int32, search, time string) ([]*ColumnSearch, int32, error)

	SendColumn(ctx context.Context, id int32, uuid string) (*ColumnDraft, error)
	SendColumnToMq(ctx context.Context, column *Column, mode string) error
	SendReviewToMq(ctx context.Context, review *ColumnReview) error

	FreezeColumnCos(ctx context.Context, id int32, uuid string) error

	SetColumnAgree(ctx context.Context, id int32, uuid string) error
	SetColumnAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error
	SendColumnStatisticToMq(ctx context.Context, uuid, mode string) error
	SetColumnView(ctx context.Context, id int32, uuid string) error
	SetColumnViewToCache(ctx context.Context, id int32, uuid string) error
	SetColumnUserCollect(ctx context.Context, id, collectionsId int32, userUuid string) error
	SetColumnCollect(ctx context.Context, id int32, uuid string) error
	SetColumnCollectToCache(ctx context.Context, id int32, uuid, userUuid string) error

	CancelColumnAgree(ctx context.Context, id int32, uuid string) error
	CancelColumnAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelColumnCollect(ctx context.Context, id int32, uuid string) error
	CancelColumnCollectFromCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelColumnUserCollect(ctx context.Context, id int32, userUuid string) error
	CancelSubscribeColumn(ctx context.Context, id int32, uuid string) error

	SubscribeColumn(ctx context.Context, id int32, author, uuid string) error
	SubscribeJudge(ctx context.Context, id int32, uuid string) (bool, error)
}

type ColumnUseCase struct {
	repo         ColumnRepo
	creationRepo CreationRepo
	tm           Transaction
	re           Recovery
	log          *log.Helper
}

func NewColumnUseCase(repo ColumnRepo, re Recovery, creationRepo CreationRepo, tm Transaction, logger log.Logger) *ColumnUseCase {
	return &ColumnUseCase{
		repo:         repo,
		creationRepo: creationRepo,
		tm:           tm,
		re:           re,
		log:          log.NewHelper(log.With(logger, "module", "creation/biz/columnUseCase")),
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

func (r *ColumnUseCase) GetColumnSearch(ctx context.Context, page int32, search, time string) ([]*ColumnSearch, int32, error) {
	columnList, total, err := r.repo.GetColumnSearch(ctx, page, search, time)
	if err != nil {
		return nil, 0, v1.ErrorGetColumnSearchFailed("get column search failed: %s", err.Error())
	}
	return columnList, total, nil
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

func (r *ColumnUseCase) SubscribeColumn(ctx context.Context, id int32, author, uuid string) error {
	if author == uuid {
		return v1.ErrorSubscribeColumnFailed("author and uuid are the same")
	}

	err := r.repo.SubscribeColumn(ctx, id, author, uuid)
	if err != nil {
		return v1.ErrorSubscribeColumnFailed("subscribe column failed: %s", err.Error())
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *ColumnUseCase) CancelSubscribeColumn(ctx context.Context, id int32, uuid string) error {
	err := r.repo.CancelSubscribeColumn(ctx, id, uuid)
	if err != nil {
		return v1.ErrorCancelSubscribeColumnFailed("cancel subscribe column failed: %s", err.Error())
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *ColumnUseCase) SubscribeJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	subscribe, err := r.repo.SubscribeJudge(ctx, id, uuid)
	if err != nil {
		return false, v1.ErrorSubscribeColumnJudgeFailed("get subscribe column judge failed: %s", err.Error())
	}
	if err != nil {
		return false, err
	}
	return subscribe, nil
}

func (r *ColumnUseCase) SendColumn(ctx context.Context, id int32, uuid, ip string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		draft, err := r.repo.SendColumn(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateColumnFailed("send column failed: %s", err.Error())
		}

		err = r.creationRepo.SetRecord(ctx, id, 2, uuid, "create", ip)
		if err != nil {
			return v1.ErrorSetRecordFailed("set record failed: %s", err.Error())
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

func (r *ColumnUseCase) SendColumnEdit(ctx context.Context, id int32, uuid, ip string) error {
	column, err := r.repo.GetColumn(ctx, id)
	if err != nil {
		return v1.ErrorGetColumnFailed("get column failed: %s", err.Error())
	}

	if column.Uuid != uuid {
		return v1.ErrorGetColumnFailed("send column edit failed: no auth")
	}

	err = r.creationRepo.SetRecord(ctx, id, 2, uuid, "edit", ip)
	if err != nil {
		return v1.ErrorSetRecordFailed("set record failed: %s", err.Error())
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
	err := r.repo.SendColumnToMq(ctx, &Column{
		ColumnId: id,
		Uuid:     uuid,
		Auth:     auth,
	}, "create_column_db_cache_and_search")
	if err != nil {
		return v1.ErrorCreateColumnFailed("create column to mq failed: %s", err.Error())
	}
	return nil
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

func (r *ColumnUseCase) CreateColumnDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
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

		err = r.repo.CreateColumnCache(ctx, id, auth, uuid)
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
	})
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

	if auth == 2 {
		return nil
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

func (r *ColumnUseCase) GetColumnStatistic(ctx context.Context, id int32) (*ColumnStatistic, error) {
	statistic, err := r.repo.GetColumnStatistic(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetStatisticFailed("get column statistic failed: %s", err.Error())
	}
	return statistic, nil
}

func (r *ColumnUseCase) ColumnStatisticJudge(ctx context.Context, id int32, uuid string) (*ColumnStatisticJudge, error) {
	agree, err := r.repo.GetColumnAgreeJudge(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorGetStatisticJudgeFailed("get column statistic judge failed: %s", err.Error())
	}
	collect, err := r.repo.GetColumnCollectJudge(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorGetStatisticJudgeFailed("get column statistic judge failed: %s", err.Error())
	}
	return &ColumnStatisticJudge{
		Agree:   agree,
		Collect: collect,
	}, nil
}

func (r *ColumnUseCase) GetSubscribeList(ctx context.Context, page int32, uuid string) ([]*Subscribe, error) {
	subscribe, err := r.repo.GetSubscribeList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetSubscribeColumnListFailed("get subscribe column list failed: %s", err.Error())
	}
	return subscribe, nil
}

func (r *ColumnUseCase) GetSubscribeListCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetSubscribeListCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get subscribe column list count failed: %s", err.Error())
	}
	return count, nil
}

func (r *ColumnUseCase) GetColumnSubscribes(ctx context.Context, uuid string, ids []int32) ([]*Subscribe, error) {
	subscribe, err := r.repo.GetColumnSubscribes(ctx, uuid, ids)
	if err != nil {
		return nil, v1.ErrorGetColumnSubscribesFailed("get column subscribes failed: %s", err.Error())
	}
	return subscribe, nil
}

func (r *ColumnUseCase) SetColumnAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetColumnAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set column agree failed: %s", err.Error())
		}
		err = r.repo.SetColumnAgreeToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set column agree to cache failed: %s", err.Error())
		}
		err = r.repo.SendColumnStatisticToMq(ctx, uuid, "agree")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set column agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) SetColumnView(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetColumnView(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetViewFailed("set column view failed: %s", err.Error())
		}
		err = r.repo.SetColumnViewToCache(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetViewFailed("set column view to cache failed: %s", err.Error())
		}
		err = r.repo.SendColumnStatisticToMq(ctx, uuid, "view")
		if err != nil {
			return v1.ErrorSetViewFailed("set column view to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) CancelColumnAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelColumnAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel column agree failed: %s", err.Error())
		}
		err = r.repo.CancelColumnAgreeFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel column agree from cache failed: %s", err.Error())
		}
		err = r.repo.SendColumnStatisticToMq(ctx, uuid, "agree_cancel")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("set column agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) SetColumnCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetColumnUserCollect(ctx, id, collectionsId, userUuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set column user collect failed: %s", err.Error())
		}
		err = r.repo.SetColumnCollect(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set column collect failed: %s", err.Error())
		}
		err = r.repo.SetColumnCollectToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set column collect to cache failed: %s", err.Error())
		}
		err = r.repo.SendColumnStatisticToMq(ctx, uuid, "collect")
		if err != nil {
			return v1.ErrorSetCollectFailed("set column collect to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) CancelColumnCollect(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelColumnUserCollect(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel column user collect failed: %s", err.Error())
		}
		err = r.repo.CancelColumnCollect(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel column collect failed: %s", err.Error())
		}
		err = r.repo.CancelColumnCollectFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel column collect to cache failed: %s", err.Error())
		}
		err = r.repo.SendColumnStatisticToMq(ctx, uuid, "collect_cancel")
		if err != nil {
			return v1.ErrorCancelCollectFailed("set column collect to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) AddColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.AddColumnIncludes(ctx, id, articleId, uuid)
		if err != nil {
			return v1.ErrorAddColumnIncludesFailed("add column includes failed: %s", err.Error())
		}

		err = r.repo.AddColumnIncludesToCache(ctx, id, articleId, uuid)
		if err != nil {
			return v1.ErrorAddColumnIncludesFailed("add column includes to cache failed: %s", err.Error())
		}
		return nil
	})
}

func (r *ColumnUseCase) DeleteColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error {
	err := r.repo.DeleteColumnIncludes(ctx, id, articleId, uuid)
	if err != nil {
		return v1.ErrorDeleteColumnIncludesFailed("delete column includes failed: %s", err.Error())
	}

	err = r.repo.DeleteColumnIncludesFromCache(ctx, id, articleId, uuid)
	if err != nil {
		return v1.ErrorDeleteColumnIncludesFailed("delete column includes from cache failed: %s", err.Error())
	}
	return nil
}
