package spec

import(
  "fmt"
  "testing"
)

type MockLogger struct {
  Message string
}

func (mockLogger *MockLogger) Fatalf(format string, args ...any) {
  mockLogger.Message = fmt.Sprintf(format, args...)

  fmt.Print(fmt.Sprintf("\n%sFatalf called with message: %s%s\n\n", greenStart, mockLogger.Message, colorEnd))
}

func TestSpecToEqIntSuccess(testing *testing.T) {
  var mockLogger *MockLogger = &MockLogger{}
  var spec = Spec{Testing: testing, logger: mockLogger}
  var expected int = 55;
  var result int = 55;

  spec.Expect(expected).ToEq(result)
  if mockLogger.Message != "" {
    testing.Fatalf(fmt.Sprintf("%sExpected no error but got: { %s }%s", magentaStart, mockLogger.Message, colorEnd))
  }
}

func TestSpecToEqIntFailure(testing *testing.T) {
  var mockLogger *MockLogger = &MockLogger{}
  var spec = Spec{Testing: testing, logger: mockLogger}
  var expected int = 55;
  var result int = 65;
  var expectedErrorMessage = fmt.Sprintf("%sExpected: %d to equal: %d%s", redStart, expected, result, colorEnd)

  spec.Expect(expected).ToEq(result)
  if mockLogger.Message != expectedErrorMessage {
    testing.Fatalf(fmt.Sprintf("Expected error message: { %s }, but got { %s }", expectedErrorMessage, mockLogger.Message))
  }
}

func TestSpecToNotEqIntSuccess(testing *testing.T) {
  var mockLogger *MockLogger = &MockLogger{}
  var spec = Spec{Testing: testing, logger: mockLogger}
  var expected int = 55;
  var result int = 65;

  spec.Expect(expected).ToNotEq(result)
  if mockLogger.Message != "" {
    testing.Fatalf(fmt.Sprintf("Expected no error, but got: { %s }", mockLogger.Message))
  }
}

func TestSpecToNotEqIntFailure(testing *testing.T) {
  var mockLogger *MockLogger = &MockLogger{}
  var spec = Spec{Testing: testing, logger: mockLogger}
  var expected int = 55;
  var result int = 55;
  var expectedErrorMessage = fmt.Sprintf("%sExpected: %d to not equal: %d%s", redStart, expected, result, colorEnd)

  spec.Expect(expected).ToNotEq(result)

  if mockLogger.Message != expectedErrorMessage {
    testing.Fatalf(fmt.Sprintf("Expected error message: { %s }, but got { %s }", expectedErrorMessage, mockLogger.Message))
  }
}
