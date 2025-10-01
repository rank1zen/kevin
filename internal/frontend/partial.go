package frontend

import (
	"net/http"

	"github.com/a-h/templ"
)

type Partial struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

type PartialFunc func() *Partial

func ResolvePartial() templ.Attributes {
	attrs := templ.Attributes{
		"hx-post":    "/profile.Header",
		"hx-trigger": "load",
	}

	return attrs
}
