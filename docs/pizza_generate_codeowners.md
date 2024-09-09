## pizza generate codeowners

Generate a CODEOWNERS file for a GitHub repository using a "~/.sauced.yaml" config

### Synopsis

Generates a CODEOWNERS file for a given git repository. This uses a ~/.sauced.yaml configuration to attribute emails with given entities.

The generated file specifies up to 3 owners for EVERY file in the git tree based on the number of lines touched in that specific file over the specified range of time.

```
pizza generate codeowners path/to/repo [flags]
```

### Examples

```

		# Generate CODEOWNERS file for the current directory
		pizza generate codeowners .

		# Generate CODEOWNERS file for a specific repository
		pizza generate codeowners /path/to/your/repo

		# Generate CODEOWNERS file analyzing the last 180 days
		pizza generate codeowners . --range 180

		# Generate an OWNERS style file instead of CODEOWNERS
		pizza generate codeowners . --owners-style-file

		# Specify a custom location for the .sauced.yaml file
		pizza generate codeowners . --config /path/to/.sauced.yaml
		
```

### Options

```
  -h, --help                help for codeowners
      --owners-style-file   Generate an agnostic OWNERS style file instead of CODEOWNERS.
  -r, --range int           The number of days to analyze commit history (default 90) (default 90)
```

### Options inherited from parent commands

```
  -c, --config string       The codeowners config (default "~/.sauced.yaml")
      --disable-telemetry   Disable sending telemetry data to OpenSauced
  -l, --log-level string    The logging level. Options: error, warn, info, debug (default "info")
      --tty-disable         Disable log stylization. Suitable for CI/CD and automation
```

### SEE ALSO

* [pizza generate](pizza_generate.md)	 - Generates documentation and insights from your codebase

