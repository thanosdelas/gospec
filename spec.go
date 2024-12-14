package spec

import (
  "fmt"
  "testing"
  "reflect"
  "strings"
  "runtime"
)

type TestingFatalfLogger interface {
  Fatalf(format string, args ...any)
}

type Spec struct {
  leftOperand any
  rightOperand any
  leftOperandError error
  rightOperandError error
  Testing *testing.T
  logger TestingFatalfLogger
}

const color bool = true
const greenStart string = "\033[92m"
const greenBoldStart string = "\033[1;32m"
const redStart string = "\033[31;1m"
const magentaStart string = "\033[95m"
const yellowStart string = "\033[33m"
const yellowStartBold string = "\033[33;1m"
const blueStartBold string = "\033[1;34m"
const colorEnd string = "\033[0m"

/**
 * Must ensure testing is provided
 * Must assign left operand
 */
func (spec *Spec) Expect(leftOperand any) *Spec {
  spec.ensureTestingIsProvided()
  spec.leftOperand = leftOperand

  return spec
}

/**
 * Must ensure testing is provided
 * Must assign left operand
 */
func (spec *Spec) ExpectNoError(leftOperand any) {
  spec.ensureTestingIsProvided()
  spec.leftOperand = leftOperand

  if spec.leftOperand == nil {
    return
  }

  e, isError := spec.leftOperand.(error)

  if !isError {
    printStack()
    spec.logger.Fatalf(errorMessage("Wrong error type provided. Expected to be a type of: *errors.errorString"))
  }

  printStack()
  spec.logger.Fatalf(errorMessage("Expected no error, but got: %v"), e.Error())
}

/**
 * Must assign right operand
 */
func (spec *Spec) ToEq(rightOperand any) *Spec {
  spec.rightOperand = rightOperand

  if !spec.verifytTypeMismatch() {
    printStack()

    spec.logger.Fatalf(
      errorMessage("Type mismatch between operands; cannot compare %v with %v"),
      reflect.TypeOf(spec.leftOperand),
      reflect.TypeOf(spec.rightOperand),
    )

    return spec
  }

  // Exclusive comparison for primitive data types.
  if spec.operandsAreErrors() {
    if spec.leftOperandError.Error() != spec.rightOperandError.Error() {
      printStack()
      spec.logger.Fatalf(errorMessage("Expected error: \"%v\" but got: \"%v\""), spec.leftOperandError.Error(), spec.rightOperandError.Error())
    }

    return spec
  }

  if spec.operandsAreStructs() {
    leftOperandMap := convertStructToMap(spec.leftOperand)
    rightOperandMap := convertStructToMap(spec.rightOperand)

    for key, _ := range(leftOperandMap) {
      if leftOperandMap[key] != rightOperandMap[key] {
        printStack()
        spec.logger.Fatalf(errorMessage("Expected %v.%v to equal %v but got %v"), reflect.TypeOf(spec.leftOperand), key, rightOperandMap[key], leftOperandMap[key])
      }
    }
  }

  // Normal comparison for primitive data types.
  if spec.leftOperand != rightOperand {
    printStack()
    spec.logger.Fatalf(errorMessage("Expected: %v to equal: %v"), spec.leftOperand, rightOperand)
  }

  return spec
}

/**
 * Must assign right operand
 */
func (spec *Spec) ToNotEq(rightOperand any) *Spec {
  spec.rightOperand = rightOperand

  if !spec.verifytTypeMismatch() {
    printStack()

    spec.logger.Fatalf(
      errorMessage("Type mismatch between operands; cannot compare %v with %v"),
      reflect.TypeOf(spec.leftOperand),
      reflect.TypeOf(spec.rightOperand),
    )

    return spec
  }

  // Exclusive comparison for primitive data types.
  if spec.operandsAreErrors() {
    if spec.leftOperandError.Error() == spec.rightOperandError.Error() {
      printStack()
      spec.logger.Fatalf(errorMessage("Expected error: \"%v\" but got: \"%v\""), spec.leftOperandError.Error(), spec.rightOperandError.Error())
    }

    return spec
  }

  if spec.operandsAreStructs() {
    leftOperandMap := convertStructToMap(spec.leftOperand)
    rightOperandMap := convertStructToMap(spec.rightOperand)

    for key, _ := range(leftOperandMap) {
      if leftOperandMap[key] == rightOperandMap[key] {
        printStack()
        spec.logger.Fatalf(errorMessage("Expected %v.%v to not equal %v but got %v"), reflect.TypeOf(spec.leftOperand), key, rightOperandMap[key], leftOperandMap[key])
      }
    }
  }

  // Normal comparison for primitive data types.
  if spec.leftOperand == rightOperand {
    printStack()
    spec.logger.Fatalf(errorMessage("Expected: %v to not equal: %v"), spec.leftOperand, rightOperand)
  }

  return spec
}

//
// Private
//

func (spec *Spec) verifytTypeMismatch() bool {
  if reflect.TypeOf(spec.leftOperand) != reflect.TypeOf(spec.rightOperand) {
    return false
  }

  return true
}

func (spec *Spec) operandsAreErrors() bool {
  var leftOperandIsError, rightOperandIsError bool

  spec.leftOperandError, leftOperandIsError = spec.leftOperand.(error)
  spec.rightOperandError, rightOperandIsError = spec.rightOperand.(error)

  return leftOperandIsError && rightOperandIsError
}

func (spec *Spec) operandsAreStructs() bool {
  leftOperandValueOf := reflect.ValueOf(spec.leftOperand)
  rightOperandValueOf := reflect.ValueOf(spec.rightOperand)

  return leftOperandValueOf.Kind().String() == "struct" && rightOperandValueOf.Kind().String() == "struct"
}

func convertStructToMap(data interface{}) map[string]interface{} {
  var structMap = make(map[string]interface{})

  valueOfData := reflect.ValueOf(data)
  typeOfData := reflect.TypeOf(data)

  for x := 0; x < typeOfData.NumField(); x++ {
    field := typeOfData.Field(x)

    structMap[field.Name] = valueOfData.Field(x).Interface()
  }

  return structMap
}

func errorMessage(message string) string {
  if !color {
    return message
  }

  return redStart + message + colorEnd
}

func printStack() {
  var buffer []uint8 = make([]byte, 1024)
  var stack int = runtime.Stack(buffer, false)
  var stackString string = string(buffer[:stack])
  var lines []string = strings.Split(stackString, "\n")

  for _, line := range lines {
    if strings.Contains(line, "Test") || strings.Contains(line, "_test.go") {
      fmt.Print(redStart)
      fmt.Println(strings.ReplaceAll(line, "\t", "  "))
      fmt.Print(colorEnd)
    } else {
      fmt.Print(yellowStart)
      fmt.Println(strings.ReplaceAll(line, "\t", "  "))
      fmt.Print(colorEnd)
    }
  }
}

func (spec *Spec) ensureTestingIsProvided() {
  if spec.Testing == nil {
    panic(errorMessage("You must provide *testing.T to Spec.Testing in order to use Spec (spec := spec.Spec{Testing: <*testing.T>})"))
  }

  if spec.logger == nil {
    spec.logger = spec.Testing
  }
}
