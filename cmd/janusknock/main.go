package main

import (
	"fmt"
	"os"

	"janusknock/internal/client"
	"janusknock/internal/common"
	"janusknock/internal/server"
)

func main() {
	config := common.ParseFlags()

	switch config.Mode {
	case common.ModeServer:
		server.Run(config)
	case common.ModeClient:
		client.Run(config)
	default:
		fmt.Println("Invalid mode")
		os.Exit(1)
	}
}
