package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"golang.org/x/sync/errgroup"
	"sort"
)

type CreationRepo interface {
	GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error)
	GetCollectArticle(ctx context.Context, id, page int32) ([]*Article, error)
	GetCollectArticleCount(ctx context.Context, id int32) (int32, error)
	GetCollectTalk(ctx context.Context, id, page int32) ([]*Talk, error)
	GetCollectTalkCount(ctx context.Context, id int32) (int32, error)
	GetCollectColumn(ctx context.Context, id, page int32) ([]*Column, error)
	GetCollectColumnCount(ctx context.Context, id int32) (int32, error)
	GetCollectCount(ctx context.Context, id int32) (int64, error)
	GetCollection(ctx context.Context, id int32, uuid string) (*Collections, error)
	GetCollectionListInfo(ctx context.Context, ids []int32) ([]*Collections, error)
	GetCollections(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsAll(ctx context.Context, uuid string) ([]*Collections, error)
	GetCollectionsByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsCount(ctx context.Context, uuid string) (int32, error)
	GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error)
	GetCreationUser(ctx context.Context, uuid string) (*CreationUser, error)
	GetCreationUserVisitor(ctx context.Context, uuid string) (*CreationUser, error)
	CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error
	EditCollections(ctx context.Context, id int32, uuid, name, introduce string, auth int32) error
	DeleteCollections(ctx context.Context, id int32, uuid string) error
	SetRecord(ctx context.Context, id, mode int32, uuid, operation, ip string) error
	SetLeaderBoardToCache(ctx context.Context, boardList []*LeaderBoard)
}

type CreationUseCase struct {
	repo        CreationRepo
	articleRepo ArticleRepo
	talkRepo    TalkRepo
	columnRepo  ColumnRepo
	re          Recovery
	log         *log.Helper
}

func NewCreationUseCase(repo CreationRepo, articleRepo ArticleRepo, talkRepo TalkRepo, columnRepo ColumnRepo, re Recovery, logger log.Logger) *CreationUseCase {
	return &CreationUseCase{
		repo:        repo,
		articleRepo: articleRepo,
		talkRepo:    talkRepo,
		columnRepo:  columnRepo,
		re:          re,
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

func (r *CreationUseCase) GetCollectArticle(ctx context.Context, id, page int32) ([]*Article, error) {
	articleList, err := r.repo.GetCollectArticle(ctx, id, page)
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

func (r *CreationUseCase) GetCollectTalk(ctx context.Context, id, page int32) ([]*Talk, error) {
	talkList, err := r.repo.GetCollectTalk(ctx, id, page)
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

func (r *CreationUseCase) GetCollectColumn(ctx context.Context, id, page int32) ([]*Column, error) {
	columnList, err := r.repo.GetCollectColumn(ctx, id, page)
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

func (r *CreationUseCase) GetCollection(ctx context.Context, id int32, uuid string) (*Collections, error) {
	collection, err := r.repo.GetCollection(ctx, id, uuid)
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

func (r *CreationUseCase) GetCollections(ctx context.Context, uuid string, page int32) ([]*Collections, error) {
	collections, err := r.repo.GetCollections(ctx, uuid, page)
	if err != nil {
		return nil, v1.ErrorGetCollectionsFailed("get collections failed: %s", err.Error())
	}
	return collections, nil
}

func (r *CreationUseCase) GetCollectionsAll(ctx context.Context, uuid string) ([]*Collections, error) {
	collections, err := r.repo.GetCollectionsAll(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetCollectionsFailed("get collections all failed: %s", err.Error())
	}
	return collections, nil
}

func (r *CreationUseCase) GetCollectionsByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error) {
	collections, err := r.repo.GetCollectionsByVisitor(ctx, uuid, page)
	if err != nil {
		return nil, v1.ErrorGetCollectionsFailed("get collections failed: %s", err.Error())
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

func (r *CreationUseCase) CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error {
	err := r.repo.CreateCollections(ctx, uuid, name, introduce, auth)
	if err != nil {
		return v1.ErrorCreateCollectionsFailed("create collections failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) EditCollections(ctx context.Context, id int32, uuid, name, introduce string, auth int32) error {
	err := r.repo.EditCollections(ctx, id, uuid, name, introduce, auth)
	if err != nil {
		return v1.ErrorEditCollectionsFailed("edit collections failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) DeleteCollections(ctx context.Context, id int32, uuid string) error {
	count, err := r.repo.GetCollectCount(ctx, id)
	if err != nil {
		return v1.ErrorGetCountFailed("get collection count failed: %s", err.Error())
	}

	if count == 1 {
		return v1.ErrorNotEmpty("collections is not empty")
	}

	err = r.repo.DeleteCollections(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteCollectionsFailed("delete collections failed: %s", err.Error())
	}
	return nil
}
