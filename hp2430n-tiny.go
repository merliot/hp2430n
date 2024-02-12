//go:build tinygo

package hp2430n

import "github.com/merliot/device/uart"

func newTransport() uart.Uart {
	return uart.New()
}
