## pizza completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	pizza completion fish | source

To load completions for every new session, execute once:

	pizza completion fish > ~/.config/fish/completions/pizza.fish

You will need to start a new shell for this setup to take effect.


```
pizza completion fish [flags]
```

### Options

```
  -h, --help              help for fish
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

