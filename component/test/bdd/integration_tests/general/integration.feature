Feature: Testing integration of bdd tester
  # This scenario should run and pass
  Scenario:
    Given This test will pass

  # This scenario should not run
  @exception
  Scenario:
    Given This test will fail
