package handler

import (
	"context"
	"time"

	"github.com/temporalio/nexus-playground/service"
	"github.com/temporalio/nexus-playground/utils"
	"go.temporal.io/sdk/activity"
)

var a Activities

type Activities struct{}

func (a *Activities) OK(ctx context.Context, input service.Input) (*service.Output, error) {
	return newOutput(input, "OK"), nil

	//return "", fmt.Errorf("other activity error")

	//return "", temporal.NewNonRetryableApplicationError("activity error: test 123", "activity_error_type", fmt.Errorf("other internal error"))
}

func (a *Activities) WaitForCancel(ctx context.Context, input service.Input) (*service.Output, error) {
	logger := activity.GetLogger(ctx)
	logger.Info(utils.Message(input, "a.WaitForCancel started"))
	for {
		select {
		case <-time.After(1 * time.Second):
			logger.Info(utils.Message(input, "a.WaitForCancel activity heartbeating..."))
			activity.RecordHeartbeat(ctx, "")
		case <-ctx.Done():
			logger.Info(utils.Message(input, "a.WaitForCancel activity canceled"))
			return newOutput(input, "canceled by Done"), nil
		}
	}
}
