## pizza generate insight

Generate an OpenSauced Contributor Insight based on GitHub logins in a CODEOWNERS file

### Synopsis

Generate an OpenSauced Contributor Insight based on GitHub logins in a CODEOWNERS file
to get metrics and insights on those users.

The provided path must be a local git repo with a valid CODEOWNERS file and GitHub "@login"
for each codeowner.

After logging in, the generated Contributor Insight on OpenSauced will have insights on
active contributors, contributon velocity, and more.

```
pizza generate insight path/to/repo/with/CODEOWNERS/file [flags]
```

### Examples

```
  # Use CODEOWNERS file in explicit directory
  $ pizza generate insight /path/to/repo

  # Use CODEOWNERS file in local directory
  $ pizza generate insight .
```

### Options

```
  -h, --help   help for insight
```

### Options inherited from parent commands

```
  -c, --config string       The codeowners config
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza generate](pizza_generate.md)	 - Generates documentation and insights from your codebase

