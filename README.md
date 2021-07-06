# Kafkatail

[![build](https://github.com/HandOfGod94/kafkatail/actions/workflows/build.yml/badge.svg)](https://github.com/HandOfGod94/kafkatail/actions/workflows/build.yml)
[![test](https://github.com/HandOfGod94/kafkatail/actions/workflows/test.yml/badge.svg)](https://github.com/HandOfGod94/kafkatail/actions/workflows/test.yml)

:construction: **Work In Progress** :construction:  
:construction: **Proceed with caution** :construction:

## Overview

CLI app to tail kafka logs on console from any topic having messages in any format: Plaintext, proto, ...

## Contents
- [Installation](#installation)
- [QuickStart](#quickstart)
- [Usage](#usage)
- [Known Limitation](#known-limitation)
- [Development](#development)

### Installation

Download binaries from `Release`

To install via go toolchain

* For Go 1.16 or higher
```sh
go install github.com/handofgod94/kafkatail@v0.1.1
```

* For go version <= 1.15
```sh
GO111MODULE=on go get github.com/handofgod94/kafkatail@v0.1.1
```

> Soon will be available via package managers

### QuickStart

Examples
```sh
# check kafka version installed
kafkatail --version

# Get help
kafkatail --help # or kafkatail -h

# tail a topic for a kafka cluster
kafkatail --bootstrap_servers=localhost:9093 foo-topic # or kafkatail -b localhost:9093 foo-topic

# tail a topic with protobuf encoded messages from a kafka cluster
# message_type = type of `message` to use for decoding. This must be defined in `.proto` file.
kafkatail -b localhost:9093 --proto_file=foo.proto --include_paths=/usr/dir1,/usr/dir2 --message_type=Bar foo-topic

# tail from specific offset
# Note. if you don't provide partition args, it defaults to Partition 0
kafkatail --bootstrap_servers=localhost:9093 --offeset 30 foo-topic 

# tail from specific datetime
# Note. if you don't provide partition args, it defaults to Partition 0
kafkatail --bootstrap_servers=localhost:9093 --from_datetime=2021-05-03T12:30:00Z foo-topic 

# consume from all paritions on a topic
# Note. Currently only works by registering a kafka-consumer-group
kafkatail -b localhost:9093 --group_id=mygroup foo-topic
```

### Usage

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

### Known Limitation

* No/Limited support for consuming from multiple partitions at same time:
  When you don't specific `--partition` arg while starting `kafkatail`, it defaults to `partition 0`.  
  Although if you specify `--group_id`, it will tail from all the partition, but you will have an extra consumer group entry
  on kafka topic.
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
