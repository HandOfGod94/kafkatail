Feature: kafkatail [flags] topic

  Quick overview of what all things are possible with the cli

  As a developer
  I want to view messages coming to kafka topic
  So that it will help me while debugging messaging system

  @no-kafka
  Scenario: kafkatail
    topic is required args

    When I run `kafkatail`
    Then the exit status should not be 0
    And the output should contain:
    """
    Error: requires at least 1 arg(s)
    """

  @no-kafka
  Scenario: kafatail test-topic
    `bootstrap_servers` flag is required

    When I run `kafkatail test-topic`
    Then the exit status should not be 0
    And the output should contain:
    """
    Error: required flag(s) "bootstrap_servers" not set
    """

  Scenario: kafkatail --bootstrap_servers=localhost:9093 test-topic
    prints message when pushed to topic

    Given "test-topic" is present on kafka-broker "localhost:9093" with 1 partition
    And I wait 2.0 seconds for a command to start up
    When I run `kafkatail --bootstrap_servers=localhost:9093 test-topic` in background
    And "hello-world" message is pushed to "test-topic" on partition 0
    And I stop the command started last
    Then the output should contain "hello-world"

  Scenario: kafkatail --bootstrap_servers=localhost:9093 --partition=1 test-topic
    prints message from a specific partition with `--partition` flag

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

  Scenario: kafkatail --bootstrap_servers=localhost:9093 --group_id=test-group test-topic
    prints messages from all partition with `--group_id` flag

    Given "test-topic" is present on kafka-broker "localhost:9093" with 2 partition
    And I wait 5.0 seconds for a command to start up
    When I run `kafkatail --bootstrap_servers=localhost:9093 --group_id=test-group test-topic` in background
    And "hello-world" message is pushed to "test-topic" on partition 1
    And "foo-bar" message is pushed to "test-topic" on partition 0
    Then the output should contain:
    """
    ====================Message====================
    ============Partition: 1, Offset: 0==========
    hello-world
    """
    And the output should contain:
    """
    ====================Message====================
    ============Partition: 0, Offset: 0==========
    foo-bar
    """

  Scenario: kafkatail -b=localhost:9093 test-topic
    using shorthand flag for `bootstrap_servers` will yield same result as with full flag name

    Given "test-topic" is present on kafka-broker "localhost:9093" with 1 partition
    And I wait 2.0 seconds for a command to start up
    When I run `kafkatail -b=localhost:9093 test-topic` in background
    And "hello-world" message is pushed to "test-topic" on partition 0
    And I stop the command started last
    Then the output should contain "hello-world"
