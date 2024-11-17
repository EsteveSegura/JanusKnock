package client

import (
	"fmt"
	"janusknock/internal/common"
	"net"
	"os"
	"time"
)

func Run(config *common.Config) {
	// Read the secret from the specified file
	keyBytes, err := os.ReadFile(config.KeyFile)
	if err != nil {
		common.ErrorLogger.Printf("Error reading key file: %v\n", err)
		return
	}

	secret := string(keyBytes) // Use a local variable for the secret
	common.InfoLogger.Printf("Read secret from file: %s\n", secret)

	alignedTime := common.GetAlignedTimestamp()
	ports := common.GeneratePorts(secret, alignedTime)

	common.InfoLogger.Printf("Attempting port sequence: %v\n", ports)

	for _, port := range ports {
		knock(config.Host, port)
		time.Sleep(time.Second)
	}
}

func knock(host string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		common.InfoLogger.Printf("Knock on port %d\n", port)
		return
	}
	conn.Close()
}
