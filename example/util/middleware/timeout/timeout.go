package timeout

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// timeoutWriter implements http.Writer
// 可以参考标准库 http.timeoutHandler
type timeoutWriter struct {
	gin.ResponseWriter

	h    http.Header
	wbuf bytes.Buffer

	mu          sync.Mutex
	timeout     bool
	wroteHeader bool
	code        int
}

// Write the response
func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timeout {
		return 0, nil
	}

	return tw.wbuf.Write(b)
}

// WriteHeader In http.ResponseWriter interface
func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)

	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timeout || tw.wroteHeader {
		return
	}

	tw.writeHeader(code)
}

// writeHeader set that the header has been written
func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

// Header "relays" the header, h, set in struct
// In http.ResponseWriter interface
func (tw *timeoutWriter) Header() http.Header {
	return tw.h
}

// SetTimeout sets timedOut field to true
func (tw *timeoutWriter) SetTimeout() {
	tw.timeout = true
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid write header code: %v", code))
	}
}

// Timeout 超时控制
// 超时时间由调用方控制，默认返回 timeout
// 搬自: https://github.com/JacobSNGoodwin/memrizr/blob/master/account/handler/middleware/timeout.go
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// set Gin's writer as our custom writer
		tw := &timeoutWriter{
			ResponseWriter: c.Writer,
			h:              make(http.Header),
		}
		c.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// update gin request context
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		finished := make(chan struct{})        // to indicate handler finished
		panicChan := make(chan interface{}, 1) // used to handle panics if we can't recover

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next() // calls subsequent middleware(s) and handler
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			// if we cannot recover from panic,
			// send internal server error
			e := NewInternal()
			tw.ResponseWriter.WriteHeader(e.Status())
			eResp, _ := json.Marshal(gin.H{
				"code": -1,
				"msg":  e.Error(),
			})
			tw.ResponseWriter.Write(eResp)
		case <-finished:
			// if finished, set headers and write resp
			tw.mu.Lock()
			defer tw.mu.Unlock()
			// map Headers from tw.Header() (written to by gin)
			// to tw.ResponseWriter for response
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			// tw.wbuf will have been written to already when gin writes to tw.Write()
			tw.ResponseWriter.Write(tw.wbuf.Bytes())
		case <-ctx.Done():
			// timeout has occurred, send errTimeout and write headers
			tw.mu.Lock()
			defer tw.mu.Unlock()
			et := NewRequestTimeout("timeout")
			// ResponseWriter from gin
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(et.Status())
			eResp, _ := json.Marshal(gin.H{
				"code": -1,
				"msg":  et.Error(),
			})
			tw.ResponseWriter.Write(eResp)
			c.Abort()
			tw.SetTimeout()
		}
	}
}
