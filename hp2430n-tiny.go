//go:build tinygo

package hp2430n

import (
	"machine"
)

type targetStruct struct {
	uart *machine.UART
}

func (h *Hp2430n) targetNew() {
	h.uart = machine.UART0
	h.uart.Configure(machine.UARTConfig{
		TX:       machine.UART0_TX_PIN,
		RX:       machine.UART0_RX_PIN,
		BaudRate: 9600,
	})
}

func (h *Hp2430n) Write(buf []byte) (n int, err error) {
	return h.uart.Write(buf)
}

func (h *Hp2430n) Read(buf []byte) (n int, err error) {
	return h.uart.Read(buf)
}
