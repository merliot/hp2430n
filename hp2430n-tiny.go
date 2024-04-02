//go:build tinygo

package hp2430n

import (
	"embed"

	"github.com/merliot/device/uart"
)

var fs embed.FS

func newTransport() *uart.Uart {
	return uart.New()
}
