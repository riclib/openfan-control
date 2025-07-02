# OpenFan Micro Control

A Go command-line tool for controlling multiple OpenFan Micro devices.

## Features

- Control multiple fans simultaneously or individually
- Toggle fans on/off with memory of previous speed
- Adjust fan speeds incrementally or set specific values
- View current fan status including speed percentage and RPM
- Persistent storage of last known fan speeds

## Installation

### Using go install (recommended)

```bash
go install github.com/riclib/openfan-control@latest
```

This will install the `openfan-control` binary to your `$GOPATH/bin` directory (or `$HOME/go/bin` if GOPATH is not set).

### Building from source

```bash
git clone https://github.com/riclib/openfan-control.git
cd openfan-control
go build -o openfan-control
```

## Configuration

Create a configuration file at `~/.config/openfan/fans.yaml`:

```yaml
fans:
  left: http://openfan-left.local
  right: http://openfan-right.local
  top: http://192.168.1.100
```

Fan names can be any identifier you choose. URLs should point to your OpenFan devices (using mDNS hostnames or IP addresses).

## Usage

```bash
# Show help
openfan-control --help
openfan-control -h

# Show configuration help
openfan-control config

# Show status of all fans
openfan-control status

# Show status of specific fan
openfan-control -fan left status

# List all configured fans
openfan-control list-fans

# Set all fans to 50%
openfan-control set 50

# Set specific fan to 75%
openfan-control -fan right set 75

# Increase all fans by 10%
openfan-control dial 10

# Decrease specific fan by 20%
openfan-control -fan top dial -20

# Toggle all fans on/off
openfan-control toggle

# Toggle specific fan
openfan-control -fan left toggle
```

## Commands

- **status** - Display current fan speed (%) and RPM
- **set <speed>** - Set fan speed to a specific value (0-100)
- **dial <value>** - Increase/decrease fan speed by given amount
- **toggle** - Turn fan on/off, remembering previous speed
- **list-fans** - List all configured fans and their URLs
- **config** - Show example configuration and setup instructions

## Options

- **-fan <name>** - Target a specific fan instead of all fans
- **-h, --help** - Show help message

## Files

- Configuration: `~/.config/openfan/fans.yaml`
- Status storage: `~/.config/openfan/fan-status`

## Notes

- When controlling multiple fans, the `dial` command reads the current speed from the first fan and applies the same resulting speed to all fans, keeping them synchronized
- The `toggle` command remembers the last non-zero speed for each fan
- If a fan has never been set to a non-zero speed, toggle will default to 50%

## Requirements

- Go 1.21 or later
- Network access to OpenFan Micro devices
- OpenFan Micro API (runs on port 80)
