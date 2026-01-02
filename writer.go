package nesgress

import (
	"bytes"
	"io"
	"sync"
)

// synchronizedWriter wraps an io.Writer with a mutex to prevent concurrent writes.
type synchronizedWriter struct {
	writer io.Writer
	mutex  sync.Mutex
}

func (sw *synchronizedWriter) Write(p []byte) (n int, err error) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	return sw.writer.Write(p)
}

// safeBytesBuffer provides thread-safe access to a bytes.Buffer for both reads and writes.
type safeBytesBuffer struct {
	buf   *bytes.Buffer
	mutex sync.RWMutex
}

func (sbb *safeBytesBuffer) Write(p []byte) (n int, err error) {
	sbb.mutex.Lock()
	defer sbb.mutex.Unlock()
	return sbb.buf.Write(p)
}

// SafeString provides thread-safe read access to buffer content
func (sbb *safeBytesBuffer) SafeString() string {
	sbb.mutex.RLock()
	defer sbb.mutex.RUnlock()
	return sbb.buf.String()
}
