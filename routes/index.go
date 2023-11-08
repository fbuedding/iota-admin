package routes

import (
	"net/http"

	iotagentsdk "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/fbuedding/iota-admin/web/template"
	"github.com/go-chi/chi/v5"
)

func Index() chi.Router {
	r := chi.NewRouter()
  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/index", http.StatusSeeOther)
  })
  r.Get("/index", func(w http.ResponseWriter, r *http.Request) {
    template.Index([]iotagentsdk.IoTA{iotagentsdk.IoTA{Host: "iot-agent", Port: 4061 }}).Render(r.Context(), w)
    //http.ServeFile(w, r, "web/static/index.html")

	})
  
	return r
}

