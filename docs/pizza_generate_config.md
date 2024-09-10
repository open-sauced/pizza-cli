## pizza generate config

Generates a ".sauced.yaml" config based on the current repository

### Synopsis

Generates a ".sauced.yaml" configuration file. The attribution of emails to given entities
is based on the repository this command is ran in.

```
pizza generate config path/to/repo [flags]
```

### Options

```
  -h, --help                       help for config
  -i, --interactive                Whether to be interactive
  -o, --output-path .sauced.yaml   Directory to create the .sauced.yaml file. (default "./")
```

### Options inherited from parent commands

```
  -c, --config string       The codeowners config (default "~/.sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza generate](pizza_generate.md)	 - Generates documentation and insights from your codebase

