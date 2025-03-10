package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	contentType := w.Header().Get("Content-Type")
	if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
		return w.Writer.Write(b)
	}

	return w.ResponseWriter.Write(b)
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request", http.StatusBadRequest)
				return
			}
			defer gzipReader.Close()
			r.Body = io.NopCloser(gzipReader)
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		gzResponseWriter := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gzipWriter,
		}

		w.Header().Set("Content-Encoding", "gzip")

		h.ServeHTTP(gzResponseWriter, r)
	})
}
