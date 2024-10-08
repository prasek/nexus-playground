package handler

import (
	"errors"
	"time"

	"github.com/temporalio/nexus-playground/service"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var w Workflow

type Workflow struct{}

func (w *Workflow) OK(ctx workflow.Context, input service.Input) (*service.Output, error) {
	return newOutput(input, "OK"), nil
}

func (w *Workflow) WaitForCancel(ctx workflow.Context, input service.Input) (*service.Output, error) {
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Hour,
		},
	)

	result := service.Output{}
	err := workflow.ExecuteActivity(ctx,
		a.WaitForCancel,
		input,
	).Get(ctx, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (w *Workflow) WaitForSignal(ctx workflow.Context, input service.Input) (*service.Output, error) {
	var signal DoneSignal
	signalChan := workflow.GetSignalChannel(ctx, "done")
	signalChan.Receive(ctx, &signal)
	if signal.Done {
		return newOutput(input, "got done signal"), nil
	}
	return nil, errors.New(message(input, "got bad signal"))
}

type DoneSignal struct {
	Done bool
}

func (w *Workflow) Error(ctx workflow.Context, input service.Input) (*service.Output, error) {
	logWorkflowInfo(ctx, input, "starting ...")
	errorType, err := getFirstArg(input, "error type not found, usage: starter async-op-workflow-error <error type>")
	if err != nil {
		logWorkflowError(ctx, input, "getFirstArg", err)
		return nil, temporal.NewNonRetryableApplicationError(message(input, getAvailableErrorTypes()), "error_type_not_found", err)
	}

	if errorType == "help" {
		return newOutput(input, getAvailableErrorTypes()), nil
	}

	err = newErrorFromType(input, errorType)
	logWorkflowError(ctx, input, "Result", err)
	return nil, err
}
