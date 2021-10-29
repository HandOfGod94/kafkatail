# Kafkatail

[![build](https://github.com/HandOfGod94/kafkatail/actions/workflows/build.yml/badge.svg)](https://github.com/HandOfGod94/kafkatail/actions/workflows/build.yml)
[![test](https://github.com/HandOfGod94/kafkatail/actions/workflows/test.yml/badge.svg)](https://github.com/HandOfGod94/kafkatail/actions/workflows/test.yml)

:construction: **Work In Progress** :construction:  
:construction: **Proceed with caution** :construction:

## Overview

CLI app to tail kafka logs on console from any topic having messages in any format: Plaintext, proto, ...

**Contents**  
- [Installation](#installation)
- [QuickStart](#quickstart)
- [Usage](#usage)
- [Examples](#examples)
- [Known Limitation](#known-limitation)
- [Development](#development)

### Installation

Feel free to use any of the available methods

* Download executable binaries from [Release](https://github.com/HandOfGod94/kafkatail/releases) page
* Install via `go` toolchain
  + For Go 1.16 or higher
    ```sh
    go install github.com/handofgod94/kafkatail@v0.1.2
    ```

  + For go version <= 1.15
    ```sh
    GO111MODULE=on go get github.com/handofgod94/kafkatail@v0.1.2
    ```
    The go binary will be installed in `$GOPATH/bin`

> Soon will be available via package managers

### QuickStart

`$ kafkatail --version`
```
kafkatail version 0.1.3
```

### Usage

`$ kafkatail --help`
```
Tail kafka messages from any topic, of any wire format on console (plaintext, protobuf)

Usage:
  kafkatail [flags] topic

Flags:
  -b, --bootstrap_servers bootstrap_servers   list of kafka bootstrap_servers separated by comma
      --from_datetime string                  tail from specific past datetime in RFC3339 format (default "0001-01-01T00:00:00Z")
      --group_id group_id                     [Optional] kafka consumer group_id to be used for subscribing to topic
  -h, --help                                  help for kafkatail
      --include_paths include_paths           include_paths containing dependencies of proto. Required for `wire_format=proto`
      --message_type type                     proto message type to use for decoding . Required for `wire_format=proto`
      --offset int                            kafka offset to start consuming from. Possible Values: -1=latest, -2=earliest, n=nth offset (default -1)
      --partition int                         kafka partition to consume from
      --proto_file proto_file                 proto_file to be used for decoding kafka message. Required for `wire_format=proto`
  -v, --version                               version for kafkatail
      --wire_format wire_format[=plaintext]   Wire format of messages in topic (default plaintext)
```

### Examples
```sh
# tail messages from a topic
kafkatail --bootstrap_servers=localhost:9093 kafkatail-test

# tail proto messages from a topic
kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --proto_file=starwars.proto --include_paths="../testdata" --message_type=Human kafkatail-test-proto

# tail messages from an offset. Default: -1 (latest). For earliets, use offset=-2
kafkatail --bootstrap_servers=localhost:9093 --offset=12 kafkatail-test-base

# tail messages from specific time
kafkatail --bootstrap_servers=localhost:9093 --from_datetime=2021-06-28T15:04:23Z kafkatail-test-base

# tail messages from specific partition. Default: 0
kafkatail --bootstrap_servers=localhost:9093 --partition=5 kafkatail-test-base

# tail from multiple partitions, using group_id
kafkatail --bootstrap_servers=localhost:9093 --group_id=myfoo kafka-consume-gorup-id-int-test
```

### Known Limitation

* No/Limited support for consuming from multiple partitions at same time:
  When you don't specific `--partition` arg while starting `kafkatail`, it defaults to `partition 0`.  
  Although if you specify `--group_id`, it will tail from all the partition but there are caveats:
  * you will have an extra consumer group entry on kafka topic
  * `--offset` option doesn't work with `--gorup_id`
* No `avro`, `thirft`, `<custom>` decoding support yet.  
  Only supports `protobuf` and `plaintext` for now

### Development
* Minimum go version: `1.16`

```sh
# to run dev checks and create binary
make

# to run test
make test

# to run integration tests
make integration-tests

# to clean up integration test env
make clean-integration
```
