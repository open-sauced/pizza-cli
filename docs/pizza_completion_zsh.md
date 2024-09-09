## pizza completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(pizza completion zsh)

To load completions for every new session, execute once:

#### Linux:

	pizza completion zsh > "${fpath[1]}/_pizza"

#### macOS:

	pizza completion zsh > $(brew --prefix)/share/zsh/site-functions/_pizza

You will need to start a new shell for this setup to take effect.


```
pizza completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string       The codeowners config (default "~/.sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza completion](pizza_completion.md)	 - Generate the autocompletion script for the specified shell

