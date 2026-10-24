package authn

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type apiContentWriter struct {
	gin.ResponseWriter
}

func (w *apiContentWriter) restrictActiveContent(sample []byte) {
	contentType := w.Header().Get("Content-Type")
	if contentType == "" && len(sample) > 0 {
		contentType = http.DetectContentType(sample)
	}
	mediaType, _, _ := strings.Cut(strings.ToLower(contentType), ";")
	switch strings.TrimSpace(mediaType) {
	case "text/html", "application/xhtml+xml", "image/svg+xml", "text/xml", "application/xml":
		w.Header().Set("Content-Security-Policy", "sandbox; default-src 'none'; frame-ancestors 'none'")
		w.Header().Set("Content-Disposition", "attachment")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}
}

func (w *apiContentWriter) WriteHeader(status int) {
	w.restrictActiveContent(nil)
	w.ResponseWriter.WriteHeader(status)
}

func (w *apiContentWriter) WriteHeaderNow() {
	w.restrictActiveContent(nil)
	w.ResponseWriter.WriteHeaderNow()
}

func (w *apiContentWriter) Write(data []byte) (int, error) {
	w.restrictActiveContent(data[:min(512, len(data))])
	return w.ResponseWriter.Write(data)
}

func (w *apiContentWriter) WriteString(data string) (int, error) {
	w.restrictActiveContent([]byte(data[:min(512, len(data))]))
	return w.ResponseWriter.WriteString(data)
}
