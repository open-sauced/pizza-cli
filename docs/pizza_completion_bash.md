## pizza completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(pizza completion bash)

To load completions for every new session, execute once:

#### Linux:

	pizza completion bash > /etc/bash_completion.d/pizza

#### macOS:

	pizza completion bash > $(brew --prefix)/etc/bash_completion.d/pizza

You will need to start a new shell for this setup to take effect.


```
pizza completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string       The saucectl config (default "~/.sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza completion](pizza_completion.md)	 - Generate the autocompletion script for the specified shell

