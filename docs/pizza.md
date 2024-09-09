## pizza

OpenSauced CLI

### Synopsis

A command line utility for insights, metrics, and generating CODEOWNERS documentation for your open source projects

```
pizza <command> <subcommand> [flags]
```

### Options

```
  -c, --config string       The saucectl config (default "~/.sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -h, --help                help for pizza
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza completion](pizza_completion.md)	 - Generate the autocompletion script for the specified shell
* [pizza generate](pizza_generate.md)	 - Generates documentation and insights from your codebase
* [pizza insights](pizza_insights.md)	 - Gather insights about git contributors, repositories, users and pull requests
* [pizza login](pizza_login.md)	 - Log into the CLI via GitHub
* [pizza version](pizza_version.md)	 - Displays the build version of the CLI

