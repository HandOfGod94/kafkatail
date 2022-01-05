Feature: kafkatail --wire_format=proto test-topic

  As a developer
  I want to view kafka messages encoded as protobuf bytes for a topic
  So that I can see what's happening on my topic


  @no-kafka
  Scenario: kafkatail --wire_format=proto --bootstrap_servers=localhost:9093 test-topic
    It requires `include_paths`, `message_type` and `proto_file` flags

    When I run `kafkatail --wire_format=proto --bootstrap_servers=localhost:9093 test-topic`
    Then the exit status should not be 0
    And the output should contain:
    """
    Error: required flag(s) "include_paths", "message_type", "proto_file" not set
    """

  Scenario: kafkatail -b=localhost:9093 --wire_format=proto --include_paths=./testdata --proto_file=starwars.proto --message_type=Human test-topic

  Prints plain text message after decoding protobuf message

    Given "test-topic" is present on kafka-broker "localhost:9093" with 1 partition
    And I wait 2.0 seconds for a command to start up
    When I run `kafkatail -b=localhost:9093 --wire_format=proto --include_paths=../../testdata --proto_file=starwars.proto --message_type=Human test-topic` in background
    And starwars Human proto message is pushed to "test-topic"
    And I stop the command started last
    Then the output should contain 'homePlanet: "earth"'
