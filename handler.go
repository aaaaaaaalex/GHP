package main

import (
  log "github.com/sirupsen/logrus"

  "net/http"
  "text/template"
)

/*
responseBodyRenderer, implements: http.ResponseWriter
  intercepts bytes to a wrapped writer, renders as a go template.
  a model is exposed to the template to aid in dynamic rendering.
*/
type responseBodyRenderer struct {
  responseWriter http.ResponseWriter
  model interface{}
  body []byte
}

// Write the body to a buffer, then attempt to render it as a template
func (r *responseBodyRenderer) Write (body []byte) (int, error){
  r.body = body

  template, err := template.New("bodyTemplate").Parse(string(r.body))
  if err != nil {
    log.Errorf("Error parsing response body as template: %s", err.Error())
    return 0, err
  }

  // count the number of bytes written to the responsewriter with a counting buffer
  countedBuffer := CountBytes( r.responseWriter )
  err = template.Execute(countedBuffer, r.model)
  if err != nil {
    log.Errorf("Error executing response template with model: %s", err.Error())
  }

  return countedBuffer.BytesWritten(), err
}

// Header does not alter the wrapped writers functionality
func (r *responseBodyRenderer) Header() http.Header {
  return r.responseWriter.Header()
}

// WriterHeader does not alter the wrapped writers functionality
func (r *responseBodyRenderer) WriteHeader(statusCode int){
  r.responseWriter.WriteHeader(statusCode)
}


// rendered builds a middleware for rendering a response body as a Go template
func rendered(next http.Handler) http.Handler {
  return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
    renderer := responseBodyRenderer{
       responseWriter: w,
       model: r,
    }
    next.ServeHTTP(&renderer, r)
  })
}

