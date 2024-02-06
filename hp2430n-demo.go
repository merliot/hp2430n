//go:build !tinygo && !rpi

package hp2430n

type targetRelayStruct struct {
}

type targetStruct struct {
	osStruct
}

func (h *Hp2430n) targetNew() {
	h.osNew()
}

func (h *Hp2430n) readRegisters(start, words uint16) ([]byte, error) {
	regs := make([]byte, words*2)
	switch start {
	case regMaxVoltage:
	case regBatteryCapacity:
		copy(regs[2:4], unvolts(13.4))   // battery.volts
		copy(regs[4:6], unamps(3.5))     // battery.amps
		copy(regs[8:10], unvolts(12.8))  // load.volts
		copy(regs[10:12], unamps(2.1))   // load.amps
		copy(regs[14:16], unvolts(15.7)) // solar.volts
		copy(regs[16:18], unamps(1.4))   // solar.amps
	case regLoadInfo:
	case regLoadCmd:
	case regOpDays:
	case regAlarm:
	}
	return regs, nil
}
