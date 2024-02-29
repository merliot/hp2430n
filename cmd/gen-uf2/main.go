package main

import (
	"log"
	"os"

	"github.com/merliot/hp2430n"
)

func main() {
	hp2430n := hp2430n.New("proto", "hp2430n", "proto").(*hp2430n.Hp2430n)
	if err := hp2430n.GenerateUf2s(); err != nil {
		log.Println("Error generating UF2s:", err)
		os.Exit(1)
	}
}
