// go run ./cmd
// go run -tags prime ./cmd
// tinygo flash -target xxx ./cmd

package main

import (
	"github.com/merliot/dean"
	"github.com/merliot/device/runner"
	"github.com/merliot/hp2430n"
)

var (
	id           = dean.GetEnv("ID", "hp2430n01")
	name         = dean.GetEnv("NAME", "Solisto hp2430n")
	deployParams = dean.GetEnv("DEPLOY_PARAMS", "")
	wsScheme     = dean.GetEnv("WS_SCHEME", "ws://")
	port         = dean.GetEnv("PORT", "8000")
	portPrime    = dean.GetEnv("PORT_PRIME", "8001")
	user         = dean.GetEnv("USER", "")
	passwd       = dean.GetEnv("PASSWD", "")
	dialURLs     = dean.GetEnv("DIAL_URLS", "")
	ssids        = dean.GetEnv("WIFI_SSIDS", "")
	passphrases  = dean.GetEnv("WIFI_PASSPHRASES", "")
)

func main() {
	hp2430n := hp2430n.New(id, "hp2430n", name).(*hp2430n.Hp2430n)
	hp2430n.SetDeployParams(deployParams)
	hp2430n.SetWifiAuth(ssids, passphrases)
	hp2430n.SetDialURLs(dialURLs)
	hp2430n.SetWsScheme(wsScheme)
	runner.Run(hp2430n.Device, port, portPrime, user, passwd, dialURLs, wsScheme)
}
