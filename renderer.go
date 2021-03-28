package main

import (
  log "github.com/sirupsen/logrus"

  "io"
  "net/http"
  "text/template"
)


/*
responseBodyRenderer, implements: http.ResponseWriter
  this middleware buffers bytes written for rendering as a go template.
  A Model exposes additional functions/data to the template.
*/
type responseBodyRenderer struct {
  body []byte
  model interface{}
  responseWriter http.ResponseWriter
}

// Write simply cat's bytes to the renderers body
func (r *responseBodyRenderer) Write(b []byte) (int, error) {
  if r.body == nil || len(r.body) == 0 {
    r.body = b
  } else {
    r.body = append(r.body, b...)
  }

  return len(b), nil
}

// Header does not alter the wrapped-writers functionality
func (r *responseBodyRenderer) Header() http.Header {
  return r.responseWriter.Header()
}

// WriterHeader does not alter the wrapped-writers functionality
func (r *responseBodyRenderer) WriteHeader(statusCode int){
  r.responseWriter.WriteHeader(statusCode)
}


// execute the content written to r as a template, writing out to w
func execute(w io.Writer, r responseBodyRenderer) error {
  template, err := template.New("bodyTemplate").Parse(string(r.body))
  if err != nil {
    log.Errorf("Error parsing response body as template: %s", err.Error())
    return err
  }

  err = template.Execute(w, r.model)
  if err != nil {
    log.Errorf("Error executing response template with model: %s", err.Error())
    return err
  }

  return nil;
}

/*
  RenderedBody decorates a Handler with a responseBodyRenderer middleware.
    model: a ghp.Model for rendering a template against.
*/
func RenderedBody(next http.Handler, model interface{}) http.Handler {
  return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
      renderer := responseBodyRenderer{
	      model: model,
	      responseWriter: w,
      }
      next.ServeHTTP(&renderer, r)
      _ = execute(w, renderer)
  })
}

