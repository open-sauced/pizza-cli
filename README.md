# 🍕 Pizza CLI

This CLI can be used for all things OpenSauced!

```
❯ pizza

A command line utility for insights, metrics, and all things OpenSauced

Usage:
  pizza <command> <subcommand> [flags]

Available Commands:
  bake        Use a pizza-oven to source git commits into OpenSauced
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  login       Log into the CLI application via GitHub
  repo-query  Ask questions about a GitHub repository

Flags:
  -h, --help   help for pizza

Use "pizza [command] --help" for more information about a command.
```

---

### 🖥️ Local Development

You'll need a few tools to get started:

- The [Go toolchain](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/) (for linting and other tooling)
- Make

To lint, run `make lint`. To run tests, run `make test`. To build, run `make build`.

### 🏗️ Installation

#### Manual build and install

```
make install
```

This is a convenience target for building and dropping the pizza CLI into
`/usr/local/bin/` (which requires `sudo` permissions).
Make sure you have that directory in your path: `export PATH="$PATH:/usr/local/bin"`.
Otherwise, you can build it manually with `make build` and `mv build/pizza <somewhere-in-your-path>`.

In the future, we will be dropping regular GitHub releases where you can easily
download and install pre-built binaries.
