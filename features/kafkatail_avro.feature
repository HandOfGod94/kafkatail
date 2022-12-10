Feature: kafkatail --wire_format=avro test-topic

  Avro is pretty common binary format for kafka.
  The messages are encoded to binary format and then pushed to kafka.

  As a developer
  I want to see the content of avro messages in plaintext
  So that it helps me while debugging in kafka-based async system

  @no-kafka
  Scenario: kafkatail --wire_format=avro test-topic
    It requires `schema_file` flag to be set

    When I run `kafkatail --wire_format=avro --bootstrap_servers=localhost:9093 test-topic`
    Then the exit status should not be 0
    And the output should contain:
    """
    Error: required flag(s) "schema_file" not set
    """

  Scenario: kafkatail -b=localhost:9093 --wire_format=avro --schema_path=../../testdata/starwars_avro_schemas --message_type=character test-topic

  Prints plain text message after decoding avro message

    Given "test-topic" is present on kafka-broker "localhost:9093" with 1 partition
    And I wait 2.0 seconds for a command to start up
    When I run `kafkatail -b=localhost:9093 --wire_format=avro --schema_path=../../testdata/starwars_avro_schemas --message_type=character test-topic` in background
    And starwars Human proto message is pushed to "test-topic"
    And I stop the command started last
    Then the output should contain 'homePlanet:'
    And the output should contain '"earth"'

