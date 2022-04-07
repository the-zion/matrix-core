package data

import (
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewUserRepo)

// Data .
type Data struct {
	db *gorm.DB
}

func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "module", "user-service/data.mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("failed creat or update table resources: %v", err)
	}
	return db
}

func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "user-service/data"))

	d := &Data{
		db: db,
	}
	return d, func() {
		l.Info("closing the data resources")
	}, nil
}
