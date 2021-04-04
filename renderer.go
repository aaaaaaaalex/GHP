package main

import (
  log "github.com/sirupsen/logrus"

  "io"
  "net/http"
  "strings"
  "text/template"
)


const(
  mebibyte int = 1<<20
  defaultMaxBytes int = 2*mebibyte
)

/*
responseBodyRenderer, implements: http.ResponseWriter
  this middleware buffers bytes written for rendering as a go template.
  A Model exposes additional functions/data to the template.
*/
type responseBodyRenderer struct {
  body []byte
  maxBytes int
  model interface{}
  responseWriter http.ResponseWriter
}

// shouldRender decides if it's a good idea to render the renderer's contents
func (r *responseBodyRenderer) shouldRender() bool {
  contentType := r.Header().Get("Content-Type")
  log.Debugf("response content-type: %s", contentType)
  return strings.Contains(contentType, "/html") || len(contentType) == 0
}

// write out any remaining bytes in renderer body and empty it
func (r *responseBodyRenderer) flush() {
  if (r.body != nil){
    log.Debug("flushing renderer body without executing...")
    r.responseWriter.Write(r.body)
  }
}

/*
Write will buffer bytes for later rendering up to a maximum of `maxBytes`.
If body size exceeds `maxBytes`, excess bytes are written out (FIFO)
  ~Mem complx: O(maxBytes + len(b))
*/
func (r *responseBodyRenderer) Write(b []byte) (int, error) {
  if r.body == nil {
    r.body = make([]byte, 0, r.maxBytes)
  }
  log.Debugf("Appending %d bytes to renderer body...", len(b))
  r.body = append(r.body, b...) // concatenate body

  bodyLen := len(r.body)
  if overLimit := bodyLen - r.maxBytes; overLimit > 0 {
    log.Warn("renderer body size limit exceeded: response may be only partially rendered.")
    r.responseWriter.Write( r.body[0 : overLimit] ) // flush excess
    r.body = r.body[overLimit : bodyLen] // preserve tail
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
execute r's contents as a template, writing out to w.
  note: if r's maxBytes has been exceeded, only the last maxBytes worth of data will be executed.
*/
func execute(w io.Writer, r responseBodyRenderer) error {
  log.Debugf("Executing %d bytes...", len(r.body))
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

  return nil
}

/*
  RenderedBody decorates a Handler with a responseBodyRenderer middleware.
  Requests to render are identified based on the Content-Type set by 'next's response
    next:	the handler to wrap
    model:	the datastructure to expose to templates.
*/
func RenderedBody(next http.Handler, model interface{}) http.Handler {
  return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
      renderer := responseBodyRenderer{
	      model: model,
	      responseWriter: w,
	      maxBytes: defaultMaxBytes,
      }
      next.ServeHTTP(&renderer, r)

      if renderer.shouldRender() {
	_ = execute(w, renderer)
      } else {
        log.Info("Serving unrendered response based on non-html content header")
        renderer.flush()
      }
  })
}

