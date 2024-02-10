//go:build !tinygo

package hp2430n

import (
	"net/http"
)

func (h *Hp2430n) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.API(w, r, h)
}

func (h *Hp2430n) DescHtml() []byte {
	desc, _ := fs.ReadFile("html/desc.html")
	return desc
}

func (h *Hp2430n) SupportedTargets() string {
	return h.Targets.FullNames()
}
