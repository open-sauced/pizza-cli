## pizza insights contributors

Gather insights about contributors of indexed git repositories

### Synopsis

Gather insights about contributors of indexed git repositories. This command will show new, recent, alumni, repeat contributors for each git repository

```
pizza insights contributors url... [flags]
```

### Options

```
  -f, --file string   Path to yaml file containing an array of git repository urls
  -h, --help          help for contributors
  -r, --range int     Number of days to look-back (7,30,90) (default 30)
```

### Options inherited from parent commands

```
  -c, --config string       The codeowners config
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
  -o, --output string       The formatting for command output. One of: (table, yaml, csv, json) (default "table")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza insights](pizza_insights.md)	 - Gather insights about git contributors, repositories, users and pull requests

