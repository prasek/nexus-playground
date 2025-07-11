package caller

import (
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
	//	operation  string
	//	businessid string
	//	args       []string
	Timeout                   int64 //seconds
	Concurrency               int64 //num nexus ops to create
	BadInput                  bool  //should client pass a bad input
	CallerCancelTimeout       int64 //should caller cancel the Nexus op after N seconds
	CallerWaitForCancellation bool  //should caller wait for cancellation via fut.Get()
}

func CallerWorkflow(ctx workflow.Context, input CallerWorkflowInput) (string, error) {

	childCtx, cancelFunc := workflow.WithCancel(ctx)

	logWorkflowInfo(childCtx, input, "starting ...")

	c := workflow.NewNexusClient(input.Endpoint, input.Service)
	if input.BadInput {
		c = workflow.NewNexusClient("12441241241224", "332219028")
	}

	if input.Concurrency <= 1 {

		logWorkflowInfo(childCtx, input, "concurrency <= 1",
			"Concurrency", input.Concurrency,
			"BadInput", input.BadInput)

		fut := c.ExecuteOperation(childCtx,
			input.Operation,
			input.Input,
			workflow.NexusOperationOptions{
				ScheduleToCloseTimeout: time.Duration(input.Timeout) * time.Second,
			})

		// Optionally wait for the operation to be started. NexusOperationExecution will contain the operation ID in
		// case this operation is asynchronous.
		var exec workflow.NexusOperationExecution
		if err := fut.GetNexusOperationExecution().Get(childCtx, &exec); err != nil {

			logWorkflowError(childCtx, input, "GetNexusOperationExecution", err,
				"IsApplicationError", temporal.IsApplicationError(err),
				"IsCancelled", temporal.IsCanceledError(err),
				"IsTerminatedError", temporal.IsTerminatedError(err),
				"IsTimeoutError", temporal.IsTimeoutError(err),
			)

			return "", err
		}

		logWorkflowInfo(childCtx, input, "started",
			"OperationToken", exec.OperationToken)

		if input.CallerCancelTimeout > 0 {
			workflow.Sleep(childCtx, time.Duration(input.CallerCancelTimeout*int64(time.Second)))
			logWorkflowInfo(childCtx, input, "requesting cancellation via workflow.WithCancel() handler ...")
			cancelFunc() //from workflow.WithCancel()
			logWorkflowInfo(childCtx, input, "cancelFunc() invoked")

			if !input.CallerWaitForCancellation {
				logWorkflowInfo(childCtx, input, "CallerWaitForCancellation: false, returning immediately.")
				return "Nexus operation workflow.WithCancel() func invoked, returned immemediately from caller workflow func.", nil
			}
		}

		var res service.Output
		if err := fut.Get(childCtx, &res); err != nil {

			logWorkflowError(childCtx, input, "Get", err,
				"IsApplicationError", temporal.IsApplicationError(err),
				"IsCancelled", temporal.IsCanceledError(err),
				"IsTerminatedError", temporal.IsTerminatedError(err),
				"IsTimeoutError", temporal.IsTimeoutError(err),
			)

			return "", err
		}

		logWorkflowInfo(childCtx, input, "completed",
			"OperationToken", exec.OperationToken)

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

	logWorkflowInfo(childCtx, input, "concurrency > 1",
		"Concurrency", input.Concurrency)

	for i := 0; i < int(input.Concurrency); i++ {
		logWorkflowInfo(childCtx, input, "starting operation ...",
			"OpCount", i)

		fut := c.ExecuteOperation(childCtx,
			input.Operation,
			input.Input,
			workflow.NexusOperationOptions{
				ScheduleToCloseTimeout: time.Duration(input.Timeout) * time.Second,
			})

		results = append(results, fut)

	}

	for _, fut := range results {

		var res service.Output
		if err := fut.Get(childCtx, &res); err != nil {
			logWorkflowError(childCtx, input, "Get", err,
				"IsApplicationError", temporal.IsApplicationError(err),
				"IsCancelled", temporal.IsCanceledError(err),
				"IsTerminatedError", temporal.IsTerminatedError(err),
				"IsTimeoutError", temporal.IsTimeoutError(err),
			)

			return "", err
		}

		logWorkflowInfo(childCtx, input, "completed operation",
			"Message", res.Message)

	}
	return fmt.Sprintf("Successfully completed %d nexus operations", input.Concurrency), nil
}
