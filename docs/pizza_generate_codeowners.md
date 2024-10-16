## pizza generate codeowners

Generate a CODEOWNERS file for a GitHub repository using a "~/.sauced.yaml" config

### Synopsis

Generates a CODEOWNERS file for a given git repository. The generated file specifies up to 3 owners for EVERY file in the git tree based on the number of lines touched in that specific file over the specified range of time.

Configuration:
The command requires a .sauced.yaml file for accurate attribution. This file maps 
commit email addresses to GitHub usernames. The command looks for this file in two locations:

1. In the root of the specified repository path
2. In the user's home directory (~/.sauced.yaml) if not found in the repository

If you run the command on a specific path, it will first look for .sauced.yaml in that 
path. If not found, it will fall back to ~/.sauced.yaml.

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

# Specify a custom output location for the CODEOWNERS file
pizza generate codeowners . --output-path /path/to/directory
		
```

### Options

```
  -h, --help                 help for codeowners
  -o, --output-path string   Directory to create the output file.
      --owners-style-file    Generate an agnostic OWNERS style file instead of CODEOWNERS.
  -r, --range int            The number of days to analyze commit history (default 90) (default 90)
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

