# Kafkatail

:construction: Work In Progress :construction:  
:construction: Proceed with caution :construction:

- [Kafkatail](#kafkatail)
  - [Overview](#overview)
  - [Installation](#installation)
  - [QuickStart](#quickstart)
  - [Known Limitation](#known-limitation)
  - [Development](#development)

## Overview

Tail kafka logs on console from a topic having messages in any format: Plaintext, proto, ...

## Installation

> Currently it's not released yet, but soon will be available for multiple platforms  
> via idiomatic package managers

## QuickStart

Frequently used commands to start with
```sh
# check kafka version installed
kafkatail --version

# Get help
kafkatail --help # or kafkatail -h

# tail a topic for a kafka cluster
kafkatail --bootstrap_servers=localhost:9093 foo-topic # or kafkatail -b localhost:9093 foo-topic

# tail a topic with protobuf encoded messages from a kafka cluster
# message_type = type of `message` defined in `.proto` file which you want to decode in
kafkatail -b localhost:9093 --proto_file=foo.proto --include_paths=/usr/dir1,/usr/dir2 --message_type=Bar foo-topic

# consume from all paritions on a topic
# Note. Currently only works by registering a kafka-consumer-group
kafkatail -b localhost:9093 --group_id=mygroup foo-topic
```

## Known Limitation

* No/Limited support for consuming from multiple partitions at same time:
  When you don't specific `--partition` arg while starting `kafkatail`, it defaults to `partition 0`.  
  Although if you specify `--group_id`, it will tail from all the partition, but you will have an extra consumer group entry
  on kafka topic.
* No avro, thirft, <custom> decoding support
  Only supports `protobuf` and `plaintext` for now

## Development
* Minimum go version: `1.16`

```sh
# to run dev checks and create binary
make

# to run test
make test

# to run integration tests
make integration-tests

# to clean up integration teste env
make clean-integration
```
