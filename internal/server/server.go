package server

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"janusknock/internal/common"
)

type Server struct {
	config    *common.Config
	attempts  map[string][]int
	mutex     sync.RWMutex
	listeners []net.Listener
}

func Run(config *common.Config) {
	if config.OnSuccess == "" {
		common.ErrorLogger.Println("Error: --on-success flag is required in server mode")
		os.Exit(1)
	}

	server := &Server{
		config:   config,
		attempts: make(map[string][]int),
	}

	// Read the secret from the specified file
	keyBytes, err := os.ReadFile(config.KeyFile)
	if err != nil {
		common.ErrorLogger.Printf("Error reading key file: %v\n", err)
		return
	}
	secret := string(keyBytes) // Use a local variable for the secret
	common.InfoLogger.Printf("Secret read from file: %s\n", secret)

	server.start(secret)
}

func (s *Server) start(secret string) {
	common.InfoLogger.Println("Server initialized and listening...")
	common.InfoLogger.Println("Waiting for port knocking...")

	alignedTime := common.GetAlignedTimestamp()
	currentPorts := common.GeneratePorts(secret, alignedTime)
	common.InfoLogger.Printf("Initial active ports: %v\n", currentPorts)

	s.setupListeners(currentPorts, secret)

	for {
		now := time.Now()
		var waitTime time.Duration
		if now.Second() >= 30 {
			waitTime = time.Duration(60-now.Second()) * time.Second
		} else {
			waitTime = time.Duration(30-now.Second()) * time.Second
		}

		waitTime -= time.Duration(now.Nanosecond())
		time.Sleep(waitTime)

		s.closeListeners()

		alignedTime = common.GetAlignedTimestamp()
		currentPorts = common.GeneratePorts(secret, alignedTime)
		common.InfoLogger.Printf("New active ports: %v\n", currentPorts)

		s.setupListeners(currentPorts, secret)
	}
}

func (s *Server) setupListeners(ports []int, secret string) {
	for _, port := range ports {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			common.ErrorLogger.Printf("Error when listening port %d: %v\n", port, err)
			continue
		}
		s.listeners = append(s.listeners, listener)
		go s.handleListener(listener, port, secret)
	}

	nextChange := time.Now().Add(time.Duration(30-time.Now().Second()%30) * time.Second)
	common.InfoLogger.Printf("Next port replacement: %02d:%02d:%02d\n",
		nextChange.Hour(), nextChange.Minute(), nextChange.Second())
}

func (s *Server) closeListeners() {
	for _, listener := range s.listeners {
		listener.Close()
	}
	s.listeners = nil
}

func (s *Server) handleListener(listener net.Listener, port int, secret string) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				common.ErrorLogger.Printf("Error while accepting connection: %v\n", err)
			}
			return
		}
		go s.handleConnection(conn, port, secret)
	}
}

func (s *Server) handleConnection(conn net.Conn, port int, secret string) {
	defer conn.Close()

	ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()

	s.mutex.Lock()
	s.attempts[ip] = append(s.attempts[ip], port)
	attemptCount := len(s.attempts[ip])
	common.InfoLogger.Printf("Knock received from %s at port %d (attempt %d/%d)\n",
		ip, port, attemptCount, common.KnockCount)
	s.mutex.Unlock()

	if attemptCount == common.KnockCount {
		s.verifySequence(ip, secret)
	}
}

func (s *Server) verifySequence(ip string, secret string) {
	currentPorts := common.GeneratePorts(secret, common.GetAlignedTimestamp())

	s.mutex.Lock()
	defer s.mutex.Unlock()

	attempts := s.attempts[ip]
	delete(s.attempts, ip)

	if len(attempts) == len(currentPorts) {
		match := true
		for i := range attempts {
			if attempts[i] != currentPorts[i] {
				match = false
				break
			}
		}

		if match {
			common.InfoLogger.Printf("Correct port sequence from %s\n", ip)
			s.executeOnSuccess(ip)
		} else {
			common.InfoLogger.Printf("Wrong port sequence from %s\n", ip)
			common.InfoLogger.Printf("   Received: %v\n", attempts)
			common.InfoLogger.Printf("   Expected: %v\n", currentPorts)
		}
	}
}

func (s *Server) executeOnSuccess(ip string) {
	common.InfoLogger.Printf("Executing success command for %s: %s\n", ip, s.config.OnSuccess)

	// Create the command
	cmd := exec.Command("bash", "-c", s.config.OnSuccess)
	cmd.Env = append(os.Environ(), fmt.Sprintf("REMOTE_IP=%s", ip))

	// Capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		common.ErrorLogger.Printf("Error executing success command: %v\n", err)
		return
	}

	common.InfoLogger.Printf("Success command output: %s\n", string(output))
}
