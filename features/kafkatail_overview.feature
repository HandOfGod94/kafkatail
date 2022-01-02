Feature: kafkatail [flags] topic

  Quick overview of what all things are possible with the cli

  As a developer
  I want to view messages coming to kafka topic
  So that it will help me while debuggin messaging system

  Background:
    Given "test-topic" is present on kafka-broker "localhost:9093"

  Scenario: topic is a required args
    When I run `kafkatail`
    Then the exit status should be 1
    And the output should contain:
    """
    Error: requires at least 1 arg(s)
    """

  Scenario: bootstrap_servers flag is required
    When I run `kafkatail test-topic`
    Then the exit status should be 1
    And the output should contain:
    """
    Error: required flag(s) "bootstrap_servers" not set
    """

