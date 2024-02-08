import { WebSocketController, ViewMode } from './common.js'

export function run(prefix, url, viewMode) {
	const hp2430n = new Hp2430n(prefix, url, viewMode)
}

class Hp2430n extends WebSocketController {

	open() {
		super.open()

		if (this.state.DeployParams === "") {
			return
		}

		this.show()
	}

	show() {
		this.showStatus()
		this.showSystem()
		this.showController()
		this.showBattery()
		this.showLoadInfo()
		this.showSolar()
		this.showDaily()
		this.showHistorical()
	}

	showStatus() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var status = document.getElementById("status")
			status.value = ""
			status.value += "Status:                      " + this.state.Status
			break;
		}
	}

	showSystem() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("system")
			ta.value = ""
			ta.value += "Max Voltage Supported (V):   " + this.state.System.MaxVolts + "\r\n"
			ta.value += "Rated Charge Current (A):    " + this.state.System.ChargeAmps + "\r\n"
			ta.value += "Rated Discharge Current (A): " + this.state.System.DischargeAmps + "\r\n"
			ta.value += "Product Type:                " + this.state.System.ProductType + "\r\n"
			ta.value += "Model:                       " + this.state.System.Model + "\r\n"
			ta.value += "Software Version:            " + this.state.System.SWVersion + "\r\n"
			ta.value += "Hardware Version:            " + this.state.System.HWVersion + "\r\n"
			ta.value += "Serial:                      " + this.state.System.Serial
			break;
		}
	}

	showAlarms() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var list = ""
			for (let i = 0; i < this.state.Controller.Alarms.length; i++) {
				list += "\r\n    \u26A0 " + this.state.Controller.Alarms[i]
			}
			return list
			break;
		}
	}

	showController() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("controller")
			ta.value = ""
			ta.value += "* Temp (C):                  " + this.state.Controller.Temp + "\r\n"
			ta.value += "* Alarms:                    "
			if (this.state.Controller.Alarms === null) {
				ta.rows = 2
				ta.value += "<none>"
			} else {
				ta.rows = 3 + this.state.Controller.Alarms.length
				ta.value += this.showAlarms()
			}
			break;
		}
	}

	showBattery() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("battery")
			ta.value = ""
			ta.value += "* Capacity SOC:              " + this.state.Battery.SOC + "\r\n"
			ta.value += "* Voltage (V):               " + this.state.Battery.Volts + "\r\n"
			ta.value += "* Current (A):               " + this.state.Battery.Amps + "\r\n"
			ta.value += "* Temp (C):                  " + this.state.Battery.Temp + "\r\n"
			ta.value += "* Charging State:            " + this.state.Battery.ChargeState
			break;
		case ViewMode.ViewTile:
			document.getElementById("battery-volts").innerText = this.state.Battery.Volts.toFixed(2)
			document.getElementById("battery-amps").innerText = this.state.Battery.Amps.toFixed(2)
			break;
		}
	}

	showLoadInfo() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("load")
			ta.value = ""
			ta.value += "* Voltage (V):               " + this.state.LoadInfo.Volts + "\r\n"
			ta.value += "* Current (A):               " + this.state.LoadInfo.Amps + "\r\n"
			ta.value += "* Status:                    " + this.state.LoadInfo.Status + "\r\n"
			ta.value += "* Brightness:                " + this.state.LoadInfo.Brightness
			break;
		case ViewMode.ViewTile:
			document.getElementById("load-volts").innerText = this.state.LoadInfo.Volts.toFixed(2)
			document.getElementById("load-amps").innerText = this.state.LoadInfo.Amps.toFixed(2)
			break;
		}
	}

	showSolar() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("solar")
			ta.value = ""
			ta.value += "* Voltage (V):               " + this.state.Solar.Volts + "\r\n"
			ta.value += "* Current (A):               " + this.state.Solar.Amps
			break;
		case ViewMode.ViewTile:
			document.getElementById("solar-volts").innerText = this.state.Solar.Volts.toFixed(2)
			document.getElementById("solar-amps").innerText = this.state.Solar.Amps.toFixed(2)
			break;
		}
	}

	showDaily() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("daily")
			ta.value = ""
			ta.value += "Battery Min Voltage (V):     " + this.state.Daily.BattMinVolts + "\r\n"
			ta.value += "Battery Max Voltage (V):     " + this.state.Daily.BattMaxVolts + "\r\n"
			ta.value += "Charging Max Current (A):    " + this.state.Daily.ChargeMaxAmps + "\r\n"
			ta.value += "Discharging Max Current (A): " + this.state.Daily.DischargeMaxAmps + "\r\n"
			ta.value += "Charging Max Power (W):      " + this.state.Daily.ChargeMaxWatts + "\r\n"
			ta.value += "Disharging Max Power (W):    " + this.state.Daily.DischargeMaxWatts + "\r\n"
			ta.value += "Charge Current (Ah):         " + this.state.Daily.ChargeAmpHrs + "\r\n"
			ta.value += "Discharge Current (Ah):      " + this.state.Daily.DischargeAmpHrs + "\r\n"
			ta.value += "Power Generated (W):         " + this.state.Daily.GenPowerWatts + "\r\n"
			ta.value += "Power Consumed (W):          " + this.state.Daily.ConPowerWatts
			break;
		}
	}

	showHistorical() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("historical")
			ta.value = ""
			ta.value += "Total Operating Days:        " + this.state.Historical.OpDays + "\r\n"
			ta.value += "Total Batt Over Discharges:  " + this.state.Historical.OverDischarges + "\r\n"
			ta.value += "Total Battery Full Charges:  " + this.state.Historical.FullCharges + "\r\n"
			ta.value += "Total Charge Current (Ah):   " + this.state.Historical.ChargeAmpHrs + "\r\n"
			ta.value += "Total Discharge Cur (Ah):    " + this.state.Historical.DischargeAmpHrs + "\r\n"
			ta.value += "Total Power Generated (W):   " + this.state.Historical.GenPowerWatts + "\r\n"
			ta.value += "Total Power Consumed (W):    " + this.state.Historical.ConPowerWatts
			break;
		}
	}

	handle(msg) {
		switch(msg.Path) {
		case "update/status":
			this.state.Status = msg.Status
			this.showStatus()
			break
		case "update/system":
			this.state.System = msg.System
			this.showSystem()
			break
		case "update/controller":
			this.state.Controller = msg.Controller
			this.showController()
			break
		case "update/battery":
			this.state.Battery = msg.Battery
			this.showBattery()
			break
		case "update/load":
			this.state.LoadInfo = msg.LoadInfo
			this.showLoadInfo()
			break
		case "update/solar":
			this.state.Solar = msg.Solar
			this.showSolar()
			break
		case "update/daily":
			this.state.Daily = msg.Daily
			this.showDaily()
			break
		case "update/historical":
			this.state.Historical = msg.Historical
			this.showHistorical()
			break
		}
	}
}
