## pizza generate codeowners

Generates a CODEOWNERS file for a given repository using a "~/.sauced.yaml" config

### Synopsis

WARNING: Proof of concept feature.

Generates a CODEOWNERS file for a git repository. This uses a ~/.sauced.yaml
configuration to attribute emails with given entities.

The generated file specifies up to 3 owners for EVERY file in the git tree based on the
number of lines touched in that specific file over the specified range of time.

```
pizza generate codeowners path/to/repo [flags]
```

### Options

```
  -h, --help                help for codeowners
      --owners-style-file   Whether to generate an agnostic OWNERS style file.
  -r, --range int           The number of days to lookback (default 90)
```

### Options inherited from parent commands

```
  -c, --config string       The saucectl config (default "~/.sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza generate](pizza_generate.md)	 - Generates something

