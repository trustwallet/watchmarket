package middleware

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/storage"
	"io/ioutil"
	"time"
)

type cachedWriter struct {
	gin.ResponseWriter
	status  int
	written bool
	expire  int64
	key     string
}

var cache storage.Middleware

func InitCachingMiddleware(c storage.Middleware) {
	cache = c
}

func NewCachedWriter(expire int64, writer gin.ResponseWriter, key string) *cachedWriter {
	return &cachedWriter{writer, 0, false, expire, key}
}

func (w *cachedWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachedWriter) Status() int {
	return w.ResponseWriter.Status()
}

func (w *cachedWriter) Written() bool {
	return w.ResponseWriter.Written()
}

func (w *cachedWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err != nil {
		return 0, errors.E(err, "fail to cache write string", errors.Params{"data": data})
	}
	if w.Status() != 200 {
		return 0, errors.E("Write: invalid cache status", errors.Params{"data": data})
	}
	cacheResp := storage.CacheResponse{
		Status: w.Status(),
		Header: w.Header(),
		Data:   data,
	}
	cacheData := storage.CacheData{
		Data:    cacheResp,
		Expired: time.Now().Unix(),
	}

	result, err := cache.Set(w.key, cacheData)
	if err != nil || result != storage.SaveResultSuccess {
		return 0, errors.E("Failed to Set", errors.Params{"data": data, "saving_error": err})
	}

	return ret, nil
}

func (w *cachedWriter) WriteString(data string) (n int, err error) {
	ret, err := w.ResponseWriter.WriteString(data)
	if err != nil {
		return 0, errors.E(err, "fail to cache write string", errors.Params{"data": data})
	}
	if w.Status() != 200 {
		return 0, errors.E("WriteString: invalid cache status", errors.Params{"data": data})
	}
	cacheResp := storage.CacheResponse{
		Status: w.Status(),
		Header: w.Header(),
		Data:   []byte(data),
	}
	cacheData := storage.CacheData{
		Data:    cacheResp,
		Expired: time.Now().Unix(),
	}
	result, err := cache.Set(w.key, cacheData)
	if err != nil || result != storage.SaveResultSuccess {
		return 0, errors.E("Failed to Set", errors.Params{"data": data, "saving_error": err})
	}

	return ret, err
}

func generateKey(c *gin.Context) string {
	url := c.Request.URL.String()
	var b []byte
	if c.Request.Body != nil {
		b, _ = ioutil.ReadAll(c.Request.Body)
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}
	hash := sha1.Sum(append([]byte(url), b...))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// CacheMiddleware encapsulates a gin handler function and caches the response with an expiration time.
func CacheMiddleware(expiration int64, handle gin.HandlerFunc) gin.HandlerFunc {
	if cache == nil {
		logger.Fatal("gin cache middleware is created with empty cache")
	}
	return func(c *gin.Context) {
		defer c.Next()
		key := generateKey(c)
		mc, err := cache.Get(key)
		if err != nil || mc.Data.Data == nil || time.Now().Unix()-mc.Expired > expiration {
			writer := NewCachedWriter(expiration, c.Writer, key)
			c.Writer = writer
			handle(c)

			if c.IsAborted() {
				cache.Delete(key)
			}
			return
		}

		c.Writer.WriteHeader(mc.Data.Status)
		for k, vals := range mc.Data.Header {
			for _, v := range vals {
				c.Writer.Header().Set(k, v)
			}
		}
		_, err = c.Writer.Write(mc.Data.Data)
		if err != nil {
			cache.Delete(key)
			logger.Error(err, "cannot write data", mc)
		}
	}
}
