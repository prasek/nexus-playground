package caller

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/temporalio/nexus-playground/service"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueue = "my-caller-workflow-task-queue"
)

type CallerWorkflowInput struct {
	Endpoint string
	Service  string
	service.Input
	//	Operation  string
	//	BusinessID string
	//	Args       []string
	Timeout     int64 //seconds
	Concurrency int64 //num Nexus ops to create
	BadInput    bool  //should client pass a bad input

}

func CallerWorkflow(ctx workflow.Context, input CallerWorkflowInput) (string, error) {
	logWorkflowInfo(ctx, input, "starting ...")

	c := workflow.NewNexusClient(input.Endpoint, input.Service)
	if input.BadInput {
		c = workflow.NewNexusClient("12441241241224", "332219028")
	}

	if input.Concurrency <= 1 {

		logWorkflowInfo(ctx, input, "concurrency <= 1",
			"Concurrency", input.Concurrency,
			"BadInput", input.BadInput)

		fut := c.ExecuteOperation(ctx,
			input.Operation,
			input.Input,
			workflow.NexusOperationOptions{
				ScheduleToCloseTimeout: time.Duration(input.Timeout) * time.Second,
			})

		// Optionally wait for the operation to be started. NexusOperationExecution will contain the operation ID in
		// case this operation is asynchronous.
		var exec workflow.NexusOperationExecution
		if err := fut.GetNexusOperationExecution().Get(ctx, &exec); err != nil {
			logWorkflowError(ctx, input, "GetNexusOperationExecution", err,
				"IsApplicationError", temporal.IsApplicationError(err),
				"IsCancelled", temporal.IsCanceledError(err),
				"IsTerminatedError", temporal.IsTerminatedError(err),
				"IsTimeoutError", temporal.IsTimeoutError(err),
			)

			var nexusError *temporal.NexusOperationError
			if errors.As(err, &nexusError) {
				logWorkflowError(ctx, input, "NexusError", nexusError)
				logWorkflowError(ctx, input, "NexusError.Cause", nexusError.Cause)
			}
			return "", err
		}

		logWorkflowInfo(ctx, input, "started",
			"OperationID", exec.OperationID)

		var res service.Output
		if err := fut.Get(ctx, &res); err != nil {
			logWorkflowError(ctx, input, "Get", err,
				"IsApplicationError", temporal.IsApplicationError(err),
				"IsCancelled", temporal.IsCanceledError(err),
				"IsTerminatedError", temporal.IsTerminatedError(err),
				"IsTimeoutError", temporal.IsTimeoutError(err),
			)

			var nexusError *temporal.NexusOperationError
			if errors.As(err, &nexusError) {
				logWorkflowError(ctx, input, "NexusError", nexusError)
				logWorkflowError(ctx, input, "NexusError.Cause", nexusError.Cause)
			}
			return "", err
		}

		logWorkflowInfo(ctx, input, "completed",
			"OperationID", exec.OperationID)

		return res.Message, nil
	}

	//concurrency > 1

	match, err := regexp.MatchString("^sync-op-wait-for-five-seconds$", input.Operation)
	if err != nil {
		return "", temporal.NewNonRetryableApplicationError("concurrency op name regex validation error", "concurrency_op_name_regex", err)
	}
	if !match {
		return "", temporal.NewNonRetryableApplicationError(
			"concurrency > 1 is only for sync-op-wait-for-second",
			"concurrency_command_not_supported",
			fmt.Errorf("%s not supported with concurrency > 1", input.Operation))
	}

	var results []workflow.NexusOperationFuture

	logWorkflowInfo(ctx, input, "concurrency > 1",
		"Concurrency", input.Concurrency)

	for i := 0; i < int(input.Concurrency); i++ {
		logWorkflowInfo(ctx, input, "starting operation ...",
			"OpCount", i)

		fut := c.ExecuteOperation(ctx,
			input.Operation,
			input.Input,
			workflow.NexusOperationOptions{
				ScheduleToCloseTimeout: time.Duration(input.Timeout) * time.Second,
			})

		results = append(results, fut)

	}

	for _, fut := range results {

		var res service.Output
		if err := fut.Get(ctx, &res); err != nil {
			logWorkflowError(ctx, input, "Get", err,
				"IsApplicationError", temporal.IsApplicationError(err),
				"IsCancelled", temporal.IsCanceledError(err),
				"IsTerminatedError", temporal.IsTerminatedError(err),
				"IsTimeoutError", temporal.IsTimeoutError(err),
			)

			var nexusError *temporal.NexusOperationError
			if errors.As(err, &nexusError) {
				logWorkflowError(ctx, input, "NexusError", nexusError)
				logWorkflowError(ctx, input, "NexusError.Cause", nexusError.Cause)
			}
			return "", err
		}

		logWorkflowInfo(ctx, input, "completed operation",
			"Message", res.Message)

	}
	return fmt.Sprintf("Successfully completed %d nexus operations", input.Concurrency), nil
}
