package timeout

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Writer struct {
	gin.ResponseWriter

	body *bytes.Buffer
	h    http.Header

	mu          sync.Mutex
	timeout     bool
	wroteHeader bool
	code        int
}

func (w *Writer) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.timeout {
		return 0, nil
	}

	return w.body.Write(b)
}

func (w *Writer) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.timeout {
		return
	}

	w.writeHeader(code)
}

func (w *Writer) writeHeader(code int) {
	w.wroteHeader = true
	w.code = code
}

func (w *Writer) WriteHeaderNow() {}

func (w *Writer) Header() http.Header {
	return w.h
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid write header code: %v", code))
	}
}
