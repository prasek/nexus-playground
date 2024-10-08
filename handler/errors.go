package handler

import (
	"errors"
	"fmt"
	"sort"

	"github.com/nexus-rpc/sdk-go/nexus"
	"github.com/temporalio/nexus-playground/service"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/temporal"
)

/*
simulatedErrors can be injected with one of the following operations:
- sync-op-error - into a NewSyncOperation
- async-op-error - into a NewWorkflowRunOperationWithOptions custom Handler
- async-op-workflow-error - into the underlying handler workflow

if you'd like to inject different errors just add them to the list below
and they'll automatically be available in the 3 operations above.
*/

type errorFunc func(service.Input) error

var simulatedErrors = map[string]errorFunc{
	"fmt.Errorf": func(input service.Input) error {
		return fmt.Errorf(message(input, "unknown error"))
	},
	"ApplicationError": func(input service.Input) error {
		return temporal.NewApplicationError(message(input, "temporal app error"), "my temporal app error type")
	},
	"ApplicationErrorNonRetryable": func(input service.Input) error {
		return temporal.NewNonRetryableApplicationError(message(input, "temporal app error"), "my temporal app error type", errors.New("cause: unknown"))
	},
	"HandlerErrorTypeInternal": func(input service.Input) error {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeInternal, message(input, "nexus intenral error"))
	},
	"HandlerErrorTypeResourceExhausted": func(input service.Input) error {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeResourceExhausted, message(input, "nexus resource exhausted error"))
	},
	"HandlerErrorTypeNotImplemented": func(input service.Input) error {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeNotImplemented, message(input, "nexus not implemented error"))
	},
	"HandlerErrorTypeUnavailable": func(input service.Input) error {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeUnavailable, message(input, "nexus unavailable error"))
	},
	"UnsuccessfulOperationError": func(input service.Input) error {
		return &nexus.UnsuccessfulOperationError{
			State: nexus.OperationStateFailed,
			Failure: nexus.Failure{
				Message: message(input, "unsuccessful nexus op error"),
			}}
	},
	"HandlerErrorTypeBadRequest": func(input service.Input) error {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, message(input, "nexus bad request error"))
	},
	"HandlerErrorTypeUnauthenticated": func(input service.Input) error {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeUnauthenticated, message(input, "nexus unauthenticated error"))
	},
	"HandlerErrorTypeUnauthorized": func(input service.Input) error {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeUnauthorized, message(input, "nexus unauthorized error"))
	},
	"HandlerErrorTypeNotFound": func(input service.Input) error {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeNotFound, message(input, "nexus not found error"))
	},
	"serviceerror.NamespaceNotFound": func(input service.Input) error {
		return &serviceerror.NamespaceNotFound{
			Message:   message(input, "simulated serviceerror namespace not found"),
			Namespace: "serviceerror",
		}
	},
}

func newErrorFromType(input service.Input, errorType string) error {
	f := simulatedErrors[errorType]
	if f == nil {
		return nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, "%s is not a valid error type", errorType)
	}
	return f(input)

}

func getAvailableErrorTypes() string {
	keys := make([]string, 0)
	for errType, _ := range simulatedErrors {
		keys = append(keys, errType)
	}

	sort.Strings(keys)

	s := "available error types:"
	for _, errType := range keys {
		s += fmt.Sprintf("\n- %s", errType)
	}
	return s

}

