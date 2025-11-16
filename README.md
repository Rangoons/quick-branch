# quick-branch

A fast CLI tool for streamlining your workflow with [Linear](https://linear.app) issues and git. Assign yourself to issues, update statuses, and create branches with proper naming—all from your terminal. Use **turbo mode** to do everything in one command.

## Features

- **Turbo mode** - Assign yourself, update status to "In Dev", and checkout branch in one command
- **Issue assignment** - Assign yourself to Linear issues from the terminal
- **Status updates** - Update issue status directly from the CLI
- **Quick branch creation** - Create and checkout git branches using Linear's branch naming conventions
- **Copy to clipboard** - Instantly copy issue URLs or branch names
- **Beautiful issue display** - View issue details with formatted markdown descriptions
- **Secure authentication** - Store your Linear API key locally in platform-appropriate config directories
- **Fast** - No UI overhead, works entirely in your terminal

## Installation

### Homebrew (recommended)

```bash
brew install rangoons/tap/quick-branch
```

### From source

Requires Go 1.24 or later:

```bash
go install github.com/rangoons/quick-branch/cmd/quick-branch@latest
```

### Download binaries

Download pre-built binaries for your platform from the [releases page](https://github.com/rangoons/quick-branch/releases).

## Quick Start

1. **Authenticate with Linear:**

   ```bash
   quick-branch auth
   ```

   Enter your Linear API key when prompted. Get your API key from [Linear Settings → API](https://linear.app/settings/api).

2. **Start working on an issue (the fast way):**

   ```bash
   quick-branch start ABC-123 --turbo
   ```

   This assigns you to the issue, updates the status to "In Dev", and checks out a new branch—all in one command!

## Usage

### Authentication

Store your Linear API key:

```bash
quick-branch auth
```

Your API key is stored securely in:
- **macOS**: `~/Library/Application Support/quick-branch/config.yaml`
- **Linux**: `~/.config/quick-branch/config.yaml`
- **Windows**: `%AppData%\quick-branch\config.yaml`

### Working with Issues

#### View and interact with issues

```bash
quick-branch issue <issue-id> [flags]
```

**Examples:**

```bash
# View issue details
quick-branch issue ABC-123 -v

# Copy issue URL to clipboard
quick-branch issue ABC-123 --url

# Copy branch name to clipboard
quick-branch issue ABC-123 --branch

# Create and checkout a new branch
quick-branch issue ABC-123 --checkout

# Combine flags: view details and checkout
quick-branch issue ABC-123 -v -c
```

**Flags:**

- `-u, --url` - Copy issue URL to clipboard
- `-b, --branch` - Copy branch name to clipboard
- `-c, --checkout` - Create and checkout a new branch with the Linear branch name
- `-v, --verbose` - Display issue description with formatted markdown

#### Start working on an issue

```bash
quick-branch start <issue-id> [flags]
```

Assign yourself to an issue with optional status update and branch checkout.

**Examples:**

```bash
# Turbo mode: do everything in one command
quick-branch start ABC-123 --turbo

# Assign yourself to an issue
quick-branch start ABC-123

# Assign and update status to "In Dev"
quick-branch start ABC-123 --status

# Assign and checkout the branch
quick-branch start ABC-123 --checkout

# Assign, update status, and checkout branch
quick-branch start ABC-123 -s -c
```

**Flags:**

- `-t, --turbo` - Assign yourself, update status to "In Dev", and checkout branch (all-in-one!)
- `-s, --status` - Update issue status to "In Dev"
- `-c, --checkout` - Create and checkout a new branch with the Linear branch name

## Workflow Examples

### The Fast Way (Turbo Mode)

Start working on an issue instantly:

```bash
# Get started with everything in one command
quick-branch start PRJ-456 --turbo

# Output:
# Success! Assigned John Doe to PRJ-456
# Success! Updated PRJ-456 to In Dev
# Success! Now working on rangoons/prj-456-add-user-authentication
```

### The Traditional Way

Break down each step for more control:

```bash
# 1. See what you're working on
quick-branch issue PRJ-456 -v

# Output:
# Add user authentication: Backlog
#
# ## Description
# Implement OAuth2 authentication with support for GitHub and Google...

# 2. Assign yourself and update status
quick-branch start PRJ-456 --status

# Output:
# Success! Assigned John Doe to PRJ-456
# Success! Updated PRJ-456 to In Dev

# 3. Create a branch and start working
quick-branch issue PRJ-456 -c

# Output:
# Success! Now working on rangoons/prj-456-add-user-authentication

# 4. Later: quickly grab the issue URL for a PR description
quick-branch issue PRJ-456 -u

# Output:
# Copied issue url to clipboard
```

## Configuration

Configuration is stored in YAML format at the locations mentioned above.

**Example config:**

```yaml
api_key: lin_api_your_key_here
```

You can also set configuration via environment variables with the `QUICK_BRANCH_` prefix:

```bash
export QUICK_BRANCH_API_KEY=lin_api_your_key_here
```

## Development

### Prerequisites

- Go 1.24 or later
- Linear API key

### Building from source

```bash
git clone https://github.com/rangoons/quick-branch.git
cd quick-branch
go build -o quick-branch ./cmd/quick-branch
```

### Running tests

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Glamour](https://github.com/charmbracelet/glamour) - Terminal markdown rendering
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [genqlient](https://github.com/Khan/genqlient) - GraphQL client generation
