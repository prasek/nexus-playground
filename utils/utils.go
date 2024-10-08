package utils

import (
	"fmt"
	"strings"

	"github.com/temporalio/nexus-playground/service"
)

func Message(input service.Input, msg string) string {
	return fmt.Sprintf("%s(%s): %s", input.Operation, strings.Join(input.Args, ";"), msg)
}

func ErrorMessage(input service.Input, err error) string {
	return Message(input, fmt.Sprintf("%s", err))
}

func ErrorDump(input service.Input, msg string, err interface{}) string {
	return Message(input, fmt.Sprintf("%s: %#v", msg, err))
}
