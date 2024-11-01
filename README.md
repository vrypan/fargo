`fargo` is a CLI interface to Farcaster written in Go.

# Install

## Generic
Clone the repo, and run `make bin`. Copy the generated binary `fargo` to a location in your `$PATH`.

## Homebrew

Fargo can be installed with Brew by using the command below.

brew install vrypan/fargo/fargo

# Use

```
fargo --help
```

The most interesting command is `fargo get`. 

Try `fargo get --help` or `fargo get @vrypan.eth/0x3e9f6825dc23a14efb4c5d71723f5bea2f89095f -e`.

The second one will also make you appreciate how much spam is suppressed by Warpcast.
