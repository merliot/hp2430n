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
	return make([]byte, words*2), nil
	/*
		switch start {
		case regMaxVoltage:
				0:        14,
				17*2 - 1: 0,
			}, nil
		case regBatteryCapacity:
			return []byte{
				0: 14,
				9: 0,
			}, nil
		case regLoadInfo:
			return []byte{
				0:  14,
				17: 0,
			}, nil
		case regLoadCmd:
			return []byte{
				0:  14,
				17: 0,
			}, nil
		case regOpDays:
			return []byte{
				0:  14,
				17: 0,
			}, nil
		case regAlarm:
			return []byte{
				0:  14,
				17: 0,
			}, nil
		}
		return nil, fmt.Errorf("unknown start")
	*/
}
