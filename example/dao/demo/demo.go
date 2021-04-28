package demo

import (
	"context"
	"time"

	"github.com/x-lambda/protoc-gen-gin-example/util/db"
)

type Item struct {
	ID         int32
	Name       string
	CreateTime time.Time
	ModifyTime time.Time
}

func QueryByID(ctx context.Context, id int32) (item Item, err error) {
	conn := db.Get(ctx, "")
	sql := "select id, name, create_time, modify_time from t_demo where id = ?"
	q := db.SQLSelect("t_demo", sql)
	err = conn.QueryRowContext(ctx, q, id).Scan(&item.ID, &item.Name, &item.CreateTime, &item.ModifyTime)
	return
}
