# how to used

## 使用http server方式
```go
package main

import (
	"github.com/gzq0616/httpproxy"
	"net/http"
)

func main() {
	srv := &http.Server{
		Addr: ":8080",
		Handler: httpproxy.HandleProxy(httpproxy.Options{
			Target: "http://example.xxx.com",
			PathRewrite: map[string]string{
				"^/proxy/": "/",
			},
		}),
	}

	srv.ListenAndServe()
}

```

## 使用gin框架
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gzq0616/httpproxy"
)

func proxyMiddleware() gin.HandlerFunc {
	proxy := httpproxy.HandleProxy(httpproxy.Options{
		Target: "http://example.xxx.com",
		PathRewrite: map[string]string{
			"^/proxy-middle/": "/",
		},
	})

	return func(c *gin.Context) {
		proxy(c.Writer, c.Request)
		return
	}
}

func handle(c *gin.Context) {
	proxy := httpproxy.HandleProxy(httpproxy.Options{
		Target: "http://example.xxx.com",
		PathRewrite: map[string]string{
			"^/proxy-handle/": "/",
		},
	})
	proxy(c.Writer, c.Request)
	c.Next()
}

func main() {
	g := gin.Default()
	
	// 使用中间件方式
	g.Any("/proxy-middle/*action", proxyMiddleware())

	// 使用http handlerFunc方式
	g.Any("/proxy/*action", gin.WrapF(httpproxy.HandleProxy(httpproxy.Options{
		Target: "http://example.xxx.com",
		PathRewrite: map[string]string{
			"^/proxy/": "/",
		},
	})))

	// 直接使用handler方式
	g.Any("/proxy-handle/*action", handle)

	g.Run(":8080")
}
```

