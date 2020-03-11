package middleware

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/storage"
	"net/http/httptest"
	"testing"
	"time"
)

var Cache *storage.Storage

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestWrite(t *testing.T) {
	s := setupRedis(t)
	Cache = internal.InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	InitCachingMiddleware(Cache)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	writer := NewCachedWriter(60*3, c.Writer, "mykey")
	c.Writer = writer

	c.Writer.WriteHeader(204)
	c.Writer.WriteHeaderNow()
	c.Writer.Write([]byte("foo"))
	assert.Equal(t, 204, c.Writer.Status())
	assert.Equal(t, "foo", w.Body.String())
	assert.True(t, c.Writer.Written())
}

func TestCachePage(t *testing.T) {
	InitCachingMiddleware(Cache)
	router := gin.New()
	router.GET("/cache_ping", GinCachingMiddleware(60*3, func(c *gin.Context) {
		ginutils.RenderSuccess(c, "pong "+fmt.Sprint(time.Now().UnixNano()))
	}))

	w1 := performRequest("GET", "/cache_ping", router)
	w2 := performRequest("GET", "/cache_ping", router)

	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, w1.Body.String(), w2.Body.String())
}

func TestCachePageExpire(t *testing.T) {
	InitCachingMiddleware(Cache)
	router := gin.New()
	router.GET("/cache_ping", GinCachingMiddleware(1, func(c *gin.Context) {
		ginutils.RenderSuccess(c, "pong "+fmt.Sprint(time.Now().UnixNano()))
	}))

	w1 := performRequest("GET", "/cache_ping", router)
	time.Sleep(time.Second * 3)
	w2 := performRequest("GET", "/cache_ping", router)

	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, 200, w2.Code)
	assert.NotEqual(t, w1.Body.String(), w2.Body.String())
}
