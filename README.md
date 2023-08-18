<div align="center">
  <br>
  <img alt="Open Sauced" src="https://i.ibb.co/7jPXt0Z/logo1-92f1a87f.png" width="300px">
  <h1>🍕 Pizza CLI 🍕</h1>
  <strong>A Go command line interface for all things OpenSauced!!</strong>
  <br>
</div>
<br>
<p align="center">
  <img src="https://img.shields.io/github/languages/code-size/open-sauced/pizza" alt="GitHub code size in bytes">
  <a href="https://github.com/open-sauced/pizza/issues">
    <img src="https://img.shields.io/github/issues/open-sauced/pizza" alt="GitHub issues">
  </a>
  <a href="https://github.com/open-sauced/api.opensauced.pizza/releases">
    <img src="https://img.shields.io/github/v/release/open-sauced/pizza.svg?style=flat" alt="GitHub Release">
  </a>
  <a href="https://discord.gg/U2peSNf23P">
    <img src="https://img.shields.io/discord/714698561081704529.svg?label=&logo=discord&logoColor=ffffff&color=7389D8&labelColor=6A7EC2" alt="Discord">
  </a>
  <a href="https://twitter.com/saucedopen">
    <img src="https://img.shields.io/twitter/follow/saucedopen?label=Follow&style=social" alt="Twitter">
  </a>
</p>

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

### 📦 Install

There are several methods for downloading and installing the `pizza` CLI:

#### Homebrew

```sh
brew install open-sauced/tap/pizza
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
> It may not be advisable to pipe scripts from GitHub directly into
> a command line interpreter! If you do not fully trust the source, first
> download the script, inspect it manually to ensure its integrity, and then
> run it:
> ```sh
> curl -fsSL https://raw.githubusercontent.com/open-sauced/pizza-cli/main/install.sh > install.sh
> vim install.sh
> ./install.sh
> ```

#### Manual build and install

```
make install
```

This is a convenience `make` target for building and dropping the pizza CLI into
`/usr/local/bin/` (which requires `sudo` permissions).
Make sure you have that directory in your path: `export PATH="$PATH:/usr/local/bin"`.
Otherwise, you can build it manually with `make build` and `mv build/pizza <somewhere-in-your-path>`.

In the future, we will be dropping regular GitHub releases where you can easily
download and install pre-built binaries.

### 🖥️ Local Development

You'll need a few tools to get started:

- The [Go toolchain](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/) (for linting and other tooling)
- Make

To lint, run `make lint`. To run tests, run `make test`. To build, run `make build`.
