package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/temporalio/nexus-playground/service"
)

func init() {
	spew.Config = spew.ConfigState{
		MaxDepth:                1,
		DisableMethods:          false,
		ContinueOnMethod:        true,
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
		Indent:                  "\t",
	}
}

func Message(input service.Input, msg string) string {
	return fmt.Sprintf("%s(%s): %s", input.Operation, strings.Join(input.Args, ";"), msg)
}

func ErrorMessage(input service.Input, err error) string {
	return Message(input, fmt.Sprintf("%s", err))
}

func ErrorDump(prefix string, err error) string {
	res := fmt.Sprintf("%s: %#v", prefix, err)
	for {
		res += fmt.Sprintf("\n -- CAUSE: %v", spew.Sprintf("%#v", err))
		err = errors.Unwrap(err)
		if err == nil {
			break
		}
	}

	return res
}

func ErrorDumpMessage(input service.Input, msg string, err interface{}) string {
	return ErrorDump(msg, err.(error))
}
