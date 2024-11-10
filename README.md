![fargo-logo](https://github.com/user-attachments/assets/9283719d-e434-4de3-b2e3-16f607e46e8f)

`fargo` is a Farcaster CLI written in Go.

# Install

## Generic
Clone the repo, and run `make local`. Copy the generated binary `fargo` to a location in your `$PATH`.

## Homebrew

Fargo can be installed with Brew by using the command below.

`brew install vrypan/fargo/fargo`

# Use

```
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
  post        Submit messages to the network
  version     Get the current version

Flags:
  -h, --help   help for fargo

Use "fargo [command] --help" for more information about a command.
```

The most interesting command is `fargo get`.

Try `fargo get --help` or `fargo get @vrypan.eth/0x3e9f6825dc23a14efb4c5d71723f5bea2f89095f -e`.

The second one will also make you appreciate how much spam is suppressed by Warpcast.

## Interacting with the network

To interact with the network (currently `fargo post cast`) you will need an app keypair (private/public).
At the moment, Fargo does not provide a way to generate these keys, but you can use a service like
https://www.castkeys.xyz to generate them.

Check out `fargo post cast --help` for more info.
