package template

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/rs/zerolog/log"
)

func Prepare(r *http.Request, component templ.Component) templ.Component {

	if r.Header.Get("HX-Request") == "true" {
		log.Debug().Msg("HTMX Request")
		return component
	}
	return Main(component)
}
