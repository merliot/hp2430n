//go:build !tinygo

package hp2430n

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/merliot/device"
)

//go:embed css html js template
var fs embed.FS

type osStruct struct {
	templates *template.Template
}

func (h *Hp2430n) osNew() {
	h.CompositeFs.AddFS(fs)
	h.templates = h.CompositeFs.ParseFS("template/*")
}

func (h *Hp2430n) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.TrimPrefix(r.URL.Path, "/") {
	case "state":
		device.ShowState(h.templates, w, h)
	default:
		h.API(h.templates, w, r)
	}
}

func (h *Hp2430n) DescHtml() []byte {
	desc, _ := fs.ReadFile("html/desc.html")
	return desc
}

func (h *Hp2430n) SupportedTargets() string {
	return h.Targets.FullNames()
}
