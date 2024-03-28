package templates

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fbuedding/iota-admin/web/templates/components"
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

func HandleError(ctx context.Context, w http.ResponseWriter, err error, status int) {
	w.Header().Set("hx-Reswap", "beforeend")
	w.WriteHeader(status)
	err1 := components.Error(err).Render(ctx, w)
	if err1 != nil {
		log.Err(err).Msg("Could not redner error alert")
		return
	}
}
