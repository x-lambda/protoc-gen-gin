package demo

import (
	"context"
	"fmt"
	"time"
)

func TestTimeout(ctx context.Context) {
	time.Sleep(1 * time.Second)
	fmt.Println("你看不到我😛")
}
