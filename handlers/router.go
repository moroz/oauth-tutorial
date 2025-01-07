package handlers

import "net/http"

func Router() http.Handler {
	r := http.NewServeMux()

	pages := PageController()
	r.HandleFunc("GET /", pages.Index)

	return r
}
