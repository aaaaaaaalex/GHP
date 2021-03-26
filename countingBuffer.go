package main

import (
  "net/http"
)

type countingBuffer struct {
  Writer	http.ResponseWriter
  bytesWritten 	int
}

// write to the responseWriter, making note of how many bytes
func (c countingBuffer) Write(b []byte) (int, error) {
  var written int
  written, err := c.Writer.Write(b)
  c.bytesWritten += written

  return written, err
}

// Header satisfies the ResponseWriter interface
func (c countingBuffer) Header() http.Header {
  return c.Writer.Header()
}

// WriteHeader satisfies the ResponseWriter interface
func (c countingBuffer) WriteHeader(statusCode int) {
  c.Writer.WriteHeader(statusCode)
}

// BytesWritten encapsulates the internal count of bytes written
func (c *countingBuffer) BytesWritten() int {
  return c.bytesWritten
}

// countBytes modifies a responseWriter by wrapping it with a counter
func CountBytes(r http.ResponseWriter) countingBuffer {
  return countingBuffer {
    Writer: r,
  }
}
