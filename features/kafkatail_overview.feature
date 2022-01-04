Feature: kafkatail [flags] topic

  Quick overview of what all things are possible with the cli

  As a developer
  I want to view messages coming to kafka topic
  So that it will help me while debugging messaging system

  @no-kafka
  Scenario: topic is a required args
    When I run `kafkatail`
    Then the exit status should be 1
    And the output should contain:
    """
    Error: requires at least 1 arg(s)
    """

  @no-kafka
  Scenario: bootstrap_servers flag is required
    When I run `kafkatail test-topic`
    Then the exit status should be 1
    And the output should contain:
    """
    Error: required flag(s) "bootstrap_servers" not set
    """

  Scenario: prints message when pushed to topic
    Given "test-topic" is present on kafka-broker "localhost:9093" with 1 partition
    And I wait 2.0 seconds for a command to start up
    When I run `kafkatail --bootstrap_servers=localhost:9093 test-topic` in background
    And "hello-world" message is pushed to "test-topic" on partition 0
    And I stop the command started last
    Then the output should contain "hello-world"

  Scenario: prints message from a specific partition with `--partition` flag
    Given "test-topic" is present on kafka-broker "localhost:9093" with 2 partition
    And I wait 2.0 seconds for a command to start up
    When I run `kafkatail --bootstrap_servers=localhost:9093 --partition=1 test-topic` in background
    And "hello-world" message is pushed to "test-topic" on partition 1
    And I stop the command started last
    Then the output should contain:
    """
    ====================Message====================
    ============Partition: 1, Offset: 0==========
    hello-world
    """
