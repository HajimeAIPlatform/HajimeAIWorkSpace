package proxy

import (
	"HajimeAIWorkSpace/common/apps/hajime_center/initializers"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// DifyHandler forwards requests after removing the "dify" prefix
func DifyHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the "dify" prefix from the URL path
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/dify")

	// Parse the target URL
	config, _ := initializers.LoadEnv(".")
	targetURL, err := url.Parse(config.DifyHost)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Serve the request using the proxy
	proxy.ServeHTTP(w, r)
}
