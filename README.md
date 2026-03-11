# ABAPer CLI

Command line interface for the [ABAPer](https://abaper.bluefunda.com) platform. Communicates with ABAPer APIs exposed through the ABAPer gateway (`abaper-gw`).

## Quick Start

```bash
# Install (macOS)
brew tap bluefunda/tap
brew install abaper

# Authenticate
abaper login

# Verify connection
abaper status

# Create and deploy an ABAP program
abaper generate --type program --name ZMY_REPORT
abaper deploy --type program --name ZMY_REPORT --source-file report.abap
```

## Installation

### Homebrew (macOS)

```bash
brew tap bluefunda/tap
brew install abaper
```

This installs the binary and the man page automatically.

### From GitHub Releases

Download the binary for your platform from the [releases page](https://github.com/bluefunda/abaper-cli/releases).

| Platform       | Archive                              |
|----------------|--------------------------------------|
| macOS (ARM64)  | `abaper_<version>_darwin_arm64.tar.gz` |
| macOS (AMD64)  | `abaper_<version>_darwin_amd64.tar.gz` |
| Linux (AMD64)  | `abaper_<version>_linux_amd64.tar.gz`  |
| Linux (ARM64)  | `abaper_<version>_linux_arm64.tar.gz`  |
| Windows (AMD64)| `abaper_<version>_windows_amd64.zip`   |
| Windows (ARM64)| `abaper_<version>_windows_arm64.zip`   |

```bash
# macOS / Linux
curl -L https://github.com/bluefunda/abaper-cli/releases/latest/download/abaper_<version>_<os>_<arch>.tar.gz | tar xz
sudo mv abaper /usr/local/bin/
```

### From Source

```bash
go install github.com/bluefunda/abaper-cli/cmd/abaper@latest
```

### Docker

```bash
docker pull bluefunda/abaper
docker run --rm bluefunda/abaper version
```

## Commands

### `abaper login`

Authenticate with the ABAPer platform using the OAuth2 device authorization flow. Opens a browser window for interactive login. Credentials are stored locally in `~/.abaper/tokens.yaml` with restricted permissions (0600).

```bash
abaper login
```

### `abaper logout`

Clear stored credentials by removing the local token file.

```bash
abaper logout
```

### `abaper status`

Show connection and authentication status. Reports the configured base URL, realm, organization, authentication state, and API health.

```bash
abaper status
abaper status -o json
```

### `abaper generate`

Create ABAP objects on the target SAP system. Supports programs, classes, and interfaces. Accepts source from a file or generates default templates.

```bash
# Generate with default template
abaper generate --type program --name ZMY_PROGRAM

# Generate from source file
abaper generate --type class --name ZCL_MY_CLASS --source-file my_class.abap

# Generate an interface
abaper generate --type interface --name ZIF_MY_INTERFACE
```

**Flags:**

| Flag            | Required | Default   | Description                        |
|-----------------|----------|-----------|------------------------------------|
| `--name`        | Yes      | —         | Object name                        |
| `--type`        | No       | `program` | Object type: program, class, interface |
| `--source-file` | No       | —         | Path to ABAP source file           |

### `abaper deploy`

Upload source code and activate an ABAP object in a single step. Performs a save followed by activation, matching the workflow in ABAPer Editor.

```bash
abaper deploy --type program --name ZMY_PROGRAM --source-file program.abap
```

**Flags:**

| Flag            | Required | Default   | Description                        |
|-----------------|----------|-----------|------------------------------------|
| `--name`        | Yes      | —         | Object name                        |
| `--type`        | No       | `program` | Object type: program, class, interface |
| `--source-file` | Yes      | —         | Path to ABAP source file           |

### `abaper version`

Print the CLI version.

```bash
abaper version
```

## Man Page

A Unix man page is included. After installing the binary from a release archive:

```bash
# Install the man page (included in release archives)
sudo install -m 644 man/abaper.1 /usr/local/share/man/man1/abaper.1

# Or using make
sudo make install-man

# Then use it like any other Unix command
man abaper
```

## Global Flags

These flags are available on all commands:

| Flag              | Description                                        |
|-------------------|----------------------------------------------------|
| `--base-url`      | ABAPer API base URL (default: `https://api.bluefunda.com`) |
| `--realm`         | Keycloak realm (default: `trm`)                    |
| `-o`, `--output`  | Output format: `text`, `json` (default: `text`)    |

## Output Formats

All commands support `--output json` for machine-readable output, useful for scripting and CI/CD pipelines:

```bash
abaper status -o json | jq '.authenticated'
```

## Configuration

Configuration is loaded in the following order of precedence:

1. **CLI flags** — `--base-url`, `--realm`
2. **Environment variables** — `ABAPER_BASE_URL`, `ABAPER_REALM`, `ABAPER_ORG`
3. **Config file** — `~/.abaper/config.yaml`

### Config File

```yaml
# ~/.abaper/config.yaml
base_url: https://api.bluefunda.com
realm: trm
org: default
```

### Environment Variables

| Variable          | Description                |
|-------------------|----------------------------|
| `ABAPER_BASE_URL` | Override the API base URL  |
| `ABAPER_REALM`    | Override the Keycloak realm|
| `ABAPER_ORG`      | Override the organization  |

### Files

| Path                      | Description                           |
|---------------------------|---------------------------------------|
| `~/.abaper/config.yaml`   | Configuration file                   |
| `~/.abaper/tokens.yaml`   | OAuth2 tokens (permissions: 0600)    |

## Authentication

ABAPer CLI uses the **OAuth2 device authorization flow** via Keycloak (same flow as ABAPer VS Code extension):

1. `abaper login` requests a device code from the authorization server
2. Your browser opens to the verification URL
3. The CLI polls for authorization completion
4. Access and refresh tokens are stored locally

Tokens are **automatically refreshed** when expired. The CLI uses the `cai-cli` OAuth2 client ID.

## Docker Usage

The CLI is available as a Docker image on Docker Hub under [`bluefunda/abaper`](https://hub.docker.com/r/bluefunda/abaper).

```bash
# Pull the latest image
docker pull bluefunda/abaper:latest

# Pull a specific version
docker pull bluefunda/abaper:v1.0.0

# Run a command
docker run --rm bluefunda/abaper version
docker run --rm bluefunda/abaper status -o json

# Mount config for authenticated commands
docker run --rm -v ~/.abaper:/root/.abaper bluefunda/abaper status
```

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

The CLI follows the same API patterns as [ABAPer Editor](https://github.com/bluefunda/abaper-editor) and [ABAPer VS Code](https://github.com/bluefunda/abaper-vscode):

- **Authentication**: OAuth2 device authorization flow via Keycloak
- **API Client**: Calls ABAPer APIs through the KrakenD gateway (`abaper-gw`) at `/abaper/api/v1/*`
- **Request Headers**: `Authorization: Bearer <token>`, `X-Realm`, `X-SAP-*` headers
- **Response Format**: `{ "success": bool, "data": T, "error": string }`

## Release Process

Releases are automated via [release-please](https://github.com/googleapis/release-please):

1. Merge PRs with [conventional commit](https://www.conventionalcommits.org/) titles
2. release-please creates a Release PR with version bump and changelog
3. Merging the Release PR triggers:
   - GoReleaser builds multi-platform binaries (named `abaper_<version>_<os>_<arch>`)
   - Binaries attached to GitHub Release
   - Docker image pushed to [`bluefunda/abaper`](https://hub.docker.com/r/bluefunda/abaper)
