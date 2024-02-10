//go:build !tinygo

package hp2430n

import (
	"net/http"
)

func (h *Hp2430n) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.API(w, r, h)
}
