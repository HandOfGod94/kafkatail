Feature: kafkatail --wire_format=plaintext test-topic

  As a developer
  I want to view kafka messages encoded as plaintext for a topic
  So that I can see what's happening on my topic


  Scenario: kafkatail --wire_format=plaintext --bootstrap_servers=localhost:9093 test-topic
    Given "test-topic" is present on kafka-broker "localhost:9093" with 1 partition
    And I wait 2.0 seconds for a command to start up
    When I run `kafkatail --wire_format=plaintext --bootstrap_servers=localhost:9093 test-topic` in background
    And "hello-world" message is pushed to "test-topic" on partition 0
    And I stop the command started last
    Then the output should contain "hello-world"
