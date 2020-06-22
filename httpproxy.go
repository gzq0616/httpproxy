package httpproxy

import (
	"github.com/gobwas/glob"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const StartChar = "^"
const GlobChar = "*"

/*
target: 'http://www.example.org', // target host
pathRewrite: {
	'^/api/old-path': '/api/new-path', // rewrite path
	'^/api/remove/path': '/path',      // remove base path
}
*/
type Options struct {
	Target string // target host. eg. http://www.example.org
	//ChangeOrigin bool   // needed for virtual hosted sites
	PathRewrite map[string]string
}

func HandleProxy(opt Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remote, _ := url.Parse(opt.Target)
		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.Host = remote.Host

		for k, v := range opt.PathRewrite {
			matched := matchPath(k, r.RequestURI)
			if matched {
				rp := rewritePath(k, v, r.RequestURI)
				r.URL.Path = rp
				break
			}
		}

		proxy.ServeHTTP(w, r)
	}
}

func rewritePath(pattern, rewritePath, reqPath string) string {
	if strings.HasPrefix(pattern, StartChar) {
		pattern = strings.Split(pattern, StartChar)[1]
	}
	return strings.Replace(reqPath, pattern, rewritePath, 1)
}

func matchPath(pattern, reqPath string) bool {
	if strings.HasPrefix(pattern, StartChar) {
		pattern = strings.Split(pattern, StartChar)[1]
	}
	return glob.MustCompile(pattern + GlobChar).Match(reqPath)
}
