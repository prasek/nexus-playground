package caller

import (
	"github.com/temporalio/nexus-playground/utils"
	"go.temporal.io/sdk/workflow"
)

func logWorkflowInfo(ctx workflow.Context, input CallerWorkflowInput, msg string, keyvals ...interface{}) {
	logger := workflow.GetLogger(ctx)
	logger.Info(message(input, msg), keyvals...)
}

func logWorkflowError(ctx workflow.Context, input CallerWorkflowInput, msg string, err error, keyvals ...interface{}) error {
	logger := workflow.GetLogger(ctx)
	logger.Error(utils.ErrorDump(input.Input, msg, err), keyvals...)
	return err
}

func message(input CallerWorkflowInput, msg string) string {
	return utils.Message(input.Input, msg)
}

func errorMessage(input CallerWorkflowInput, err error) string {
	return utils.ErrorMessage(input.Input, err)
}
