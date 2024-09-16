## pizza generate

Generates documentation and insights from your codebase

### Synopsis

The 'generate' command provides tools to automate the creation of important project documentation and derive insights from your codebase.

```
pizza generate [subcommand] [flags]
```

### Options

```
  -h, --help   help for generate
```

### Options inherited from parent commands

```
  -c, --config string       The codeowners config
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza](pizza.md)	 - OpenSauced CLI
* [pizza generate codeowners](pizza_generate_codeowners.md)	 - Generate a CODEOWNERS file for a GitHub repository using a "~/.sauced.yaml" config
* [pizza generate config](pizza_generate_config.md)	 - Generates a ".sauced.yaml" config based on the current repository
* [pizza generate insight](pizza_generate_insight.md)	 - Generate an OpenSauced Contributor Insight based on GitHub logins in a CODEOWNERS file

