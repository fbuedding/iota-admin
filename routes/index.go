package routes

import (
	"net/http"

	"github.com/fbuedding/iota-admin/internal/globals"
	iotagentsdk "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/fbuedding/iota-admin/web/templates"
	"github.com/fbuedding/iota-admin/web/templates/pages"
	"github.com/go-chi/chi/v5"
)

func Index() chi.Router {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	})

	router.Get("/index", func(w http.ResponseWriter, r *http.Request) {
		templates.Prepare(r, pages.Index([]iotagentsdk.IoTA{{Host: globals.Conf.IoTAHost, Port: globals.Conf.IoTAPort}})).Render(r.Context(), w)
		// TODO add multiple IoT-Agent support
	})

	return router
}
