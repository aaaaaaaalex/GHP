package main

import (
  log "github.com/sirupsen/logrus"

  "net/http"
  "text/template"
)

/*
responseBodyRenderer, implements: http.ResponseWriter
  buffers written bytes for rendering as a go template.
  a Model exposes additional functions/data to the template.
*/
type responseBodyRenderer struct {
  body []byte
}

// Write the body to a buffer, then attempt to render it as a template
func (r *responseBodyRenderer) Render(w http.ResponseWriter, model Model) error {
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

// Write simply cat's bytes to the renderers body
func (r *responseBodyRenderer) Write(b []bytes) (int, error) {
  if r.body == nil || len(r.body) == 0 {
    r.body = b
  } else {
    r.body = append(r.body, b)
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


/*
  RenderedBody decorates a Handler with a responseBodyRenderer middleware.
    model: a ghp.Model for rendering a template against.
*/
func RenderedBody(next http.Handler, model Model) http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      renderer := responseBodyRenderer{
        responseWriter: w,
        request: r,
      }
      next.ServeHTTP(&renderer, r)
      _ = renderer.Render(w, model)
    }
  )
}

