# OpenFan Control Project Information

## Project Overview
OpenFan Control is a CLI tool for controlling fan speeds via HTTP API.

## Key Commands
- `openfan status` - Shows status of all configured fans
- `openfan set <speed>` - Sets speed for all fans (0-100%)
- `openfan set <speed> -fan <name>` - Sets speed for specific fan
- `openfan resume` - Resumes last known speeds
- `openfan add <name> <url>` - Adds a new fan
- `openfan remove <name>` - Removes a fan
- `openfan list` - Lists all configured fans

## Project Structure
- `main.go` - Main entry point and command handling
- `api.go` - HTTP API client functions
- `config.go` - Configuration management (fans stored in ~/.config/openfan/fans.yaml)
- `status.go` - Status tracking (last speeds stored in ~/.config/openfan/status.yaml)

## Key Functions
- `handleSet()` (main.go:159-186) - Sets fan speeds, processes all fans even if some fail
- `selectFans()` (main.go:215-223) - Determines which fans to target
- `setFanSpeed()` (api.go:51-78) - Makes HTTP API calls to control fans
- `handleStatus()` (main.go:189-212) - Gets current fan status

## Testing & Linting
- No specific test framework or linting commands found in the codebase yet
- Ask user for preferred Go linting/testing commands if needed

## Recent Fixes
- Fixed issue where `openfan set` would only change one fan if another failed
- Now continues processing all fans and reports errors at the end
- Fixed API response parsing in setFanSpeed() - API returns status "ok" not "success"