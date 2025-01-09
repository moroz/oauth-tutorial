package github

import "net/url"

const INIT_BASE_URL = "https://github.com/login/oauth/authorize"

func BuildOAuthInitURL(clientID, redirectURL, state string) string {
	query := url.Values{
		"client_id":    {clientID},
		"redirect_uri": {redirectURL},
		"state":        {state},
		"scope":        {"user"},
	}
	return INIT_BASE_URL + "?" + query.Encode()
}
