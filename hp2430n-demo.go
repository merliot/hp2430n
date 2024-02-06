//go:build !tinygo && !rpi

package hp2430n

type targetRelayStruct struct {
}

type targetStruct struct {
	osStruct
	start uint16
	words uint16
}

func (h *Hp2430n) targetNew() {
	h.osNew()
}

func (h *Hp2430n) Write(buf []byte) (n int, err error) {
	// get start and words from Modbus request
	h.start = (uint16(buf[2]) << 8) | uint16(buf[3])
	h.words = (uint16(buf[4]) << 8) | uint16(buf[5])
	return n, nil
}

func (h *Hp2430n) Read(buf []byte) (n int, err error) {
	res := buf[3:]
	switch h.start {
	case regMaxVoltage:
	case regBatteryCapacity:
		copy(res[2:4], unvolts(13.4))   // battery.volts
		copy(res[4:6], unamps(3.5))     // battery.amps
		copy(res[8:10], unvolts(12.8))  // load.volts
		copy(res[10:12], unamps(2.1))   // load.amps
		copy(res[14:16], unvolts(15.7)) // solar.volts
		copy(res[16:18], unamps(1.4))   // solar.amps
	case regLoadInfo:
	case regLoadCmd:
	case regOpDays:
	case regAlarm:
	}
	return int(5 + h.words*2), nil
}
