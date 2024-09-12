<div align="center">
  <br>
  <img alt="Open Sauced" src="https://github.com/open-sauced/assets/blob/main/logos/logo-on-dark.png" >
  <h1>🍕 Pizza CLI 🍕</h1>
  <strong>A Go command line interface for managing code ownership and project insights with OpenSauced!</strong>
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
  generate    Generate configurations and codeowner files
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

#### NPM

```sh
npm i -g pizza
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
