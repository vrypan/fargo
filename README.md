`fargo` is a CLI interface to Farcaster written in Go.

# Install

## Generic
Clone the repo, and run `make bin`. Copy the generated binary `fargo` to a location in your `$PATH`.

## Homebrew

Fargo can be installed with Brew by using the command below.

`brew install vrypan/fargo/fargo`

# Use

```
fargo --help

A command line tool to interact with Farcaster

A command line tool to interact with Farcaster

Usage:
  fargo [command]

Available Commands:
  cache       Cache management
  completion  Generate the autocompletion script for the specified shell
  config      Get/Set fargo configuration parameters
  download    Download Farcaster-embedded URLs
  get         Get Farcaster data
  help        Help about any command
  inspect     Inspect a cast
  send
  version     Get the current version.

Flags:
  -h, --help   help for fargo

Use "fargo [command] --help" for more information about a command.
```

The most interesting command is `fargo get`.

Try `fargo get --help` or `fargo get @vrypan.eth/0x3e9f6825dc23a14efb4c5d71723f5bea2f89095f -e`.

The second one will also make you appreciate how much spam is suppressed by Warpcast.
