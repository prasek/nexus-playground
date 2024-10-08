package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/temporalio/nexus-playground/service"
	"github.com/temporalio/nexus-playground/utils"
	"go.temporal.io/sdk/temporalnexus"
	"go.temporal.io/sdk/workflow"
)

func logServiceInfo(ctx context.Context, input service.Input, msg string, keyvals ...interface{}) {
	logger := temporalnexus.GetLogger(ctx)
	logger.Info(message(input, msg), keyvals...)
}

func logServiceError(ctx context.Context, input service.Input, msg string, err error, keyvals ...interface{}) error {
	logger := temporalnexus.GetLogger(ctx)
	logger.Error(errorDump(input, msg, err), keyvals...)
	return err
}

func logWorkflowInfo(ctx workflow.Context, input service.Input, msg string, keyvals ...interface{}) {
	logger := workflow.GetLogger(ctx)
	logger.Info(message(input, msg), keyvals...)
}

func logWorkflowError(ctx workflow.Context, input service.Input, msg string, err error, keyvals ...interface{}) error {
	logger := workflow.GetLogger(ctx)
	logger.Error(errorDump(input, msg, err), keyvals...)
	return err
}

func message(input service.Input, msg string) string {
	return utils.Message(input, msg)
}

func errorDump(input service.Input, msg string, err error) string {
	return utils.ErrorDump(input, msg, err)
}

func errorMessage(input service.Input, err error) string {
	return utils.ErrorMessage(input, err)
}

func newOutput(input service.Input, msg string) *service.Output {
	return &service.Output{Message: message(input, msg)}
}

func getFirstArg(input service.Input, errorMsg string) (string, error) {
	if len(input.Args) > 0 {
		return input.Args[0], nil
	}
	return "", errors.New(errorMsg)
}

func workflowID(input service.Input) string {
	return fmt.Sprintf("%s-%s", input.Operation, input.BusinessID)
}

func workflowIDWaitForSignal(txID string) string {
	return fmt.Sprintf("%s-%s", "async-op-workflow-wait-for-signal", txID)
}
