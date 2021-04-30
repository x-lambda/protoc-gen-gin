package demo

import (
	"context"
	"fmt"
	"time"

	"github.com/x-lambda/protoc-gen-gin/example/dao/demo"
)

func TestTimeout(ctx context.Context) (err error) {
	time.Sleep(1 * time.Second)
	// fmt.Println("你看不到我😛")
	item, err := demo.QueryByID(ctx, 11)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)
	return
}
