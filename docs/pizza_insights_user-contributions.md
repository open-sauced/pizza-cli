## pizza insights user-contributions

Gather insights on individual contributors for given repo URLs

### Synopsis

Gather insights on individual contributors given a list of repository URLs

```
pizza insights user-contributions url... [flags]
```

### Options

```
  -f, --file string     Path to yaml file containing an array of git repository urls
  -h, --help            help for user-contributions
  -p, --range int32     Number of days, used for query filtering (default 30)
  -s, --sort string     Sort user contributions by (total, commits, prs) (default "none")
  -u, --users strings   Inclusive comma separated list of GitHub usernames to filter for
```

### Options inherited from parent commands

```
  -c, --config string       The codeowners config (default ".sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
  -o, --output string       The formatting for command output. One of: (table, yaml, csv, json) (default "table")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza insights](pizza_insights.md)	 - Gather insights about git contributors, repositories, users and pull requests

