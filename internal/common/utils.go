package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	ModeServer = "server"
	ModeClient = "client"
	KnockCount = 3
	TimeWindow = 30
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

func GenerateSecret() string {
	timestamp := time.Now().String()
	hash := sha256.Sum256([]byte(timestamp))
	return hex.EncodeToString(hash[:])
}

func GeneratePorts(secret string, timestamp int64) []int {
	timeWindow := (timestamp / 30) * 30
	data := fmt.Sprintf("%s%d", secret, timeWindow)
	hash := sha256.Sum256([]byte(data))

	ports := make([]int, KnockCount)
	for i := 0; i < KnockCount; i++ {
		portNum := (int(hash[i]) << 8) | int(hash[i+1])
		ports[i] = 10000 + (portNum % 55535)
	}

	InfoLogger.Printf("Timestamp aligned: %d", timeWindow)
	InfoLogger.Printf("Generated ports: %v", ports)

	return ports
}

func GetAlignedTimestamp() int64 {
	now := time.Now()
	seconds := now.Second()
	if seconds >= 30 {
		now = now.Add(time.Duration(30-seconds) * time.Second)
	} else {
		now = now.Add(time.Duration(-seconds) * time.Second)
	}
	return now.Unix()
}
