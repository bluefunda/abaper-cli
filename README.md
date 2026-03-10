# ABAPer CLI

Command line interface for the [ABAPer](https://abaper.bluefunda.com) platform. Communicates with ABAPer APIs exposed through the ABAPer gateway (`abaper-gw`).

## Installation

### From GitHub Releases

Download the binary for your platform from the [releases page](https://github.com/bluefunda/abaper-cli/releases).

```bash
# macOS (Apple Silicon)
curl -L https://github.com/bluefunda/abaper-cli/releases/latest/download/abaper-cli_<version>_darwin_arm64.tar.gz | tar xz
sudo mv abaper /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/bluefunda/abaper-cli/releases/latest/download/abaper-cli_<version>_linux_amd64.tar.gz | tar xz
sudo mv abaper /usr/local/bin/
```

### From Source

```bash
go install github.com/bluefunda/abaper-cli/cmd/abaper@latest
```

### Docker

```bash
docker run --rm bluefunda/abaper version
```

## Usage

### Authentication

```bash
# Login via browser-based device flow
abaper login

# Check status
abaper status

# Logout
abaper logout
```

### Developer Workflows

```bash
# Generate a new ABAP program
abaper generate --type program --name ZMY_PROGRAM

# Generate a class with source from file
abaper generate --type class --name ZCL_MY_CLASS --source-file my_class.abap

# Deploy (upload + activate)
abaper deploy --type program --name ZMY_PROGRAM --source-file program.abap

# Check CLI version
abaper version
```

### Output Formats

All commands support `--output json` for machine-readable output:

```bash
abaper status -o json
```

## Configuration

Configuration is loaded from (in order of precedence):

1. CLI flags (`--base-url`, `--realm`)
2. Environment variables (`ABAPER_BASE_URL`, `ABAPER_REALM`, `ABAPER_ORG`)
3. Config file `~/.abaper/config.yaml`

### Config File

```yaml
# ~/.abaper/config.yaml
base_url: https://api.bluefunda.com
realm: trm
org: default
```

Tokens are stored separately in `~/.abaper/tokens.yaml` with restricted permissions (0600).

## Developer Setup

### Prerequisites

- Go 1.25+
- golangci-lint (for linting)

### Build

```bash
make build          # Build for current platform
make build-all      # Cross-compile for all platforms
make test           # Run tests
make lint           # Run linter
make vet            # Run go vet
```

### Docker

```bash
make docker-build   # Build Docker image
make docker-run     # Run CLI in container
```

## Architecture

The CLI follows the same API patterns as [abaper-editor](https://github.com/bluefunda/abaper-editor) and [abaper-vscode](https://github.com/bluefunda/abaper-vscode):

- **Authentication**: OAuth2 device authorization flow via Keycloak (same as `abaper-vscode`)
- **API Client**: Calls ABAPer APIs through the KrakenD gateway (`abaper-gw`) at `/abaper/api/v1/*`
- **Request Headers**: `Authorization: Bearer <token>`, `X-Realm`, `X-SAP-*` headers
- **Response Format**: `{ "success": bool, "data": T, "error": string }`

## Release Process

Releases are automated via [release-please](https://github.com/googleapis/release-please) and follow the standards from [release-foundry](https://github.com/bluefunda/release-foundry).

1. Merge PRs with [conventional commit](https://www.conventionalcommits.org/) titles
2. release-please creates a Release PR with version bump and changelog
3. Merging the Release PR triggers:
   - GoReleaser builds multi-platform binaries
   - Binaries attached to GitHub Release
   - Docker image pushed to `bluefunda/abaper`
