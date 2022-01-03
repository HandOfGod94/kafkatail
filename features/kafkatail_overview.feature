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
    Given "test-topic" is present on kafka-broker "localhost:9093"
    When I run `kafkatail --bootstrap_servers=localhost:9093 test-topic` in background
    And "hello-world" message is pushed to "test-topic"
    And I stop the command started last
    Then the output should contain "hello-world"
