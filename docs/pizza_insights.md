## pizza insights

Gather insights about git contributors, repositories, users and pull requests

### Synopsis

Gather insights about git contributors, repositories, user and pull requests and display the results

```
pizza insights <command> [flags]
```

### Options

```
  -h, --help            help for insights
  -o, --output string   The formatting for command output. One of: (table, yaml, csv, json) (default "table")
```

### Options inherited from parent commands

```
  -c, --config string       The saucectl config (default ".sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza](pizza.md)	 - OpenSauced CLI
* [pizza insights contributors](pizza_insights_contributors.md)	 - Gather insights about contributors of indexed git repositories
* [pizza insights repositories](pizza_insights_repositories.md)	 - Gather insights about indexed git repositories
* [pizza insights user-contributions](pizza_insights_user-contributions.md)	 - Gather insights on individual contributors for given repo URLs

