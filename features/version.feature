Feature: kafkatail version

    Scenario: Check version
    When I successfully run `kafkatail --version`
    Then the output should contain:
    """
    kafkatail version dev
    """
