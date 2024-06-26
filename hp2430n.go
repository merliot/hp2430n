package hp2430n

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/merliot/dean"
	"github.com/merliot/device"
	"github.com/merliot/device/modbus"
)

const (
	regMaxVoltage      = 0x000A
	regBatteryCapacity = 0x0100
	regLoadInfo        = 0x0120
	regLoadCmd         = 0x010A
	regOpDays          = 0x0115
	regAlarm           = 0x0121
)

type System struct {
	MaxVolts      uint8
	ChargeAmps    uint8
	DischargeAmps uint8
	ProductType   uint8
	Model         string
	SWVersion     string
	HWVersion     string
	Serial        string
}

type Controller struct {
	Temp   uint8    // deg C
	Alarms []string // faults and warnings
}

func (c Controller) isDiff(other Controller) bool {
	if c.Temp != other.Temp {
		return true
	}
	if len(c.Alarms) != len(other.Alarms) {
		return true
	}
	for i, v := range c.Alarms {
		if v != other.Alarms[i] {
			return true
		}
	}
	return false
}

type Battery struct {
	SOC         uint8
	Volts       float32
	Amps        float32
	Temp        uint8 // deg C
	ChargeState string
}

type LoadInfo struct {
	Volts      float32
	Amps       float32
	Status     bool
	Brightness uint8
}

type Solar struct {
	Volts float32
	Amps  float32
}

type Daily struct {
	BattMinVolts      float32
	BattMaxVolts      float32
	ChargeMaxAmps     float32
	DischargeMaxAmps  float32
	ChargeMaxWatts    uint16
	DischargeMaxWatts uint16
	ChargeAmpHrs      uint16
	DischargeAmpHrs   uint16
	GenPowerWatts     uint16
	ConPowerWatts     uint16
}

type Historical struct {
	OpDays          uint16
	OverDischarges  uint16
	FullCharges     uint16
	ChargeAmpHrs    uint32
	DischargeAmpHrs uint32
	GenPowerWatts   uint32
	ConPowerWatts   uint32
}

type msgStatus struct {
	Path   string
	Status string
}

type msgSystem struct {
	Path   string
	System System
}

type msgController struct {
	Path       string
	Controller Controller
}

type msgBattery struct {
	Path    string
	Battery Battery
}

type msgLoadInfo struct {
	Path     string
	LoadInfo LoadInfo
}

type msgSolar struct {
	Path  string
	Solar Solar
}

type msgDaily struct {
	Path  string
	Daily Daily
}

type msgHistorical struct {
	Path       string
	Historical Historical
}

type Hp2430n struct {
	*device.Device
	*modbus.Modbus `json:"-"`
	Status         string
	System         System
	Controller     Controller
	Battery        Battery
	LoadInfo       LoadInfo
	Solar          Solar
	Daily          Daily
	Historical     Historical
}

var targets = []string{"demo", "rpi", "nano-rp2040"}

func New(id, model, name string) dean.Thinger {
	fmt.Println("NEW HP2430N\r")
	h := &Hp2430n{}
	h.Device = device.New(id, model, name, fs, targets).(*device.Device)
	h.Modbus = modbus.New(newTransport())
	h.Status = "OK"
	return h
}

func (h *Hp2430n) save(msg *dean.Msg) {
	msg.Unmarshal(h).Broadcast()
}

func (h *Hp2430n) getState(msg *dean.Msg) {
	h.Path = "state"
	msg.Marshal(h).Reply()
}

func (h *Hp2430n) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":             h.save,
		"get/state":         h.getState,
		"update/status":     h.save,
		"update/system":     h.save,
		"update/controller": h.save,
		"update/battery":    h.save,
		"update/load":       h.save,
		"update/solar":      h.save,
		"update/dialy":      h.save,
		"update/historical": h.save,
	}
}

func (h *Hp2430n) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.API(w, r, h)
}

func version(b []byte) string {
	return fmt.Sprintf("%02d.%02d.%02d", b[1], b[2], b[3])
}

func serial(b []byte) string {
	return fmt.Sprintf("%02X%02X-%02X%02X", b[0], b[1], b[2], b[3])
}

func swap(b []byte) uint16 {
	return (uint16(b[0]) << 8) | uint16(b[1])
}

func swap4(b []byte) uint32 {
	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3])
}

func volts(b []byte) float32 {
	return float32(swap(b)) * 0.1
}

func unvolts(volts float32) []byte {
	value := uint16(volts / 0.1)
	return []byte{byte(value >> 8), byte(value & 0xFF)}
}

func amps(b []byte) float32 {
	return float32(swap(b)) * 0.01
}

func unamps(amps float32) []byte {
	value := uint16(amps / 0.01)
	return []byte{byte(value >> 8), byte(value & 0xFF)}
}

func chargeState(b byte) string {
	switch b {
	case 0:
		return "Deactivated"
	case 1:
		return "Activated"
	case 2:
		return "Mode MPPT"
	case 3:
		return "Mode Equalizing"
	case 4:
		return "Mode Boost"
	case 5:
		return "Mode Float"
	case 6:
		return "Current Limiting (Overpower)"
	}
	return "Unknown"
}

func (h *Hp2430n) readSystem(s *System) error {
	// System Info (34 bytes)
	regs, err := h.ReadRegisters(regMaxVoltage, 17)
	if err != nil {
		return err
	}
	s.MaxVolts = uint8(regs[0])
	s.ChargeAmps = uint8(regs[1])
	s.DischargeAmps = uint8(regs[2])
	s.ProductType = uint8(regs[3])
	s.Model = strings.ReplaceAll(string(regs[4:20]), "\000", "")
	s.SWVersion = version(regs[20:24])
	s.HWVersion = version(regs[24:28])
	s.Serial = serial(regs[28:32])
	// skip dev addr regs[32:34]
	return nil
}

var alarms = []string{
	"Battery over-discharge",
	"Battery over-voltage",
	"Battery under-voltage",
	"Load short circuit",
	"Load over-power or load over-current",
	"Controller temperature too high",
	"Battery high temperature protection (temperature higher than the upper discharge limit); prohibit charging",
	"Solar input over-power",
	"(reserved)",
	"Solar input side over-voltage",
	"(reserved)",
	"Solar panel working point over-voltage",
	"Solar panel reverse connected",
	"(reserved)",
	"(reserved)",
	"(reserved)",
	"(reserved)",
	"(reserved)",
	"(reserved)",
	"(reserved)",
	"(reserved)",
	"(reserved)",
	"Power main supply",
	"OO battery detected (SLD)",
	"Battery high temperature protection (temperature higher than the upper discharge limit); prohibit discharging",
	"Battery low temperature protection (temperature lower than the lower discharge limit); prohibit discharging",
	"Over-charge protection; stop charging",
	"Battery low temperature protection (temperature is lower than the lower limit of charging; stop charging",
	"Battery reverse connected",
	"Capacitor over-voltage (reserved)",
	"Induction probe damaged (street light)",
	"Load open-circuit (street light)",
}

func parseAlarms(b []byte) (active []string) {
	value := swap4(b)
	for i := 0; i < 32; i++ {
		// Check if the bit is set
		if value&(1<<i) != 0 {
			// Add corresponding alarm
			active = append(active, alarms[i])
		}
	}
	return
}

func (h *Hp2430n) readDynamic(c *Controller, b *Battery, l *LoadInfo, s *Solar) error {

	// Controller Dynamic Info (20 bytes)
	regs, err := h.ReadRegisters(regBatteryCapacity, 10)
	if err != nil {
		return err
	}
	// reserved regs[0]
	b.SOC = uint8(regs[1])
	b.Volts = volts(regs[2:4])
	b.Amps = amps(regs[4:6])
	c.Temp = uint8(regs[6])
	b.Temp = uint8(regs[7])
	l.Volts = volts(regs[8:10])
	l.Amps = amps(regs[10:12])
	// skip load power regs[12:14]
	s.Volts = volts(regs[14:16])
	s.Amps = amps(regs[16:18])
	// skip solar power regs[18:20]

	// Load Information (2 bytes)
	regs, err = h.ReadRegisters(regLoadInfo, 1)
	if err != nil {
		return err
	}
	l.Status = (regs[0] & 0x80) == 0x80
	l.Brightness = uint8(regs[0] & 0x7F)
	b.ChargeState = chargeState(regs[1])

	// Controller alarm information
	regs, err = h.ReadRegisters(regAlarm, 2)
	if err != nil {
		return err
	}
	c.Alarms = parseAlarms(regs)

	return nil
}

func (h *Hp2430n) readDaily(d *Daily) error {

	// Current Day Info (22 bytes)
	regs, err := h.ReadRegisters(regLoadCmd, 11)
	if err != nil {
		return err
	}
	// skip load cmd regs[0:2]
	d.BattMinVolts = volts(regs[2:4])
	d.BattMaxVolts = volts(regs[4:6])
	d.ChargeMaxAmps = amps(regs[6:8])
	d.DischargeMaxAmps = amps(regs[8:10])
	d.ChargeMaxWatts = swap(regs[10:12])
	d.DischargeMaxWatts = swap(regs[12:14])
	d.ChargeAmpHrs = swap(regs[14:16])
	d.DischargeAmpHrs = swap(regs[16:18])
	d.GenPowerWatts = swap(regs[18:20])
	d.ConPowerWatts = swap(regs[20:22])
	return nil
}

func (h *Hp2430n) readHistorical(d *Historical) error {

	// Historical Info (22 bytes)
	regs, err := h.ReadRegisters(regOpDays, 11)
	if err != nil {
		return err
	}
	d.OpDays = swap(regs[0:2])
	d.OverDischarges = swap(regs[2:4])
	d.FullCharges = swap(regs[4:6])
	d.ChargeAmpHrs = swap4(regs[6:10])
	d.DischargeAmpHrs = swap4(regs[10:14])
	d.GenPowerWatts = swap4(regs[14:18])
	d.ConPowerWatts = swap4(regs[18:22])
	return nil
}

func (h *Hp2430n) sendStatus(i *dean.Injector, newStatus string) {
	if h.Status == newStatus {
		return
	}

	var status = msgStatus{Path: "update/status"}
	var msg dean.Msg

	status.Status = newStatus
	i.Inject(msg.Marshal(status))
}

func (h *Hp2430n) sendSystem(i *dean.Injector) {
	var system = msgSystem{Path: "update/system"}
	var msg dean.Msg

	// sendSystem blocks until we get a good system info read

	for {
		if err := h.readSystem(&system.System); err != nil {
			h.sendStatus(i, err.Error())
			continue
		}
		i.Inject(msg.Marshal(system))
		break
	}

	h.sendStatus(i, "OK")
}

func (h *Hp2430n) sendDynamic(i *dean.Injector) {
	var controller = msgController{Path: "update/controller"}
	var battery = msgBattery{Path: "update/battery"}
	var loadInfo = msgLoadInfo{Path: "update/load"}
	var solar = msgSolar{Path: "update/solar"}
	var msg dean.Msg

	err := h.readDynamic(&controller.Controller, &battery.Battery,
		&loadInfo.LoadInfo, &solar.Solar)
	if err != nil {
		h.sendStatus(i, err.Error())
		return
	}

	// If anything has changed, send update msg(s)

	if controller.Controller.isDiff(h.Controller) {
		i.Inject(msg.Marshal(controller))
	}
	if battery.Battery != h.Battery {
		i.Inject(msg.Marshal(battery))
	}
	if loadInfo.LoadInfo != h.LoadInfo {
		i.Inject(msg.Marshal(loadInfo))
	}
	if solar.Solar != h.Solar {
		i.Inject(msg.Marshal(solar))
	}

	h.sendStatus(i, "OK")
}

func (h *Hp2430n) sendHourly(i *dean.Injector) {
	var daily = msgDaily{Path: "update/daily"}
	var historical = msgHistorical{Path: "update/historical"}
	var msg dean.Msg

	err := h.readDaily(&daily.Daily)
	if err != nil {
		h.sendStatus(i, err.Error())
		return
	}
	if daily.Daily != h.Daily {
		i.Inject(msg.Marshal(daily))
	}
	err = h.readHistorical(&historical.Historical)
	if err != nil {
		h.sendStatus(i, err.Error())
		return
	}
	if historical.Historical != h.Historical {
		i.Inject(msg.Marshal(historical))
	}

	h.sendStatus(i, "OK")
}

func (h *Hp2430n) Run(i *dean.Injector) {

	h.sendSystem(i)
	h.sendDynamic(i)
	h.sendHourly(i)

	nextHour := time.Now().Add(time.Hour)
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		h.sendDynamic(i)
		if time.Now().After(nextHour) {
			h.sendHourly(i)
			nextHour = time.Now().Add(time.Hour)
		}
	}
}
