package utils

import (
	"fmt"
	"log"
	"reflect"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func RegisterWorkflowStruct(w worker.Worker, wStruct interface{}) error {
	d := "\nWorkflows:"

	structValue := reflect.ValueOf(wStruct)
	structType := structValue.Type()
	count := 0
	for i := 0; i < structValue.NumMethod(); i++ {
		methodValue := structValue.Method(i)
		method := structType.Method(i)
		// skip private method
		if method.PkgPath != "" {
			continue
		}
		name := method.Name
		if err := validateFnFormat(method.Type); err != nil {
			log.Printf("ERROR: Invalid method %s of %s: %v", name, structType.Name(), err)
			return fmt.Errorf("method %s of %s: %w", name, structType.Name(), err)
		}

		d += fmt.Sprintf("\n - workflow: %s", name)
		w.RegisterWorkflowWithOptions(methodValue.Interface(), workflow.RegisterOptions{Name: name})

		count++
	}
	if count == 0 {
		return fmt.Errorf("no workflows (public methods) found at %v structure", structType.Name())
	}
	log.Println(d)
	return nil
}

// Validate function parameters.
func validateFnFormat(fnType reflect.Type) error {
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("expected a func as input but was %s", fnType.Kind())
	}
	if fnType.NumIn() < 1 {
		return fmt.Errorf(
			"expected at least one argument of type workflow.Context in function, found %d input arguments",
			fnType.NumIn(),
		)
	}
	if !isWorkflowContext(fnType.In(1)) {
		return fmt.Errorf("expected first argument to be workflow.Context but found %s", fnType.In(0))
	}

	// Return values
	// We expect either
	// 	<result>, error
	//	(or) just error
	if fnType.NumOut() < 1 || fnType.NumOut() > 2 {
		return fmt.Errorf(
			"expected function to return result, error or just error, but found %d return values", fnType.NumOut(),
		)
	}
	if fnType.NumOut() > 1 && !isValidResultType(fnType.Out(0)) {
		return fmt.Errorf(
			"expected function first return value to return valid type but found: %v", fnType.Out(0).Kind(),
		)
	}
	if !isError(fnType.Out(fnType.NumOut() - 1)) {
		return fmt.Errorf(
			"expected function second return value to return error but found %v", fnType.Out(fnType.NumOut()-1).Kind(),
		)
	}
	return nil
}

func isWorkflowContext(inType reflect.Type) bool {
	// NOTE: We don't expect any one to derive from workflow context.
	return inType == reflect.TypeOf((*workflow.Context)(nil)).Elem()
}

func isValidResultType(inType reflect.Type) bool {
	// https://golang.org/pkg/reflect/#Kind
	switch inType.Kind() {
	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return false
	}

	return true
}

func isError(inType reflect.Type) bool {
	errorElem := reflect.TypeOf((*error)(nil)).Elem()
	return inType != nil && inType.Implements(errorElem)
}
