# gcsio

CLI for streaming I/O with Google Cloud Storage

## Installation

Navigate to the root `gcsio/` directory (where the `Makefile` is located) and run:

```
make install
```

## Usage

### `gcsio`

```
usage: gcsio [<flags>] <command> [<args> ...]

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  upload <dst>
    Streams stdin up to a GCS object

  cat [<flags>] <source>
    Streams an object from GCS to stdout
```

### `gcsio upload --help`

```
usage: gcsio upload <dst>

Streams stdin up to a GCS object

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Args:
  <dst>  Destination object URI
```

### `gcsio cat --help`

```
usage: gcsio cat [<flags>] <source>

Streams an object from GCS to stdout

Flags:
  --help           Show context-sensitive help (also try --help-long and --help-man).
  --no-decompress  Disable automatic stream decompression

Args:
  <source>  Source object URI or glob pattern
```

## Contributing

When contributing to this repository, please follow the steps below:

1. Fork the repository
1. Submit your patch in one commit, or a series of well-defined commits
1. Submit your pull request and make sure you reference the issue you are addressing
