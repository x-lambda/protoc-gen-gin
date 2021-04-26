package timeout

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout 超时控制
// 超时时间由调用方控制
// 默认返回超时错误原因
// 参考：https://github.com/vearne/gin-timeout/blob/master/timeout.go
func Timeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := getBuff() // sync.Pool
		tw := &Writer{
			body:           buffer,
			ResponseWriter: c.Writer,
			h:              make(http.Header),
		}
		c.Writer = tw
		defer func() {
			c.Writer = tw.ResponseWriter
		}()

		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()

		// c.Request.Context()
		c.Request = c.Request.WithContext(ctx)

		finish := make(chan struct{}, 1)
		go func() {
			//
			c.Next()

			finish <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.timeout = true
			tw.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
			tw.ResponseWriter.Write([]byte(fmt.Sprint(ctx.Err())))
			// TODO
			c.Abort()
		case <-finish:
			tw.mu.Lock()
			defer tw.mu.Unlock()

			//if !tw.wroteHeader {
			//	tw.code = http.StatusOK
			//}
			//
			//tw.ResponseWriter.WriteHeader(tw.code)
			//tw.ResponseWriter.Write(buffer.Bytes())
			putBuff(buffer)
		}
	}
}
