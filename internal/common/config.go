package common

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Mode      string
	Host      string
	KeyFile   string
	OnSuccess string
}

func ParseFlags() *Config {
	config := &Config{}

	// Define flags with short and long versions
	flag.StringVar(&config.Mode, "m", "", "Mode of operation (server/client)")
	flag.StringVar(&config.Mode, "mode", "", "Mode of operation (server/client)")
	flag.StringVar(&config.Host, "h", "localhost", "Target host (default: localhost)")
	flag.StringVar(&config.Host, "host", "localhost", "Target host (default: localhost)")
	flag.StringVar(&config.KeyFile, "f", "", "Path to the key file.\nIf file not exists we will create on the specified path (only server)")
	flag.StringVar(&config.KeyFile, "file", "", "Path to the key file.\nIf file not exists we will create on the specified path (only server)")
	flag.StringVar(&config.OnSuccess, "on-success", "", "Command to execute on successful port knocking (server mode)")

	// Customize the Usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: program [options]

Options:
  -m, --mode          Mode of operation (server/client) [required]
  -h, --host          Target host (default: localhost)
  -f, --file          Path to the key file [optional for server, required for client]
                      If file not exists we will create on the specified path (only server)
  --on-success        Command to execute on successful port knocking (server mode)
`)
	}

	// Parse the flags
	flag.Parse()

	// Validate the required `mode` flag
	if config.Mode == "" {
		fmt.Fprintln(os.Stderr, "Error: -m/--mode is required")
		flag.Usage()
		os.Exit(1)
	}

	if config.KeyFile != "" && config.Mode == "server" {
		if _, err := os.Stat(config.KeyFile); err != nil {
			if os.IsNotExist(err) {
				secret := GenerateSecret()

				err := os.WriteFile(config.KeyFile, []byte(secret), 0600)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating key file: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	// If --file is not specified and the mode is server, create a new key file
	if config.KeyFile == "" && config.Mode == "server" {
		defaultFile := ".key" // Default file name
		fmt.Println("No key file specified. Generating a new key file:", defaultFile)

		// Generate a new secret
		secret := GenerateSecret()

		// Write the secret to the file
		err := os.WriteFile(defaultFile, []byte(secret), 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating key file: %v\n", err)
			os.Exit(1)
		}

		config.KeyFile = defaultFile
	}

	// For client mode, ensure --file is always provided
	if config.KeyFile == "" && config.Mode == "client" {
		fmt.Fprintln(os.Stderr, "Error: --file flag is required for client mode")
		flag.Usage()
		os.Exit(1)
	}

	return config
}
