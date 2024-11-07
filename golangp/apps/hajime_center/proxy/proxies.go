package proxy

import (
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/models"
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

	// Retrieve user from context
	user := r.Context().Value("user").(*models.User)

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ModifyResponse = func(resp *http.Response) error {
		return ModifyResponse(resp, r, *user)
	}

	// Serve the request using the proxy
	proxy.ServeHTTP(w, r)
}
