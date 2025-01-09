package handlers

import "net/http"

func Router() http.Handler {
	r := http.NewServeMux()

	pages := PageController()
	r.HandleFunc("GET /", pages.Index)

	oauth := OAuthController()
	r.HandleFunc("GET /oauth/github/start", oauth.Start)

	return r
}
