package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/user/service/internal/data"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var (
	source string
)

func NewDB(logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "update", "db"))

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	if err := db.AutoMigrate(
		&data.User{},
		&data.Profile{},
		&data.ProfileUpdate{},
		&data.Follow{},
		&data.AvatarReview{},
	); err != nil {
		l.Fatalf("failed creat or update table resources: %v", err)
	}
	return db
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)
	NewDB(logger)
}

func init() {
	flag.StringVar(&source, "source",
		"root:123456@tcp(127.0.0.1:3306)/core?charset=utf8mb4&parseTime=True&loc=Local",
		"database source, eg: -source source path")
}
