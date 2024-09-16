<div align="center">
  <br>
  <img alt="Open Sauced" src="https://github.com/open-sauced/assets/blob/main/logos/logo-on-dark.png" >
  <h1>🍕 Pizza CLI 🍕</h1>
  <strong>A Go command line interface for managing code ownership and project insights with OpenSauced!</strong>
  <br>
</div>
<br>
<div align="center">
  <img src="https://img.shields.io/github/languages/code-size/open-sauced/pizza" alt="GitHub code size in bytes">
  <a href="https://github.com/open-sauced/pizza/issues">
    <img src="https://img.shields.io/github/issues/open-sauced/pizza" alt="GitHub issues">
  </a>
  <a href="https://github.com/open-sauced/pizza-cli/releases">
    <img src="https://img.shields.io/github/v/release/open-sauced/pizza-cli.svg?style=flat" alt="GitHub Release">
  </a>
  <a href="https://twitter.com/saucedopen">
    <img src="https://img.shields.io/twitter/follow/saucedopen?label=Follow&style=social" alt="Twitter">
  </a>
    <a href="https://opensauced.pizza/docs/tools/pizza-cli">
    <img src="https://img.shields.io/badge/%F0%9F%92%A1%20OpenSauced-Docs-00ACD7.svg?style=flat-square">
  </a>
</div>

---

# 📦 Install

#### Homebrew

```sh
brew install open-sauced/tap/pizza
```

#### NPM

```sh
npm i -g pizza
```

You can also use `npx` to run one-off commands without installing anything:

```sh
npx pizza@latest generate codeowners .
```

# 🍕 Pizza Action

Use [the Pizza GitHub Action](https://github.com/open-sauced/pizza-action) for running `pizza` operations in GitHub CI/CD,
like automated `CODEOWNERS` updating and pruning:

```yaml
jobs:
  pizza-action:
    runs-on: ubuntu-latest
    steps:
      - name: Pizza Action
        uses: open-sauced/pizza-action@v2
        with:
          # Optional: Whether to commit and create a PR for "CODEOWNER" changes
          commit-and-pr: "true"
          # Optional: Title of the PR for review by team
          pr-title: "chore: update repository codeowners"
```

# 📝 Docs

- [Pizza.md](./docs/pizza.md): In depth docs on each command, option, and flag.
- [OpenSauced.pizza/docs](https://opensauced.pizza/docs/tools/pizza-cli/): Learn
  how to use the Pizza command line tool and how it works with the rest of the OpenSauced
  ecosystem.

# ✨ Usage

## Codeowners generation

Use the `codeowners` command to generate a GitHub style `CODEOWNERS` file or a more agnostic `OWNERS` file.
This can be used to granularly define what experts and entities have the
most context and knowledge on certain parts of a codebase.

It's expected that there's a `.sauced.yaml` config file in the given path or in
your home directory (as `~/.sauced.yaml`):

```sh
pizza generate codeowners /path/to/local/git/repo
```

Running this command will iterate the git ref-log to determine who to set as a code
owner based on the number of lines changed for that file within the given time range.
The first owner is the entity with the most lines changed. This command uses a `.sauced.yaml` configuration
to attribute emails in commits with the given entities in the config (like GitHub usernames or teams).
See [the section on the configuration schema for more details](#-configuration-schema)

### 🚀 New in v1.4.0: Generate Config

The `pizza generate config` command has been added to help you create `.sauced.yaml` configuration files for your projects.
This command allows you to generate configuration files with various options:

```sh
pizza generate config /path/to/local/git/repo
```

This command will iterate the git ref-log and inspect email signatures for commits
and, in interactive mode, ask you to attribute those users with GitHub handles. Once finished, the resulting
`.sauced.yaml` file can be used to attribute owners in a `CODEOWNERS` file during `pizza generate codeowners`.

#### Flags:

- `-i, --interactive`: Enter interactive mode to attribute each email manually
- `-o, --output-path string`: Set the directory for the output file
- `-h, --help`: Display help for the command

#### Examples:

1. Generate a config file in the current directory:
   ```sh
   pizza generate config ./
   ```

2. Generate a config file interactively:
   ```sh
   pizza generate config ./ -i
   ```

3. Generate a config file from the current directory and place resulting `.sauced.yaml` in a specific output directory:
   ```sh
   pizza generate config ./ -o /path/to/directory
   ```

## OpenSauced Contributor Insight from `CODEOWNERS`

You can create an [OpenSauced Contributor Insight](https://opensauced.pizza/docs/features/contributor-insights/)
from a local `CODEOWNERS` file:

```
pizza generate insight /path/to/repo/with/CODEOWNERS/file
```

This will parse the `CODEOWNERS` file and create a Contributor Insight on the OpenSauced platform.
This allows you to track insights and metrics for those codeowners, powered by OpenSauced.

## Insights

You can get metrics and insights on repositories, contributors, and more:

```
pizza insights [sub-command]
```

This powerful command lets you compose many metrics and insights together, all
powered by OpenSauced's API. Use the `--output` flag to output the results as yaml, json, csv, etc.

# 🎷 Configuration schema

```yaml
# Configuration for attributing commits with emails to individual entities.
# Used during "pizza generate codeowners".
attribution:

  # Keys can be GitHub usernames.
  jpmcb:

    # List of emails associated with the given GitHub login.
    # The commits associated with these emails will be attributed to
    # this GitHub login in this yaml map. Any number of emails may be listed.
    - john@opensauced.pizza
    - hello@johncodes.com

  # Keys may also be GitHub teams. This is useful for orchestrating multiple
  # people to a sole GitHub team.
  open-sauced/engineering:
    - john@opensauced.pizza
    - other-user@email.com
    - other-user@no-reply.github.com

  # Keys can also be agnostic names which will land as keys in "OWNERS" files
  # when the "--owners-style-file" flag is set.
  John McBride
    - john@opensauced.pizza

# Used during codeowners generation: if there are no code owners found
# for a file within the time range, the list of fallback entities 
# will be used
attribution-fallback:
  - open-sauced/engineering
  - some-other-github-login
```
