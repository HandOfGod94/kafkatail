Feature: kafkatail version

    @no-kafka
    Scenario: Check version
    When I successfully run `kafkatail --version`
    Then the output should contain:
    """
    kafkatail version 0.2.0
    """
