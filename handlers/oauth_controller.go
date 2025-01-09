package handlers

import "net/http"

type oauthController struct{}

func OAuthController() oauthController {
	return oauthController{}
}

func (c *oauthController) Start(w http.ResponseWriter, r *http.Request) {
}
