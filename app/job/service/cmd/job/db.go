package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

type DB struct {
	client *gorm.DB
}

func (db *DB) setNews(m map[string]string) error {
	user := &News{
		Id: m["id"],
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	err := db.client.WithContext(ctx).Select("Id").Create(user).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil
		} else {
			return errors.Wrapf(err, fmt.Sprintf("fail to insert a new: %s", "db-setNews"))
		}
	}
	return nil
}

func newDB(config map[string]interface{}) (*DB, error) {
	source := config["db"].(map[string]interface{})["source"].(string)
	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("failed opening connection to db: %s", "newDB"))
	}
	return &DB{
		client: db,
	}, nil
}
