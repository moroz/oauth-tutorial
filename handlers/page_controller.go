package handlers

import (
	"log"
	"net/http"

	"github.com/moroz/oauth-tutorial/templates"
)

type pageController struct{}

func PageController() pageController {
	return pageController{}
}

type pageIndexAssigns struct {
	*Context
}

func (c *pageController) Index(w http.ResponseWriter, r *http.Request) {
	err := templates.PagesIndex.Execute(w, pageIndexAssigns{
		ContextFromRequest(r),
	})
	if err != nil {
		log.Println(err)
	}
}
