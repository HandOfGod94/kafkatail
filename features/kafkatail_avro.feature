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


