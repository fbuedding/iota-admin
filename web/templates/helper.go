package templates

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/fbuedding/iota-admin/web/templates/layouts"
	"github.com/rs/zerolog/log"
)

func Prepare(r *http.Request, component templ.Component) templ.Component {

	if r.Header.Get("HX-Request") == "true" {
		log.Debug().Msg("HTMX Request")
		return component
	}
	return layouts.Main(component)
}
