package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/errgroup"
)

type CreationRepo interface {
	GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error)
	GetCollectArticleList(ctx context.Context, id, page int32) ([]*Article, error)
	GetCollectArticleCount(ctx context.Context, id int32) (int32, error)
	GetCollectTalkList(ctx context.Context, id, page int32) ([]*Talk, error)
	GetCollectTalkCount(ctx context.Context, id int32) (int32, error)
	GetCollectColumnList(ctx context.Context, id, page int32) ([]*Column, error)
	GetCollectColumnCount(ctx context.Context, id int32) (int32, error)
	GetLastCollectionsDraft(ctx context.Context, uuid string) (*CollectionsDraft, error)
	GetCollectionsContentReview(ctx context.Context, page int32, uuid string) ([]*CreationContentReview, error)
	GetUserTimeLineListVisitor(ctx context.Context, page int32, uuid string) ([]*TimeLine, error)
	GetCollections(ctx context.Context, id int32, uuid string) (*Collections, error)
	GetCollectionListInfo(ctx context.Context, collectionsList []*Collections) ([]*Collections, error)
	GetCollectionsList(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsListAll(ctx context.Context, uuid string) ([]*Collections, error)
	GetCollectionsListByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsCount(ctx context.Context, uuid string) (int32, error)
	GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error)
	GetCreationUser(ctx context.Context, uuid string) (*CreationUser, error)
	GetCreationUserVisitor(ctx context.Context, uuid string) (*CreationUser, error)
	CreateCollectionsDraft(ctx context.Context, uuid string) (int32, error)
	SendCollections(ctx context.Context, id int32, uuid, ip string) error
	SendCollectionsEdit(ctx context.Context, id int32, uuid, ip string) error
	DeleteCollections(ctx context.Context, id int32, uuid string) error
}

type ArticleRepo interface {
	CreateArticleDraft(ctx context.Context, uuid string) (int32, error)
	GetLastArticleDraft(ctx context.Context, uuid string) (*ArticleDraft, error)
	GetArticleList(ctx context.Context, page int32) ([]*Article, error)
	GetArticleListHot(ctx context.Context, page int32) ([]*Article, error)
	GetUserArticleListAll(ctx context.Context, uuid string) ([]*Article, error)
	GetColumnArticleList(ctx context.Context, id int32) ([]*Article, error)
	GetArticleCount(ctx context.Context, uuid string) (int32, error)
	GetArticleCountVisitor(ctx context.Context, uuid string) (int32, error)
	GetUserArticleList(ctx context.Context, page int32, uuid string) ([]*Article, error)
	GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*Article, error)
	GetArticleStatistic(ctx context.Context, id int32, uuid string) (*ArticleStatistic, error)
	GetUserArticleAgree(ctx context.Context, uuid string) (map[int32]bool, error)
	GetUserArticleCollect(ctx context.Context, uuid string) (map[int32]bool, error)
	GetArticleListStatistic(ctx context.Context, articleList []*Article) ([]*ArticleStatistic, error)
	GetArticleDraftList(ctx context.Context, uuid string) ([]*ArticleDraft, error)
	GetArticleSearch(ctx context.Context, page int32, search, time string) ([]*Article, int32, error)
	GetArticleImageReview(ctx context.Context, page int32, uuid string) ([]*CreationImageReview, error)
	GetArticleContentReview(ctx context.Context, page int32, uuid string) ([]*CreationContentReview, error)
	ArticleDraftMark(ctx context.Context, id int32, uuid string) error
	SendArticle(ctx context.Context, id int32, uuid, ip string) error
	SendArticleEdit(ctx context.Context, id int32, uuid, ip string) error
	DeleteArticle(ctx context.Context, id int32, uuid string) error
	DeleteArticleDraft(ctx context.Context, id int32, uuid string) error
	SetArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error
	SetArticleView(ctx context.Context, id int32, uuid string) error
	SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error
	CancelArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error
	CancelArticleCollect(ctx context.Context, id int32, uuid, userUuid string) error
	ArticleStatisticJudge(ctx context.Context, id int32, uuid string) (*ArticleStatisticJudge, error)
}

type TalkRepo interface {
	GetTalkList(ctx context.Context, page int32) ([]*Talk, error)
	GetTalkListHot(ctx context.Context, page int32) ([]*Talk, error)
	GetUserTalkList(ctx context.Context, page int32, uuid string) ([]*Talk, error)
	GetUserTalkListVisitor(ctx context.Context, page int32, uuid string) ([]*Talk, error)
	GetTalkCount(ctx context.Context, uuid string) (int32, error)
	GetTalkCountVisitor(ctx context.Context, uuid string) (int32, error)
	GetTalkListStatistic(ctx context.Context, talkList []*Talk) ([]*TalkStatistic, error)
	GetTalkStatistic(ctx context.Context, id int32, uuid string) (*TalkStatistic, error)
	GetLastTalkDraft(ctx context.Context, uuid string) (*TalkDraft, error)
	GetTalkSearch(ctx context.Context, page int32, search, time string) ([]*Talk, int32, error)
	GetUserTalkAgree(ctx context.Context, uuid string) (map[int32]bool, error)
	GetUserTalkCollect(ctx context.Context, uuid string) (map[int32]bool, error)
	GetTalkImageReview(ctx context.Context, page int32, uuid string) ([]*CreationImageReview, error)
	GetTalkContentReview(ctx context.Context, page int32, uuid string) ([]*CreationContentReview, error)
	CreateTalkDraft(ctx context.Context, uuid string) (int32, error)
	SendTalk(ctx context.Context, id int32, uuid, ip string) error
	SendTalkEdit(ctx context.Context, id int32, uuid, ip string) error
	DeleteTalk(ctx context.Context, id int32, uuid string) error
	SetTalkAgree(ctx context.Context, id int32, uuid, userUuid string) error
	CancelTalkAgree(ctx context.Context, id int32, uuid, userUuid string) error
	CancelTalkCollect(ctx context.Context, id int32, uuid, userUuid string) error
	SetTalkView(ctx context.Context, id int32, uuid string) error
	SetTalkCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error
	TalkStatisticJudge(ctx context.Context, id int32, uuid string) (*TalkStatisticJudge, error)
}

type ColumnRepo interface {
	GetLastColumnDraft(ctx context.Context, uuid string) (*ColumnDraft, error)
	GetColumnList(ctx context.Context, page int32) ([]*Column, error)
	GetColumnListHot(ctx context.Context, page int32) ([]*Column, error)
	GetUserColumnList(ctx context.Context, page int32, uuid string) ([]*Column, error)
	GetUserColumnListVisitor(ctx context.Context, page int32, uuid string) ([]*Column, error)
	GetColumnCount(ctx context.Context, uuid string) (int32, error)
	GetColumnCountVisitor(ctx context.Context, uuid string) (int32, error)
	GetColumnListStatistic(ctx context.Context, columnList []*Column) ([]*ColumnStatistic, error)
	GetColumnStatistic(ctx context.Context, id int32, uuid string) (*ColumnStatistic, error)
	GetSubscribeList(ctx context.Context, page int32, uuid string) ([]*Column, error)
	GetSubscribeListCount(ctx context.Context, uuid string) (int32, error)
	GetColumnSubscribes(ctx context.Context, uuid string, ids []int32) ([]*Subscribe, error)
	GetColumnSearch(ctx context.Context, page int32, search, time string) ([]*Column, int32, error)
	GetUserColumnAgree(ctx context.Context, uuid string) (map[int32]bool, error)
	GetUserColumnCollect(ctx context.Context, uuid string) (map[int32]bool, error)
	GetUserSubscribeColumn(ctx context.Context, uuid string) (map[int32]bool, error)
	GetColumnImageReview(ctx context.Context, page int32, uuid string) ([]*CreationImageReview, error)
	GetColumnContentReview(ctx context.Context, page int32, uuid string) ([]*CreationContentReview, error)
	SendColumn(ctx context.Context, id int32, uuid, ip string) error
	SendColumnEdit(ctx context.Context, id int32, uuid, ip string) error
	CreateColumnDraft(ctx context.Context, uuid string) (int32, error)
	SubscribeColumn(ctx context.Context, id int32, uuid string) error
	SubscribeJudge(ctx context.Context, id int32, uuid string) (bool, error)
	DeleteColumn(ctx context.Context, id int32, uuid string) error
	ColumnStatisticJudge(ctx context.Context, id int32, uuid string) (*ColumnStatisticJudge, error)
	SetColumnView(ctx context.Context, id int32, uuid string) error
	SetColumnAgree(ctx context.Context, id int32, uuid, userUuid string) error
	SetColumnCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error
	CancelColumnAgree(ctx context.Context, id int32, uuid, userUuid string) error
	CancelColumnCollect(ctx context.Context, id int32, uuid, userUuid string) error
	CancelSubscribeColumn(ctx context.Context, id int32, uuid string) error
	AddColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error
	DeleteColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error
}

type NewsRepo interface {
	GetNews(ctx context.Context, page int32) ([]*News, error)
	GetNewsSearch(ctx context.Context, page int32, search, time string) ([]*News, int32, error)
}

type CreationUseCase struct {
	repo        CreationRepo
	articleRepo ArticleRepo
	columnRepo  ColumnRepo
	talkRepo    TalkRepo
	re          Recovery
	log         *log.Helper
}

type ArticleUseCase struct {
	repo ArticleRepo
	re   Recovery
	log  *log.Helper
}

type TalkUseCase struct {
	repo TalkRepo
	re   Recovery
	log  *log.Helper
}

type ColumnUseCase struct {
	repo ColumnRepo
	re   Recovery
	log  *log.Helper
}

type NewsUseCase struct {
	repo NewsRepo
	re   Recovery
	log  *log.Helper
}

func NewCreationUseCase(repo CreationRepo, articleRepo ArticleRepo, columnRepo ColumnRepo, talkRepo TalkRepo, re Recovery, logger log.Logger) *CreationUseCase {
	return &CreationUseCase{
		repo:        repo,
		articleRepo: articleRepo,
		columnRepo:  columnRepo,
		talkRepo:    talkRepo,
		re:          re,
		log:         log.NewHelper(log.With(logger, "module", "bff/biz/CreationUseCase")),
	}
}

func NewArticleUseCase(repo ArticleRepo, re Recovery, logger log.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		repo: repo,
		re:   re,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/ArticleUseCase")),
	}
}

func NewTalkUseCase(repo TalkRepo, re Recovery, logger log.Logger) *TalkUseCase {
	return &TalkUseCase{
		repo: repo,
		re:   re,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/TalkUseCase")),
	}
}

func NewColumnUseCase(repo ColumnRepo, re Recovery, logger log.Logger) *ColumnUseCase {
	return &ColumnUseCase{
		repo: repo,
		re:   re,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/ColumnUseCase")),
	}
}

func NewNewsUseCase(repo NewsRepo, re Recovery, logger log.Logger) *NewsUseCase {
	return &NewsUseCase{
		repo: repo,
		re:   re,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/NewsUseCase")),
	}
}

func (r *CreationUseCase) GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error) {
	return r.repo.GetLeaderBoard(ctx)
}

func (r *CreationUseCase) GetCollectArticleList(ctx context.Context, id, page int32) ([]*Article, error) {
	articleList, err := r.repo.GetCollectArticleList(ctx, id, page)
	if err != nil {
		return nil, err
	}
	articleListStatistic, err := r.articleRepo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	return articleList, nil
}

func (r *CreationUseCase) GetCollectArticleCount(ctx context.Context, id int32) (int32, error) {
	return r.repo.GetCollectArticleCount(ctx, id)
}

func (r *CreationUseCase) GetCollectTalkList(ctx context.Context, id, page int32) ([]*Talk, error) {
	talkList, err := r.repo.GetCollectTalkList(ctx, id, page)
	if err != nil {
		return nil, err
	}
	talkListStatistic, err := r.talkRepo.GetTalkListStatistic(ctx, talkList)
	if err != nil {
		return nil, err
	}
	for _, item := range talkListStatistic {
		for index, listItem := range talkList {
			if listItem.Id == item.Id {
				talkList[index].Agree = item.Agree
				talkList[index].View = item.View
				talkList[index].Collect = item.Collect
				talkList[index].Comment = item.Comment
			}
		}
	}
	return talkList, nil
}

func (r *CreationUseCase) GetCollectTalkCount(ctx context.Context, id int32) (int32, error) {
	return r.repo.GetCollectTalkCount(ctx, id)
}

func (r *CreationUseCase) GetCollectColumnList(ctx context.Context, id, page int32) ([]*Column, error) {
	columnList, err := r.repo.GetCollectColumnList(ctx, id, page)
	if err != nil {
		return nil, err
	}
	columnListStatistic, err := r.columnRepo.GetColumnListStatistic(ctx, columnList)
	if err != nil {
		return nil, err
	}
	for _, item := range columnListStatistic {
		for index, listItem := range columnList {
			if listItem.Id == item.Id {
				columnList[index].Agree = item.Agree
				columnList[index].View = item.View
				columnList[index].Collect = item.Collect
			}
		}
	}
	return columnList, nil
}

func (r *CreationUseCase) GetCollectColumnCount(ctx context.Context, id int32) (int32, error) {
	return r.repo.GetCollectColumnCount(ctx, id)
}

func (r *CreationUseCase) GetLastCollectionsDraft(ctx context.Context) (*CollectionsDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetLastCollectionsDraft(ctx, uuid)
}

func (r *CreationUseCase) GetCollectionsContentReview(ctx context.Context, page int32) ([]*CreationContentReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetCollectionsContentReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}

func (r *CreationUseCase) GetUserTimeLineListVisitor(ctx context.Context, page int32, uuid string) ([]*TimeLine, error) {
	timeline, err := r.repo.GetUserTimeLineListVisitor(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	articleList := make([]*Article, 0, len(timeline))
	talkList := make([]*Talk, 0, len(timeline))
	columnList := make([]*Column, 0, len(timeline))
	for _, item := range timeline {
		switch item.Mode {
		case 1:
			articleList = append(articleList, &Article{Id: item.CreationId})
		case 2:
			columnList = append(columnList, &Column{Id: item.CreationId})
		case 3:
			talkList = append(talkList, &Talk{Id: item.CreationId})
		}
	}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		articleListStatistic, err := r.articleRepo.GetArticleListStatistic(ctx, articleList)
		if err != nil {
			return err
		}
		for _, item := range articleListStatistic {
			for index, listItem := range timeline {
				if listItem.CreationId == item.Id {
					timeline[index].Agree = item.Agree
					timeline[index].View = item.View
					timeline[index].Collect = item.Collect
					timeline[index].Comment = item.Comment
				}
			}
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		columnListStatistic, err := r.columnRepo.GetColumnListStatistic(ctx, columnList)
		if err != nil {
			return err
		}
		for _, item := range columnListStatistic {
			for index, listItem := range timeline {
				if listItem.CreationId == item.Id {
					timeline[index].Agree = item.Agree
					timeline[index].View = item.View
					timeline[index].Collect = item.Collect
				}
			}
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		talkListStatistic, err := r.talkRepo.GetTalkListStatistic(ctx, talkList)
		if err != nil {
			return err
		}
		for _, item := range talkListStatistic {
			for index, listItem := range timeline {
				if listItem.CreationId == item.Id {
					timeline[index].Agree = item.Agree
					timeline[index].View = item.View
					timeline[index].Collect = item.Collect
					timeline[index].Comment = item.Comment
				}
			}
		}
		return nil
	}))
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return timeline, nil
}

func (r *CreationUseCase) GetCollections(ctx context.Context, id int32, uuid string) (*Collections, error) {
	return r.repo.GetCollections(ctx, id, uuid)
}

func (r *CreationUseCase) GetCollectionsList(ctx context.Context, page int32) ([]*Collections, error) {
	uuid := ctx.Value("uuid").(string)
	collectionsList, err := r.repo.GetCollectionsList(ctx, uuid, page)
	if err != nil {
		return nil, err
	}
	return collectionsList, nil
}

func (r *CreationUseCase) GetCollectionsListAll(ctx context.Context) ([]*Collections, error) {
	uuid := ctx.Value("uuid").(string)
	collectionsList, err := r.repo.GetCollectionsListAll(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return collectionsList, nil
}

func (r *CreationUseCase) GetCollectionsCount(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetCollectionsCount(ctx, uuid)
}

func (r *CreationUseCase) GetCollectionsListByVisitor(ctx context.Context, page int32, uuid string) ([]*Collections, error) {
	collectionsList, err := r.repo.GetCollectionsListByVisitor(ctx, uuid, page)
	if err != nil {
		return nil, err
	}
	return collectionsList, nil
}

func (r *CreationUseCase) GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetCollectionsVisitorCount(ctx, uuid)
}

func (r *CreationUseCase) CreateCollectionsDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateCollectionsDraft(ctx, uuid)
}

func (r *CreationUseCase) GetCreationUser(ctx context.Context, uuid string) (*CreationUser, error) {
	return r.repo.GetCreationUser(ctx, uuid)
}

func (r *CreationUseCase) GetCreationUserVisitor(ctx context.Context, uuid string) (*CreationUser, error) {
	return r.repo.GetCreationUserVisitor(ctx, uuid)
}

func (r *CreationUseCase) SendCollections(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendCollections(ctx, id, uuid, ip)
}

func (r *CreationUseCase) SendCollectionsEdit(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendCollectionsEdit(ctx, id, uuid, ip)
}

func (r *CreationUseCase) DeleteCollections(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.DeleteCollections(ctx, id, uuid)
}

func (r *ArticleUseCase) GetArticleList(ctx context.Context, page int32) ([]*Article, error) {
	articleList, err := r.repo.GetArticleList(ctx, page)
	if err != nil {
		return nil, err
	}
	articleListStatistic, err := r.repo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetArticleListHot(ctx context.Context, page int32) ([]*Article, error) {
	articleList, err := r.repo.GetArticleListHot(ctx, page)
	if err != nil {
		return nil, err
	}
	articleListStatistic, err := r.repo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetColumnArticleList(ctx context.Context, id int32) ([]*Article, error) {
	articleList, err := r.repo.GetColumnArticleList(ctx, id)
	if err != nil {
		return nil, err
	}
	articleListStatistic, err := r.repo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetArticleCount(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetArticleCount(ctx, uuid)
}

func (r *ArticleUseCase) GetArticleCountVisitor(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetArticleCountVisitor(ctx, uuid)
}

func (r *ArticleUseCase) GetUserArticleList(ctx context.Context, page int32) ([]*Article, error) {
	uuid := ctx.Value("uuid").(string)
	articleList, err := r.repo.GetUserArticleList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	articleListStatistic, err := r.repo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetUserArticleListSimple(ctx context.Context, page int32) ([]*Article, error) {
	uuid := ctx.Value("uuid").(string)
	articleList, err := r.repo.GetUserArticleList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	articleListStatistic, err := r.repo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	if len(articleList) > 2 {
		return articleList[:2], nil
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*Article, error) {
	articleList, err := r.repo.GetUserArticleListVisitor(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	articleListStatistic, err := r.repo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetUserArticleListAll(ctx context.Context) ([]*Article, error) {
	uuid := ctx.Value("uuid").(string)
	articleList, err := r.repo.GetUserArticleListAll(ctx, uuid)
	if err != nil {
		return nil, err
	}
	articleListStatistic, err := r.repo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	return articleList, nil
}

func (r *ArticleUseCase) GetArticleStatistic(ctx context.Context, id int32, uuid string) (*ArticleStatistic, error) {
	return r.repo.GetArticleStatistic(ctx, id, uuid)
}

func (r *ArticleUseCase) GetUserArticleAgree(ctx context.Context) (map[int32]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserArticleAgree(ctx, uuid)
}

func (r *ArticleUseCase) GetUserArticleCollect(ctx context.Context) (map[int32]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserArticleCollect(ctx, uuid)
}

// GetArticleListStatistic todo: delete
func (r *ArticleUseCase) GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error) {
	return nil, nil
}

func (r *ArticleUseCase) GetLastArticleDraft(ctx context.Context) (*ArticleDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetLastArticleDraft(ctx, uuid)
}

func (r *ArticleUseCase) GetArticleImageReview(ctx context.Context, page int32) ([]*CreationImageReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetArticleImageReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}

func (r *ArticleUseCase) GetArticleContentReview(ctx context.Context, page int32) ([]*CreationContentReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetArticleContentReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}

func (r *ArticleUseCase) CreateArticleDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateArticleDraft(ctx, uuid)
}

func (r *ArticleUseCase) ArticleDraftMark(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.ArticleDraftMark(ctx, id, uuid)
}

func (r *ArticleUseCase) GetArticleDraftList(ctx context.Context) ([]*ArticleDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetArticleDraftList(ctx, uuid)
}

func (r *ArticleUseCase) GetArticleSearch(ctx context.Context, page int32, search, time string) ([]*Article, int32, error) {
	articleList, total, err := r.repo.GetArticleSearch(ctx, page, search, time)
	if err != nil {
		return nil, 0, err
	}
	articleListStatistic, err := r.repo.GetArticleListStatistic(ctx, articleList)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range articleListStatistic {
		for index, listItem := range articleList {
			if listItem.Id == item.Id {
				articleList[index].Agree = item.Agree
				articleList[index].View = item.View
				articleList[index].Collect = item.Collect
				articleList[index].Comment = item.Comment
			}
		}
	}
	return articleList, total, nil
}

func (r *ArticleUseCase) SendArticle(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendArticle(ctx, id, uuid, ip)
}

func (r *ArticleUseCase) SendArticleEdit(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendArticleEdit(ctx, id, uuid, ip)
}

func (r *ArticleUseCase) DeleteArticle(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.DeleteArticle(ctx, id, uuid)
}

func (r *ArticleUseCase) DeleteArticleDraft(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.DeleteArticleDraft(ctx, id, uuid)
}

func (r *ArticleUseCase) SetArticleAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetArticleAgree(ctx, id, uuid, userUuid)
}

func (r *ArticleUseCase) SetArticleView(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetArticleView(ctx, id, uuid)
}

func (r *ArticleUseCase) SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetArticleCollect(ctx, id, collectionsId, uuid, userUuid)
}

func (r *ArticleUseCase) CancelArticleAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelArticleAgree(ctx, id, uuid, userUuid)
}

func (r *ArticleUseCase) CancelArticleCollect(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelArticleCollect(ctx, id, uuid, userUuid)
}

func (r *ArticleUseCase) ArticleStatisticJudge(ctx context.Context, id int32) (*ArticleStatisticJudge, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.ArticleStatisticJudge(ctx, id, uuid)
}

func (r *TalkUseCase) GetTalkList(ctx context.Context, page int32) ([]*Talk, error) {
	talkList, err := r.repo.GetTalkList(ctx, page)
	if err != nil {
		return nil, err
	}
	talkListStatistic, err := r.repo.GetTalkListStatistic(ctx, talkList)
	if err != nil {
		return nil, err
	}
	for _, item := range talkListStatistic {
		for index, listItem := range talkList {
			if listItem.Id == item.Id {
				talkList[index].Agree = item.Agree
				talkList[index].View = item.View
				talkList[index].Collect = item.Collect
				talkList[index].Comment = item.Comment
			}
		}
	}
	return talkList, nil
}

func (r *TalkUseCase) GetTalkListHot(ctx context.Context, page int32) ([]*Talk, error) {
	talkList, err := r.repo.GetTalkListHot(ctx, page)
	if err != nil {
		return nil, err
	}
	talkListStatistic, err := r.repo.GetTalkListStatistic(ctx, talkList)
	if err != nil {
		return nil, err
	}
	for _, item := range talkListStatistic {
		for index, listItem := range talkList {
			if listItem.Id == item.Id {
				talkList[index].Agree = item.Agree
				talkList[index].View = item.View
				talkList[index].Collect = item.Collect
				talkList[index].Comment = item.Comment
			}
		}
	}
	return talkList, nil
}

func (r *TalkUseCase) GetUserTalkList(ctx context.Context, page int32) ([]*Talk, error) {
	uuid := ctx.Value("uuid").(string)
	talkList, err := r.repo.GetUserTalkList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	talkListStatistic, err := r.repo.GetTalkListStatistic(ctx, talkList)
	if err != nil {
		return nil, err
	}
	for _, item := range talkListStatistic {
		for index, listItem := range talkList {
			if listItem.Id == item.Id {
				talkList[index].Agree = item.Agree
				talkList[index].View = item.View
				talkList[index].Collect = item.Collect
				talkList[index].Comment = item.Comment
			}
		}
	}
	return talkList, nil
}

func (r *TalkUseCase) GetUserTalkListSimple(ctx context.Context, page int32) ([]*Talk, error) {
	uuid := ctx.Value("uuid").(string)
	talkList, err := r.repo.GetUserTalkList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	talkListStatistic, err := r.repo.GetTalkListStatistic(ctx, talkList)
	if err != nil {
		return nil, err
	}
	for _, item := range talkListStatistic {
		for index, listItem := range talkList {
			if listItem.Id == item.Id {
				talkList[index].Agree = item.Agree
				talkList[index].View = item.View
				talkList[index].Collect = item.Collect
				talkList[index].Comment = item.Comment
			}
		}
	}
	if len(talkList) > 2 {
		return talkList[:2], nil
	}
	return talkList, nil
}

func (r *TalkUseCase) GetUserTalkListVisitor(ctx context.Context, page int32, uuid string) ([]*Talk, error) {
	talkList, err := r.repo.GetUserTalkListVisitor(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	talkListStatistic, err := r.repo.GetTalkListStatistic(ctx, talkList)
	if err != nil {
		return nil, err
	}
	for _, item := range talkListStatistic {
		for index, listItem := range talkList {
			if listItem.Id == item.Id {
				talkList[index].Agree = item.Agree
				talkList[index].View = item.View
				talkList[index].Collect = item.Collect
				talkList[index].Comment = item.Comment
			}
		}
	}
	return talkList, nil
}

func (r *TalkUseCase) GetTalkCount(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetTalkCount(ctx, uuid)
}

func (r *TalkUseCase) GetTalkCountVisitor(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetTalkCountVisitor(ctx, uuid)
}

// GetTalkListStatistic todo: delete
func (r *TalkUseCase) GetTalkListStatistic(ctx context.Context, ids []int32) ([]*TalkStatistic, error) {
	//return r.repo.GetTalkListStatistic(ctx, ids)
	return nil, nil
}

func (r *TalkUseCase) GetTalkStatistic(ctx context.Context, id int32, uuid string) (*TalkStatistic, error) {
	return r.repo.GetTalkStatistic(ctx, id, uuid)
}

func (r *TalkUseCase) GetLastTalkDraft(ctx context.Context) (*TalkDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetLastTalkDraft(ctx, uuid)
}

func (r *TalkUseCase) GetTalkSearch(ctx context.Context, page int32, search, time string) ([]*Talk, int32, error) {
	talkList, total, err := r.repo.GetTalkSearch(ctx, page, search, time)
	if err != nil {
		return nil, 0, err
	}
	talkListStatistic, err := r.repo.GetTalkListStatistic(ctx, talkList)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range talkListStatistic {
		for index, listItem := range talkList {
			if listItem.Id == item.Id {
				talkList[index].Agree = item.Agree
				talkList[index].View = item.View
				talkList[index].Collect = item.Collect
				talkList[index].Comment = item.Comment
			}
		}
	}
	return talkList, total, nil
}

func (r *TalkUseCase) GetUserTalkAgree(ctx context.Context) (map[int32]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserTalkAgree(ctx, uuid)
}

func (r *TalkUseCase) GetUserTalkCollect(ctx context.Context) (map[int32]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserTalkCollect(ctx, uuid)
}

func (r *TalkUseCase) GetTalkImageReview(ctx context.Context, page int32) ([]*CreationImageReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetTalkImageReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}

func (r *TalkUseCase) GetTalkContentReview(ctx context.Context, page int32) ([]*CreationContentReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetTalkContentReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}

func (r *TalkUseCase) CreateTalkDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateTalkDraft(ctx, uuid)
}

func (r *TalkUseCase) SendTalk(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendTalk(ctx, id, uuid, ip)
}

func (r *TalkUseCase) SendTalkEdit(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendTalkEdit(ctx, id, uuid, ip)
}

func (r *TalkUseCase) DeleteTalk(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.DeleteTalk(ctx, id, uuid)
}

func (r *TalkUseCase) SetTalkAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetTalkAgree(ctx, id, uuid, userUuid)
}

func (r *TalkUseCase) CancelTalkAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelTalkAgree(ctx, id, uuid, userUuid)
}

func (r *TalkUseCase) CancelTalkCollect(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelTalkCollect(ctx, id, uuid, userUuid)
}

func (r *TalkUseCase) SetTalkView(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetTalkView(ctx, id, uuid)
}

func (r *TalkUseCase) SetTalkCollect(ctx context.Context, id, collectionsId int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetTalkCollect(ctx, id, collectionsId, uuid, userUuid)
}

func (r *TalkUseCase) TalkStatisticJudge(ctx context.Context, id int32) (*TalkStatisticJudge, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.TalkStatisticJudge(ctx, id, uuid)
}

func (r *ColumnUseCase) GetLastColumnDraft(ctx context.Context) (*ColumnDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetLastColumnDraft(ctx, uuid)
}

func (r *ColumnUseCase) CreateColumnDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateColumnDraft(ctx, uuid)
}

func (r *ColumnUseCase) SubscribeColumn(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SubscribeColumn(ctx, id, uuid)
}

func (r *ColumnUseCase) CancelSubscribeColumn(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CancelSubscribeColumn(ctx, id, uuid)
}

func (r *ColumnUseCase) SubscribeJudge(ctx context.Context, id int32) (bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SubscribeJudge(ctx, id, uuid)
}

func (r *ColumnUseCase) GetColumnList(ctx context.Context, page int32) ([]*Column, error) {
	columnList, err := r.repo.GetColumnList(ctx, page)
	if err != nil {
		return nil, err
	}
	columnListStatistic, err := r.repo.GetColumnListStatistic(ctx, columnList)
	if err != nil {
		return nil, err
	}
	for _, item := range columnListStatistic {
		for index, listItem := range columnList {
			if listItem.Id == item.Id {
				columnList[index].Agree = item.Agree
				columnList[index].View = item.View
				columnList[index].Collect = item.Collect
			}
		}
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetColumnListHot(ctx context.Context, page int32) ([]*Column, error) {
	columnList, err := r.repo.GetColumnListHot(ctx, page)
	if err != nil {
		return nil, err
	}
	columnListStatistic, err := r.repo.GetColumnListStatistic(ctx, columnList)
	if err != nil {
		return nil, err
	}
	for _, item := range columnListStatistic {
		for index, listItem := range columnList {
			if listItem.Id == item.Id {
				columnList[index].Agree = item.Agree
				columnList[index].View = item.View
				columnList[index].Collect = item.Collect
			}
		}
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetUserColumnList(ctx context.Context, page int32) ([]*Column, error) {
	uuid := ctx.Value("uuid").(string)
	columnList, err := r.repo.GetUserColumnList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	columnListStatistic, err := r.repo.GetColumnListStatistic(ctx, columnList)
	if err != nil {
		return nil, err
	}
	for _, item := range columnListStatistic {
		for index, listItem := range columnList {
			if listItem.Id == item.Id {
				columnList[index].Agree = item.Agree
				columnList[index].View = item.View
				columnList[index].Collect = item.Collect
			}
		}
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetUserColumnListSimple(ctx context.Context, page int32) ([]*Column, error) {
	uuid := ctx.Value("uuid").(string)
	columnList, err := r.repo.GetUserColumnList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	columnListStatistic, err := r.repo.GetColumnListStatistic(ctx, columnList)
	if err != nil {
		return nil, err
	}
	for _, item := range columnListStatistic {
		for index, listItem := range columnList {
			if listItem.Id == item.Id {
				columnList[index].Agree = item.Agree
				columnList[index].View = item.View
				columnList[index].Collect = item.Collect
			}
		}
	}
	if len(columnList) > 2 {
		return columnList[:2], nil
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetUserColumnListVisitor(ctx context.Context, page int32, uuid string) ([]*Column, error) {
	columnList, err := r.repo.GetUserColumnListVisitor(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	columnListStatistic, err := r.repo.GetColumnListStatistic(ctx, columnList)
	if err != nil {
		return nil, err
	}
	for _, item := range columnListStatistic {
		for index, listItem := range columnList {
			if listItem.Id == item.Id {
				columnList[index].Agree = item.Agree
				columnList[index].View = item.View
				columnList[index].Collect = item.Collect
			}
		}
	}
	return columnList, nil
}

func (r *ColumnUseCase) GetColumnCount(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetColumnCount(ctx, uuid)
}

func (r *ColumnUseCase) GetColumnCountVisitor(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetColumnCountVisitor(ctx, uuid)
}

// GetColumnListStatistic todo: delete
func (r *ColumnUseCase) GetColumnListStatistic(ctx context.Context, ids []int32) ([]*ColumnStatistic, error) {
	//return r.repo.GetColumnListStatistic(ctx, ids)
	return nil, nil
}

func (r *ColumnUseCase) ColumnStatisticJudge(ctx context.Context, id int32) (*ColumnStatisticJudge, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.ColumnStatisticJudge(ctx, id, uuid)
}

func (r *ColumnUseCase) SendColumn(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendColumn(ctx, id, uuid, ip)
}

func (r *ColumnUseCase) SendColumnEdit(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendColumnEdit(ctx, id, uuid, ip)
}

func (r *ColumnUseCase) DeleteColumn(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.DeleteColumn(ctx, id, uuid)
}

func (r *ColumnUseCase) GetColumnStatistic(ctx context.Context, id int32, uuid string) (*ColumnStatistic, error) {
	return r.repo.GetColumnStatistic(ctx, id, uuid)
}

func (r *ColumnUseCase) GetSubscribeList(ctx context.Context, page int32) ([]*Column, error) {
	uuid := ctx.Value("uuid").(string)
	subscribeList, err := r.repo.GetSubscribeList(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	columnListStatistic, err := r.repo.GetColumnListStatistic(ctx, subscribeList)
	if err != nil {
		return nil, err
	}
	for _, item := range columnListStatistic {
		for index, listItem := range subscribeList {
			if listItem.Id == item.Id {
				subscribeList[index].Agree = item.Agree
				subscribeList[index].View = item.View
				subscribeList[index].Collect = item.Collect
			}
		}
	}
	return subscribeList, nil
}

func (r *ColumnUseCase) GetSubscribeListCount(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetSubscribeListCount(ctx, uuid)
}

func (r *ColumnUseCase) GetColumnSubscribes(ctx context.Context, ids []int32) ([]*Subscribe, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetColumnSubscribes(ctx, uuid, ids)
}

func (r *ColumnUseCase) GetColumnSearch(ctx context.Context, page int32, search, time string) ([]*Column, int32, error) {
	columnList, total, err := r.repo.GetColumnSearch(ctx, page, search, time)
	if err != nil {
		return nil, 0, err
	}
	columnListStatistic, err := r.repo.GetColumnListStatistic(ctx, columnList)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range columnListStatistic {
		for index, listItem := range columnList {
			if listItem.Id == item.Id {
				columnList[index].Agree = item.Agree
				columnList[index].View = item.View
				columnList[index].Collect = item.Collect
			}
		}
	}
	return columnList, total, nil
}

func (r *ColumnUseCase) GetUserColumnAgree(ctx context.Context) (map[int32]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserColumnAgree(ctx, uuid)
}

func (r *ColumnUseCase) GetUserColumnCollect(ctx context.Context) (map[int32]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserColumnCollect(ctx, uuid)
}

func (r *ColumnUseCase) GetUserSubscribeColumn(ctx context.Context) (map[int32]bool, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserSubscribeColumn(ctx, uuid)
}

func (r *ColumnUseCase) GetColumnImageReview(ctx context.Context, page int32) ([]*CreationImageReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetColumnImageReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}

func (r *ColumnUseCase) GetColumnContentReview(ctx context.Context, page int32) ([]*CreationContentReview, error) {
	uuid := ctx.Value("uuid").(string)
	reviewList, err := r.repo.GetColumnContentReview(ctx, page, uuid)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}

func (r *ColumnUseCase) SetColumnAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetColumnAgree(ctx, id, uuid, userUuid)
}

func (r *ColumnUseCase) CancelColumnAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelColumnAgree(ctx, id, uuid, userUuid)
}

func (r *ColumnUseCase) SetColumnCollect(ctx context.Context, id, collectionsId int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetColumnCollect(ctx, id, collectionsId, uuid, userUuid)
}

func (r *ColumnUseCase) CancelColumnCollect(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelColumnCollect(ctx, id, uuid, userUuid)
}

func (r *ColumnUseCase) SetColumnView(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetColumnView(ctx, id, uuid)
}

func (r *ColumnUseCase) AddColumnIncludes(ctx context.Context, id, articleId int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.AddColumnIncludes(ctx, id, articleId, uuid)
}

func (r *ColumnUseCase) DeleteColumnIncludes(ctx context.Context, id, articleId int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.DeleteColumnIncludes(ctx, id, articleId, uuid)
}

func (r *NewsUseCase) GetNews(ctx context.Context, page int32) ([]*News, error) {
	return r.repo.GetNews(ctx, page)
}

func (r *NewsUseCase) GetNewsSearch(ctx context.Context, page int32, search, time string) ([]*News, int32, error) {
	newsList, total, err := r.repo.GetNewsSearch(ctx, page, search, time)
	if err != nil {
		return nil, 0, err
	}
	return newsList, total, nil
}
