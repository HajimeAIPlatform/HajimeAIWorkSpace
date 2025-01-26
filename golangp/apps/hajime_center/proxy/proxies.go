package proxy

import (
	"hajime/golangp/apps/hajime_center/constants"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/initializers"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// DifyHandler forwards requests after removing the "dify" prefix
func DifyHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the "dify" prefix from the URL path
	r.URL.Path = strings.TrimPrefix(r.URL.Path, constants.DifyServerPrefix)

	// Parse the target URL
	config, _ := initializers.LoadEnv(".")
	targetURL, err := url.Parse(config.DifyHost)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	user := models.User{}

	// 从上下文中检索用户信息
	if u, ok := r.Context().Value("user").(*models.User); ok && u != nil {
		user = *u
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ModifyResponse = func(resp *http.Response) error {
		return ModifyResponse(resp, r, user)
	}
	// Serve the request using the proxy
	proxy.ServeHTTP(w, r)
}
