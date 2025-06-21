package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	var fanName string
	var help bool
	flag.StringVar(&fanName, "fan", "", "Target specific fan by name")
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.BoolVar(&help, "h", false, "Show help message (shorthand)")
	flag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		printHelp()
		os.Exit(1)
	}

	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	status, err := LoadStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading status: %v\n", err)
		os.Exit(1)
	}

	command := args[0]

	// Special case for config command - doesn't need config file
	if command == "config" {
		handleConfig()
		return
	}

	switch command {
	case "toggle":
		err = handleToggle(config, status, fanName)
	case "dial":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "dial command requires a value\n")
			os.Exit(1)
		}
		value, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid dial value: %v\n", err)
			os.Exit(1)
		}
		err = handleDial(config, status, fanName, value)
	case "set":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "set command requires a speed value\n")
			os.Exit(1)
		}
		speed, err := strconv.Atoi(args[1])
		if err != nil || speed < 0 || speed > 100 {
			fmt.Fprintf(os.Stderr, "Invalid speed value (must be 0-100): %v\n", args[1])
			os.Exit(1)
		}
		err = handleSet(config, status, fanName, speed)
	case "status":
		err = handleStatus(config, fanName)
	case "list-fans":
		err = handleListFans(config)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleToggle(config *Config, status *Status, fanName string) error {
	fans := selectFans(config, fanName)
	if len(fans) == 0 {
		return fmt.Errorf("no fans found")
	}

	for name, url := range fans {
		currentSpeed, err := getFanSpeed(url)
		if err != nil {
			return fmt.Errorf("failed to get speed for %s: %w", name, err)
		}

		var newSpeed int
		if currentSpeed > 0 {
			status.LastSpeeds[name] = currentSpeed
			newSpeed = 0
		} else {
			newSpeed = status.LastSpeeds[name]
			if newSpeed == 0 {
				newSpeed = 50
			}
		}

		if err := setFanSpeed(url, newSpeed); err != nil {
			return fmt.Errorf("failed to set speed for %s: %w", name, err)
		}

		fmt.Printf("%s: %d%% -> %d%%\n", name, currentSpeed, newSpeed)
	}

	return status.Save()
}

func handleDial(config *Config, status *Status, fanName string, value int) error {
	fans := selectFans(config, fanName)
	if len(fans) == 0 {
		return fmt.Errorf("no fans found")
	}

	var baseSpeed int
	for name, url := range fans {
		currentSpeed, err := getFanSpeed(url)
		if err != nil {
			return fmt.Errorf("failed to get speed for %s: %w", name, err)
		}
		baseSpeed = currentSpeed
		break
	}

	newSpeed := baseSpeed + value
	if newSpeed < 0 {
		newSpeed = 0
	} else if newSpeed > 100 {
		newSpeed = 100
	}

	for name, url := range fans {
		if err := setFanSpeed(url, newSpeed); err != nil {
			return fmt.Errorf("failed to set speed for %s: %w", name, err)
		}
		if newSpeed > 0 {
			status.LastSpeeds[name] = newSpeed
		}
		fmt.Printf("%s: -> %d%%\n", name, newSpeed)
	}

	return status.Save()
}

func handleSet(config *Config, status *Status, fanName string, speed int) error {
	fans := selectFans(config, fanName)
	if len(fans) == 0 {
		return fmt.Errorf("no fans found")
	}

	for name, url := range fans {
		if err := setFanSpeed(url, speed); err != nil {
			return fmt.Errorf("failed to set speed for %s: %w", name, err)
		}
		if speed > 0 {
			status.LastSpeeds[name] = speed
		}
		fmt.Printf("%s: -> %d%%\n", name, speed)
	}

	return status.Save()
}

func handleStatus(config *Config, fanName string) error {
	fans := selectFans(config, fanName)
	if len(fans) == 0 {
		return fmt.Errorf("no fans found")
	}

	for name, url := range fans {
		currentSpeed, err := getFanSpeed(url)
		if err != nil {
			return fmt.Errorf("failed to get speed for %s: %w", name, err)
		}
		
		rpmData, err := getFanRPM(url)
		if err != nil {
			fmt.Printf("%s: %d%% (RPM: unavailable)\n", name, currentSpeed)
		} else {
			fmt.Printf("%s: %d%% (RPM: %d)\n", name, currentSpeed, rpmData)
		}
	}

	return nil
}

func handleListFans(config *Config) error {
	if len(config.Fans) == 0 {
		fmt.Println("No fans configured")
		return nil
	}

	fmt.Println("Configured fans:")
	for name, url := range config.Fans {
		fmt.Printf("  %s: %s\n", name, url)
	}

	return nil
}

func selectFans(config *Config, fanName string) map[string]string {
	if fanName != "" {
		if url, ok := config.Fans[fanName]; ok {
			return map[string]string{fanName: url}
		}
		return nil
	}
	return config.Fans
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "OpenFan Control - Control multiple OpenFan devices\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <command> [args]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  status          Show current fan status (speed and RPM)\n")
	fmt.Fprintf(os.Stderr, "  toggle          Toggle fan on/off (remembers previous speed)\n")
	fmt.Fprintf(os.Stderr, "  dial <value>    Increase/decrease speed by value (-100 to 100)\n")
	fmt.Fprintf(os.Stderr, "  set <speed>     Set speed to specific value (0-100)\n")
	fmt.Fprintf(os.Stderr, "  list-fans       List all configured fans\n")
	fmt.Fprintf(os.Stderr, "  config          Show example configuration and setup instructions\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -fan <name>     Target specific fan by name\n")
	fmt.Fprintf(os.Stderr, "  -h, --help      Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Configuration:\n")
	fmt.Fprintf(os.Stderr, "  Config file: ~/.config/openfan/fans.yaml\n")
	fmt.Fprintf(os.Stderr, "  Status file: ~/.config/openfan/fan-status\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  %s status                  # Show status of all fans\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s set 50                  # Set all fans to 50%%\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -fan right dial 10      # Increase 'right' fan by 10%%\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s toggle                  # Toggle all fans on/off\n", os.Args[0])
}

func handleConfig() {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "openfan", "fans.yaml")
	
	fmt.Println("OpenFan Control Configuration")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Printf("Configuration file location: %s\n", configPath)
	fmt.Println()
	fmt.Println("Example configuration file:")
	fmt.Println("---------------------------")
	fmt.Println("fans:")
	fmt.Println("  # Fan names can be any identifier you choose")
	fmt.Println("  # URLs should point to your OpenFan devices")
	fmt.Println("  ")
	fmt.Println("  # Using mDNS hostname (recommended)")
	fmt.Println("  left: http://openfan-left.local")
	fmt.Println("  right: http://openfan-right.local")
	fmt.Println("  ")
	fmt.Println("  # Using IP address")
	fmt.Println("  top: http://192.168.1.100")
	fmt.Println("  bottom: http://192.168.1.101")
	fmt.Println("  ")
	fmt.Println("  # Using full mDNS hostname")
	fmt.Println("  desk: http://uOpenFan-Desk-48CA43DBD758")
	fmt.Println()
	fmt.Println("Setup instructions:")
	fmt.Println("------------------")
	fmt.Println("1. Create the configuration directory:")
	fmt.Printf("   mkdir -p %s\n", filepath.Join(homeDir, ".config", "openfan"))
	fmt.Println()
	fmt.Println("2. Create the configuration file:")
	fmt.Printf("   nano %s\n", configPath)
	fmt.Println()
	fmt.Println("3. Add your fans to the configuration file using the format above")
	fmt.Println()
	fmt.Println("4. Test your configuration:")
	fmt.Printf("   %s list-fans\n", os.Args[0])
	fmt.Printf("   %s status\n", os.Args[0])
	fmt.Println()
	fmt.Println("Notes:")
	fmt.Println("------")
	fmt.Println("- Fan names are case-sensitive")
	fmt.Println("- URLs must include the protocol (http://)")
	fmt.Println("- OpenFan devices run on port 80 (default HTTP port)")
	fmt.Println("- You can find device hostnames in your router's DHCP client list")
}