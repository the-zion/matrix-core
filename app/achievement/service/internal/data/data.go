package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	_ "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/biz"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"runtime"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewTransaction, NewAchievementRepo, NewRecovery)

type Data struct {
	db       *gorm.DB
	log      *log.Helper
	redisCli redis.Cmdable
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

func (d *Data) GroupRecover(ctx context.Context, fn func(ctx context.Context) error) func() error {
	return func() error {
		defer func() {
			if rerr := recover(); rerr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.Context(ctx).Errorf("%v: %s\n", rerr, buf)
			}
		}()
		return fn(ctx)
	}
}

func NewRecovery(d *Data) biz.Recovery {
	return d
}

func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "module", "achievement/data/mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	return db
}

func NewRedis(conf *conf.Data, logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "module", "achievement/data/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		DB:           2,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
		Password:     conf.Redis.Password,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		l.Fatalf("redis connect error: %v", err)
	}
	return client
}

func NewData(db *gorm.DB, redisCmd redis.Cmdable, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "achievement/data/new-data"))

	d := &Data{
		db:       db,
		redisCli: redisCmd,
	}
	return d, func() {
		l.Info("closing the data resources")
	}, nil
}
