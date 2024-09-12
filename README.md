<div align="center">
  <br>
  <img alt="Open Sauced" src="https://github.com/open-sauced/assets/blob/main/logos/logo-on-dark.png">
  <h1>🍕 Pizza CLI 🍕</h1>
  <strong>A Go command line interface for managing code ownership and project insights with OpenSauced!</strong>
  <br>
</div>
<br>

# 📦 Install

#### Homebrew

```sh
brew install open-sauced/tap/pizza
```

#### NPM

```sh
npm i -g pizza
```

### Docker

```sh
$ docker run ghcr.io/open-sauced/pizza-cli:latest
```

For commands that require access to your file system (like `generate codeowners`), ensure
you pass a volume to the docker container:

```sh
$ docker run -v /local/path:/container/path ghcr.io/open-sauced/pizza-cli:latest \
    generate codeowners /container/path
```

For example, to mount your entire home directory (which may include a `.sauced.yaml` file
alongside the project you want to generate a `CODEOWNERS` file for):

```sh
$ docker run -v ~/:/app ghcr.io/open-sauced/pizza-cli:latest \
    codeowners /app/workspace/gopherlogs -c /app/.sauced.yaml
```

### Go install

Using the Go tool-chain, you can install the binary directly:

```sh
$ go install github.com/open-sauced/pizza-cli@latest
```

Warning! You should have the `GOBIN` env var setup to point to a persistent
location in your `PATH`. After Go 1.16, this defaults to `GOPATH[0]/bin`.

### Manual install

Download a pre-built artifact from [the GitHub releases](https://github.com/open-sauced/pizza-cli/releases):

```sh
# Make the binary executable
$ chmod +x ~/Downloads/pizza-linux-arm64

# Move the binary into a location in the PATH
# Warning: the location where you drop the binary may differ!
$ mv ~/Downloads/pizza-linux-arm64 /usr/local/share/bin/pizza
```

#### Direct script install

```sh
curl -fsSL https://raw.githubusercontent.com/open-sauced/pizza-cli/main/install.sh | sh
```

This is a convenience script that can be downloaded from GitHub directly and
piped into `sh` for conveniently downloading the latest GitHub release of the
`pizza` CLI.

Once download is completed, you can move the binary to a convenient location in
your system's `$PATH`.

> [!WARNING]
> It's _probably_ not advisable to pipe scripts from GitHub directly into
> a command line interpreter! If you do not fully trust the source, first
> download the script, inspect it manually to ensure integrity, and then
> run it:
> ```sh
> curl -fsSL https://raw.githubusercontent.com/open-sauced/pizza-cli/main/install.sh > install.sh
> vim install.sh
> ./install.sh
> ```

#### Manual build

Clone this repository. Then, using the Go tool-chain, you can build a binary:

```
$ go build -o build/pizza main.go
```

Warning! There may be unsupported features, breaking changes, or experimental
patches on the tip of the repository. Go and build with caution!

# ✨ Usage

## Codeowners generation

Use the `codeowners` command to generate an `OWNERS` file or GitHub style `CODEOWNERS` file.
This can be used to granularly define what experts and entities have the
most context and knowledge on certain parts of a codebase.

```
❯ pizza generate codeowners -h

Generates a CODEOWNERS file for a given git repository. This uses a ~/.sauced.yaml
configuration to attribute emails with given entities.

The generated file specifies up to 3 owners for EVERY file in the git tree based on the
number of lines touched in that specific file over the specified range of time.

Usage:
  pizza generate codeowners path/to/repo [flags]

Flags:
      --owners-style-file   Whether to generate an agnostic OWNERS style file.
  -h, --help                help for codeowners
  -r, --range int           The number of days to lookback (default 90)

Global Flags:
      --beta                Shorthand for using the beta OpenSauced API endpoint ("https://beta.api.opensauced.pizza").
                            Supersedes the '--endpoint' flag
  -c, --config string       The codeowners config (default ".sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -e, --endpoint string     The API endpoint to send requests to (default "https://api.opensauced.pizza")
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
  -o, --output string       The formatting for command output. One of: (table, yaml, csv, json) (default "table")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

## Configuration

```yaml
# Configuration for attributing commits with emails to individual entities.
# Used during "pizza generate codeowners".
attribution:

  # Keys can be GitHub usernames. For the "--github-codeowners" style codeowners
  # generation, these keys must be valid GitHub usernames.
  jpmcb:

    # List of emails associated with the given entity.
    # The commits associated with these emails will be attributed to
    # the entity in this yaml map. Any number of emails may be listed.
    - john@opensauced.pizza
    - hello@johncodes.com

  # Entities may also be GitHub teams.
  open-sauced/engineering:
    - john@opensauced.pizza
    - other-user@email.com
    - other-user@no-reply.github.com

  # They can also be agnostic names which will land as keys in OWNERS files
  John McBride
    - john@opensauced.pizza
```

### 🚀 New in v1.4.0: Generate Config

The `pizza generate config` command has been added to help you create configuration files for your projects. This command allows you to generate configuration files with various options:

```sh
pizza generate config [flags]
```

#### Flags:

- `-i, --interactive`: Enter interactive mode to attribute each email manually
- `-o, --output-path string`: Set the directory for the output file
- `-h, --help`: Display help for the command

#### Examples:

1. Generate a config file in the current directory:
   ```sh
   pizza generate config
   ```

2. Generate a config file interactively:
   ```sh
   pizza generate config -i
   ```

3. Generate a config file in a specific directory:
   ```sh
   pizza generate config -o /path/to/directory
   ```

# 🚜 Development

### 🔨 Requirements

There are a few things you'll need to get started:

- The [1.22 `go` programming language](https://go.dev/doc/install) toolchain and dev environment (for example, the [VScode Go plugin](https://code.visualstudio.com/docs/languages/go) is very good).
- The [`just` command runner](https://github.com/casey/just) for easy operations
