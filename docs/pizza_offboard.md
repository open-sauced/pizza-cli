## pizza offboard

CAUTION: Experimental Command. Removes users from the ".sauced.yaml" config and "CODEOWNERS" files.

### Synopsis

CAUTION: Experimental Command. Removes users from the \".sauced.yaml\" config and \"CODEOWNERS\" files.
Requires the users' name OR email.

```
pizza offboard <username/email> [flags]
```

### Options

```
  -h, --help          help for offboard
  -p, --path string   the path to the repository (required)
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

