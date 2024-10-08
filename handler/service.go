package handler

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/nexus-rpc/sdk-go/nexus"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/temporalnexus"

	"github.com/temporalio/nexus-playground/service"
	"github.com/temporalio/nexus-playground/utils"
)

var MyService = BuildNexusService()

var topLevelHelpUsage = "starter op-name ... [tag]"

/*
  New Nexus Operations can be added to BuildNexusService() below
  and that is the ONLY change needed to try different scenarios.

  The starter and caller are generic, the help operation
  will automatically include any new operations, and
  the worker registation is also done automatically.

  So feel free to add more operations below and just restart your
  handler process to try different StartWorkflowOptions.

  If you want to inject different types of errors, see errors.go

*/

func BuildNexusService() *utils.ServiceBuilder {
	s := utils.NewServiceBuilder(service.MyServiceName)

	utils.NewSyncOperation(s,
		"help",
		func(ctx context.Context, c client.Client, input service.Input, options nexus.StartOperationOptions) (*service.Output, error) {
			logServiceInfo(ctx, input, "starting ...")
			return newOutput(input, topLevelHelpUsage), nil
		})

	utils.NewSyncOperation(s,
		"sync-op-ok",
		func(ctx context.Context, c client.Client, input service.Input, options nexus.StartOperationOptions) (*service.Output, error) {
			logServiceInfo(ctx, input, "starting ...")
			return newOutput(input, "OK"), nil
		})

	utils.NewWorkflowRunOperation(s,
		"async-op-workflow-ok",
		w.OK,
		func(ctx context.Context, input service.Input, options nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
			logServiceInfo(ctx, input, "starting ...")
			return client.StartWorkflowOptions{
				ID:                                       workflowID(input),
				WorkflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
				WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
				WorkflowExecutionErrorWhenAlreadyStarted: true,
			}, nil
		})

	utils.NewWorkflowRunOperationWithOptions(s, temporalnexus.WorkflowRunOperationOptions[service.Input, *service.Output]{
		Name: "async-op-workflow-wait-for-signal",
		Handler: func(ctx context.Context, input service.Input, options nexus.StartOperationOptions) (temporalnexus.WorkflowHandle[*service.Output], error) {
			logServiceInfo(ctx, input, "starting ...")
			err := validateSignalCommandInput(input)
			if err != nil {
				logServiceError(ctx, input, "validateSignalCommandInput", err)
				return nil, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, errorMessage(input, err))
			}

			wfOpts := client.StartWorkflowOptions{
				ID:                                       workflowIDWaitForSignal(input.BusinessID),
				WorkflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
				WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
				WorkflowExecutionErrorWhenAlreadyStarted: true,
			}

			logServiceInfo(ctx, input, "starting workflow ...",
				"Workflow ID", wfOpts.ID)

			handle, err := temporalnexus.ExecuteWorkflow(ctx, options, wfOpts, w.WaitForSignal, input)
			if err != nil {
				logServiceError(ctx, input, "ExecuteWorkflow", err)

				//must do this check and return a non-retryable error or it will 500 infinite retry by default today
				if temporal.IsWorkflowExecutionAlreadyStartedError(err) {
					return nil, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, errorMessage(input, err))
				}

				//TODO: wrap other errors with an UnsuccessfulOperationError, otherwise it will 500 retry indefinitely?
				// for example: serviceerrors.NotFound
				return nil, err
			}

			logServiceInfo(ctx, input, "started workflow",
				"Workflow ID", wfOpts.ID)
			return handle, nil
		}})

	// uses arg1:txID from async-op-workflow-wait-for-signal
	utils.NewSyncOperation(s,
		"sync-op-signal",
		func(ctx context.Context, c client.Client, input service.Input, options nexus.StartOperationOptions) (*service.Output, error) {
			logServiceInfo(ctx, input, "starting ...")
			err := validateSignalCommandInput(input)
			if err != nil {
				logServiceError(ctx, input, "validateSignalCommandInput", err)
				return nil, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, errorMessage(input, err))
			}

			err = c.SignalWorkflow(ctx, workflowIDWaitForSignal(input.BusinessID), "", "done", DoneSignal{Done: true})
			if err != nil {
				logServiceError(ctx, input, "SignalWorkflow", err)
				return nil, err
			}

			return newOutput(input, "OK"), nil
		})

	//can be used in conjunction with caller `ScheduleToClose` timeout to simluate a timeout
	utils.NewSyncOperation(s,
		"sync-op-wait-for-hour",
		func(ctx context.Context, c client.Client, input service.Input, options nexus.StartOperationOptions) (*service.Output, error) {
			logServiceInfo(ctx, input, "starting ...")
			time.Sleep(1 * time.Hour)
			return newOutput(input, "OK"), nil
		})

	//use to try cancel workflow behavior - e.g. cancel the caller or handler workflow
	utils.NewWorkflowRunOperation(s,
		"async-op-workflow-wait-for-cancel",
		w.WaitForCancel,
		func(ctx context.Context, input service.Input, options nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
			logServiceInfo(ctx, input, "starting ...")
			return client.StartWorkflowOptions{
				ID:                                       workflowID(input),
				WorkflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
				WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
				WorkflowExecutionErrorWhenAlreadyStarted: true,
			}, nil
		})

	utils.NewSyncOperation(s,
		"sync-op-error",
		func(ctx context.Context, c client.Client, input service.Input, options nexus.StartOperationOptions) (*service.Output, error) {
			logServiceInfo(ctx, input, "starting ...")
			errorType, err := getFirstArg(input, "error type not found; usage: starter sync-op-error <error type>")
			if err != nil {
				logServiceError(ctx, input, "getFirstArg", err)
				return nil, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, "%s\n%s", errorMessage(input, err), getAvailableErrorTypes())
			}

			if errorType == "help" {
				return nil, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, getAvailableErrorTypes())
			}

			return nil, logServiceError(ctx, input, "Result", newErrorFromType(input, errorType))

		})

	utils.NewWorkflowRunOperationWithOptions(s, temporalnexus.WorkflowRunOperationOptions[service.Input, *service.Output]{
		Name: "async-op-error",
		Handler: func(ctx context.Context, input service.Input, options nexus.StartOperationOptions) (temporalnexus.WorkflowHandle[*service.Output], error) {
			logServiceInfo(ctx, input, "starting ...")
			errorType, err := getFirstArg(input, "error type not found, usage: starter async-op-error <error type>")
			if err != nil {
				logServiceError(ctx, input, "getFirstArg", err)
				return nil, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, "%s\n%s", errorMessage(input, err), getAvailableErrorTypes())
			}

			if errorType == "help" {
				return nil, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, getAvailableErrorTypes())
			}

			return nil, logServiceError(ctx, input, "Result", newErrorFromType(input, errorType))
		}})

	utils.NewWorkflowRunOperation(s,
		"async-op-workflow-error",
		w.Error,
		func(ctx context.Context, input service.Input, options nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
			logServiceInfo(ctx, input, "starting ...")
			return client.StartWorkflowOptions{
				ID:                                       workflowID(input),
				WorkflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
				WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
				WorkflowExecutionErrorWhenAlreadyStarted: true,
			}, nil
		})

	topLevelHelpUsage += "\navailable commands:"
	for _, opName := range s.Operations() {
		topLevelHelpUsage += fmt.Sprintf("\n- starter %s", opName)
	}

	return s
}

func validateSignalCommandInput(input service.Input) error {
	match, err := regexp.MatchString(service.CustomTxIDRegEx, input.BusinessID)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("-tx-id does not match the regex: '%s'; usage: starter -tx-id <your tx id>", service.CustomTxIDRegEx)
	}
	return nil
}
