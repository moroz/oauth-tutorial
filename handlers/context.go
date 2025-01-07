package handlers

import (
	"context"
	"net/http"
)

type Context struct {
	ctx context.Context
}

func ContextFromRequest(r *http.Request) *Context {
	return &Context{r.Context()}
}

func (c *Context) IsSignedIn() bool {
	return false
}

func (c *Context) UserName() string {
	return "Test User Name"
}
