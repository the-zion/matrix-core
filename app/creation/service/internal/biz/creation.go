package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"golang.org/x/sync/errgroup"
	"sort"
)

type CreationRepo interface {
	GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error)
	GetLastCollectionsDraft(ctx context.Context, uuid string) (*CollectionsDraft, error)
	GetCollectionsContentReview(ctx context.Context, page int32, uuid string) ([]*TextReview, error)
	GetCollectArticleList(ctx context.Context, id, page int32) ([]*Article, error)
	GetCollectArticleCount(ctx context.Context, id int32) (int32, error)
	GetCollectTalkList(ctx context.Context, id, page int32) ([]*Talk, error)
	GetCollectTalkCount(ctx context.Context, id int32) (int32, error)
	GetCollectColumnList(ctx context.Context, id, page int32) ([]*Column, error)
	GetCollectColumnCount(ctx context.Context, id int32) (int32, error)
	GetCollectCount(ctx context.Context, id int32) (int64, error)
	GetCollections(ctx context.Context, id int32, uuid string) (*Collections, error)
	GetCollectionListInfo(ctx context.Context, ids []int32) ([]*Collections, error)
	GetCollectionsList(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsListAll(ctx context.Context, uuid string) ([]*Collections, error)
	GetCollectionsListByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsCount(ctx context.Context, uuid string) (int32, error)
	GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error)
	GetCreationUser(ctx context.Context, uuid string) (*CreationUser, error)
	GetCreationUserVisitor(ctx context.Context, uuid string) (*CreationUser, error)
	GetCollectionsAuth(ctx context.Context, id int32) (int32, error)

	CreateCollectionsDraft(ctx context.Context, uuid string) (int32, error)
	CreateCollectionsFolder(ctx context.Context, id int32, uuid string) error
	CreateCollections(ctx context.Context, id, auth int32, uuid string) error
	CreateCollectionsCache(ctx context.Context, id, auth int32, uuid, mode string) error

	AddCreationUserCollections(ctx context.Context, uuid string, auth int32) error

	EditCollectionsCos(ctx context.Context, id int32, uuid string) error

	UpdateCollectionsCache(ctx context.Context, id, auth int32, uuid, mode string) error

	DeleteCollections(ctx context.Context, id int32, uuid string) error
	DeleteCollect(ctx context.Context, id int32) error
	DeleteCollectionsDraft(ctx context.Context, id int32, uuid string) error
	DeleteCreationCache(ctx context.Context, id, auth int32, uuid string) error

	ReduceCreationUserCollections(ctx context.Context, auth int32, uuid string) error

	SetRecord(ctx context.Context, id, mode int32, uuid, operation, ip string) error
	SetLeaderBoardToCache(ctx context.Context, boardList []*LeaderBoard)

	SendReviewToMq(ctx context.Context, review *CollectionsReview) error
	SendCollectionsToMq(ctx context.Context, collections *Collections, mode string) error
	SendCollections(ctx context.Context, id int32, uuid string) (*CollectionsDraft, error)
	SendCollectionsContentIrregularToMq(ctx context.Context, review *TextReview) error
	SetCollectionsContentIrregular(ctx context.Context, review *TextReview) (*TextReview, error)
	SetCollectionsContentIrregularToCache(ctx context.Context, review *TextReview) error
}

type CreationUseCase struct {
	repo        CreationRepo
	articleRepo ArticleRepo
	talkRepo    TalkRepo
	columnRepo  ColumnRepo
	tm          Transaction
	re          Recovery
	log         *log.Helper
}

func NewCreationUseCase(repo CreationRepo, articleRepo ArticleRepo, talkRepo TalkRepo, columnRepo ColumnRepo, tm Transaction, re Recovery, logger log.Logger) *CreationUseCase {
	return &CreationUseCase{
		repo:        repo,
		articleRepo: articleRepo,
		talkRepo:    talkRepo,
		columnRepo:  columnRepo,
		re:          re,
		tm:          tm,
		log:         log.NewHelper(log.With(logger, "module", "creation/biz/creationUseCase")),
	}
}

func (r *CreationUseCase) GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error) {
	boardList, err := r.repo.GetLeaderBoard(ctx)
	if err != nil {
		return nil, v1.ErrorGetLeaderBoardFailed("get leader board failed: %s", err.Error())
	}

	if len(boardList) != 0 {
		return boardList, nil
	}

	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		return r.getLeaderBoardFromArticle(ctx, &boardList)
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		return r.getLeaderBoardFromTalk(ctx, &boardList)
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		return r.getLeaderBoardFromColumn(ctx, &boardList)
	}))
	err = g.Wait()
	if err != nil {
		return nil, err
	}

	go r.repo.SetLeaderBoardToCache(context.Background(), boardList)
	sort.SliceStable(boardList, func(i, j int) bool {
		return boardList[i].Agree > boardList[j].Agree
	})
	return boardList, nil
}

func (r *CreationUseCase) GetLastCollectionsDraft(ctx context.Context, uuid string) (*CollectionsDraft, error) {
	draft, err := r.repo.GetLastCollectionsDraft(ctx, uuid)
	if kerrors.IsNotFound(err) {
		return nil, v1.ErrorRecordNotFound("last draft not found: %s", err.Error())
	}
	if err != nil {
		return nil, v1.ErrorGetTalkDraftFailed("get last draft failed: %s", err.Error())
	}
	return draft, nil
}

func (r *CreationUseCase) GetCollectionsContentReview(ctx context.Context, page int32, uuid string) ([]*TextReview, error) {
	reviewList, err := r.repo.GetCollectionsContentReview(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetContentReviewFailed("get collections content review failed: %s", err.Error())
	}
	return reviewList, nil
}

func (r *CreationUseCase) getLeaderBoardFromArticle(ctx context.Context, boardList *[]*LeaderBoard) error {
	articleList, err := r.articleRepo.GetArticleHotFromDB(ctx, 1)
	if err != nil {
		return err
	}
	for _, item := range articleList {
		*boardList = append(*boardList, &LeaderBoard{
			Id:    item.ArticleId,
			Agree: item.Agree,
			Uuid:  item.Uuid,
			Mode:  "article",
		})
	}
	return nil
}

func (r *CreationUseCase) getLeaderBoardFromTalk(ctx context.Context, boardList *[]*LeaderBoard) error {
	talkList, err := r.talkRepo.GetTalkHotFromDB(ctx, 1)
	if err != nil {
		return err
	}
	for _, item := range talkList {
		*boardList = append(*boardList, &LeaderBoard{
			Id:    item.TalkId,
			Agree: item.Agree,
			Uuid:  item.Uuid,
			Mode:  "talk",
		})
	}
	return nil
}

func (r *CreationUseCase) getLeaderBoardFromColumn(ctx context.Context, boardList *[]*LeaderBoard) error {
	columnList, err := r.columnRepo.GetColumnHotFromDB(ctx, 1)
	if err != nil {
		return err
	}
	for _, item := range columnList {
		*boardList = append(*boardList, &LeaderBoard{
			Id:    item.ColumnId,
			Agree: item.Agree,
			Uuid:  item.Uuid,
			Mode:  "column",
		})
	}
	return nil
}

func (r *CreationUseCase) GetCollectArticleList(ctx context.Context, id, page int32) ([]*Article, error) {
	articleList, err := r.repo.GetCollectArticleList(ctx, id, page)
	if err != nil {
		return nil, v1.ErrorGetCollectArticleFailed("get collect article list failed: %s", err.Error())
	}
	return articleList, nil
}

func (r *CreationUseCase) GetCollectArticleCount(ctx context.Context, id int32) (int32, error) {
	count, err := r.repo.GetCollectArticleCount(ctx, id)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get collect article count failed: %s", err.Error())
	}
	return count, nil
}

func (r *CreationUseCase) GetCollectTalkList(ctx context.Context, id, page int32) ([]*Talk, error) {
	talkList, err := r.repo.GetCollectTalkList(ctx, id, page)
	if err != nil {
		return nil, v1.ErrorGetTalkListFailed("get collect talk list failed: %s", err.Error())
	}
	return talkList, nil
}

func (r *CreationUseCase) GetCollectTalkCount(ctx context.Context, id int32) (int32, error) {
	count, err := r.repo.GetCollectTalkCount(ctx, id)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get collect talk count failed: %s", err.Error())
	}
	return count, nil
}

func (r *CreationUseCase) GetCollectColumnList(ctx context.Context, id, page int32) ([]*Column, error) {
	columnList, err := r.repo.GetCollectColumnList(ctx, id, page)
	if err != nil {
		return nil, v1.ErrorGetColumnListFailed("get collect column list failed: %s", err.Error())
	}
	return columnList, nil
}

func (r *CreationUseCase) GetCollectColumnCount(ctx context.Context, id int32) (int32, error) {
	count, err := r.repo.GetCollectColumnCount(ctx, id)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get collect column count failed: %s", err.Error())
	}
	return count, nil
}

func (r *CreationUseCase) GetCollections(ctx context.Context, id int32, uuid string) (*Collections, error) {
	collection, err := r.repo.GetCollections(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorGetCollectionFailed("get collection failed: %s", err.Error())
	}
	return collection, nil
}

func (r *CreationUseCase) GetCollectionListInfo(ctx context.Context, ids []int32) ([]*Collections, error) {
	collectionListInfo, err := r.repo.GetCollectionListInfo(ctx, ids)
	if err != nil {
		return nil, v1.ErrorGetCollectionFailed("get collection failed: %s", err.Error())
	}
	return collectionListInfo, nil
}

func (r *CreationUseCase) GetCollectionsList(ctx context.Context, uuid string, page int32) ([]*Collections, error) {
	collections, err := r.repo.GetCollectionsList(ctx, uuid, page)
	if err != nil {
		return nil, v1.ErrorGetCollectionsListFailed("get collections failed: %s", err.Error())
	}
	return collections, nil
}

func (r *CreationUseCase) GetCollectionsListAll(ctx context.Context, uuid string) ([]*Collections, error) {
	collections, err := r.repo.GetCollectionsListAll(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetCollectionsListFailed("get collections all failed: %s", err.Error())
	}
	return collections, nil
}

func (r *CreationUseCase) GetCollectionsListByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error) {
	collections, err := r.repo.GetCollectionsListByVisitor(ctx, uuid, page)
	if err != nil {
		return nil, v1.ErrorGetCollectionsListFailed("get collections failed: %s", err.Error())
	}
	return collections, nil
}

func (r *CreationUseCase) GetCollectionsCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetCollectionsCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get collections count failed: %s", err.Error())
	}
	return count, nil
}

func (r *CreationUseCase) GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetCollectionsVisitorCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get collections count failed: %s", err.Error())
	}
	return count, nil
}

func (r *CreationUseCase) GetCreationUser(ctx context.Context, uuid string) (*CreationUser, error) {
	creationUser, err := r.repo.GetCreationUser(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetCreationUserFailed("get creation user failed: %s", err.Error())
	}
	return creationUser, nil
}

func (r *CreationUseCase) GetCreationUserVisitor(ctx context.Context, uuid string) (*CreationUser, error) {
	creationUser, err := r.repo.GetCreationUserVisitor(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetCreationUserFailed("get creation user failed: %s", err.Error())
	}
	return creationUser, nil
}

func (r *CreationUseCase) SendCollections(ctx context.Context, id int32, uuid, ip string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		draft, err := r.repo.SendCollections(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateCollectionsFailed("create collections failed: %s", err.Error())
		}

		err = r.repo.SetRecord(ctx, id, 4, uuid, "create", ip)
		if err != nil {
			return v1.ErrorSetRecordFailed("set record failed: %s", err.Error())
		}

		err = r.repo.SendReviewToMq(ctx, &CollectionsReview{
			Uuid: draft.Uuid,
			Id:   draft.Id,
			Mode: "create",
		})
		if err != nil {
			return v1.ErrorCreateCollectionsFailed("send create review to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CreationUseCase) CreateCollectionsDraft(ctx context.Context, uuid string) (int32, error) {
	var id int32
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = r.repo.CreateCollectionsDraft(ctx, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create collections draft failed: %s", err.Error())
		}

		err = r.repo.CreateCollectionsFolder(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create collections folder failed: %s", err.Error())
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CreationUseCase) CreateCollections(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.SendCollectionsToMq(ctx, &Collections{
		CollectionsId: id,
		Uuid:          uuid,
		Auth:          auth,
	}, "create_collections_db_and_cache")
	if err != nil {
		return v1.ErrorCreateCollectionsFailed("create collections to mq failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) CreateCollectionsDbAndCache(ctx context.Context, id, auth int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteCollectionsDraft(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateCollectionsFailed("delete collections draft failed: %s", err.Error())
		}

		err = r.repo.CreateCollections(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorCreateCollectionsFailed("create collections failed: %s", err.Error())
		}

		err = r.repo.AddCreationUserCollections(ctx, uuid, auth)
		if err != nil {
			return v1.ErrorCreateCollectionsFailed("add creation collections failed: %s", err.Error())
		}

		err = r.repo.CreateCollectionsCache(ctx, id, auth, uuid, "create")
		if err != nil {
			return v1.ErrorCreateCollectionsFailed("create collections cache failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CreationUseCase) CollectionsContentIrregular(ctx context.Context, review *TextReview) error {
	err := r.repo.SendCollectionsContentIrregularToMq(ctx, review)
	if err != nil {
		return v1.ErrorSetContentIrregularFailed("set collections content irregular to mq failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) AddCollectionsContentReviewDbAndCache(ctx context.Context, review *TextReview) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		review, err := r.repo.SetCollectionsContentIrregular(ctx, review)
		if err != nil {
			return v1.ErrorSetContentIrregularFailed("set collections content irregular failed: %s", err.Error())
		}

		err = r.repo.SetCollectionsContentIrregularToCache(ctx, review)
		if err != nil {
			return v1.ErrorSetContentIrregularFailed("set collections content irregular to cache failed: %s", err.Error())
		}

		return nil
	})
}

func (r *CreationUseCase) SendCollectionsEdit(ctx context.Context, id int32, uuid, ip string) error {
	collections, err := r.repo.GetCollections(ctx, id, uuid)
	if err != nil {
		return v1.ErrorGetCollectionFailed("get collections failed: %s", err.Error())
	}

	if collections.Uuid != uuid {
		return v1.ErrorEditCollectionsFailed("send collections edit failed: no auth")
	}

	err = r.repo.SetRecord(ctx, id, 4, uuid, "edit", ip)
	if err != nil {
		return v1.ErrorSetRecordFailed("set record failed: %s", err.Error())
	}

	err = r.repo.SendReviewToMq(ctx, &CollectionsReview{
		Uuid: collections.Uuid,
		Id:   id,
		Mode: "edit",
	})
	if err != nil {
		return v1.ErrorEditCollectionsFailed("send edit review to mq failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) EditCollections(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.SendCollectionsToMq(ctx, &Collections{
		CollectionsId: id,
		Auth:          auth,
		Uuid:          uuid,
	}, "edit_collections_cos")
	if err != nil {
		return v1.ErrorEditTalkFailed("edit collections to mq failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) EditCollectionsCos(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.UpdateCollectionsCache(ctx, id, auth, uuid, "edit")
	if err != nil {
		return v1.ErrorEditTalkFailed("edit collections cache failed: %s", err.Error())
	}

	err = r.repo.EditCollectionsCos(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditTalkFailed("edit collections cache failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) DeleteCollections(ctx context.Context, id int32, uuid string) error {
	err := r.repo.SendCollectionsToMq(ctx, &Collections{
		CollectionsId: id,
		Uuid:          uuid,
	}, "delete_collections_cache")
	if err != nil {
		return v1.ErrorDeleteCollectionsFailed("delete collections to mq failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) DeleteCollectionsCache(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		auth, err := r.repo.GetCollectionsAuth(ctx, id)
		if err != nil {
			return v1.ErrorDeleteCollectionsFailed("get collections auth failed: %s", err.Error())
		}

		err = r.repo.DeleteCollections(ctx, id, uuid)
		if err != nil {
			return v1.ErrorDeleteCollectionsFailed("delete collections failed: %s", err.Error())
		}

		err = r.repo.DeleteCollect(ctx, id)
		if err != nil {
			return v1.ErrorDeleteCollectionsFailed("delete collect failed: %s", err.Error())
		}

		err = r.repo.ReduceCreationUserCollections(ctx, auth, uuid)
		if err != nil {
			return v1.ErrorDeleteCollectionsFailed("delete creation user collections failed: %s", err.Error())
		}

		err = r.repo.DeleteCreationCache(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorDeleteCollectionsFailed("delete creation cache failed: %s", err.Error())
		}
		return nil
	})
}
