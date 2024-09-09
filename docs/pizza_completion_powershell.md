## pizza completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	pizza completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
pizza completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string       The saucectl config (default ".sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza completion](pizza_completion.md)	 - Generate the autocompletion script for the specified shell

