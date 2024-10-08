package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/nexus-rpc/sdk-go/nexus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporalnexus"
	"go.temporal.io/sdk/workflow"
)

type ServiceBuilder struct {
	nexusService *nexus.Service
	operations   []string
	errors       []string
}

func NewServiceBuilder(name string) *ServiceBuilder {
	log.Printf("Nexus service %s\n", name)
	nexusService := nexus.NewService(name)
	return &ServiceBuilder{
		nexusService: nexusService,
		operations:   make([]string, 0),
		errors:       make([]string, 0),
	}
}

func (sb *ServiceBuilder) registerOperation(op nexus.RegisterableOperation) {
	sb.operations = append(sb.operations, op.Name())
	err := sb.nexusService.Register(op)
	sb.handleError(err)
}

func (sb *ServiceBuilder) handleError(err error) {
	if err != nil {
		sb.errors = append(sb.errors, fmt.Sprintf("- %s", err))
	}
}

func (sb *ServiceBuilder) Operations() []string {
	return sb.operations
}

func (sb *ServiceBuilder) Describe() string {
	d := fmt.Sprintf("\nService: %s", sb.nexusService.Name)

	for _, v := range sb.operations {
		d += fmt.Sprintf("\n - operation: %s", v)
	}
	return d
}

func (sb *ServiceBuilder) error() error {
	if len(sb.errors) == 0 {
		return nil
	}
	return fmt.Errorf("service builder errors:\n%s", strings.Join(sb.errors, "\n"))
}

func (sb *ServiceBuilder) GetNexusService() (*nexus.Service, error) {
	return sb.nexusService, sb.error()
}

func NewSyncOperation[I any, O any](
	sb *ServiceBuilder,
	name string,
	handler func(context.Context, client.Client, I, nexus.StartOperationOptions) (O, error),
) {
	op := temporalnexus.NewSyncOperation(name, handler)
	sb.registerOperation(op)
}

func NewWorkflowRunOperation[I, O any](
	sb *ServiceBuilder,
	name string,
	workflow func(workflow.Context, I) (O, error),
	getOptions func(context.Context, I, nexus.StartOperationOptions) (client.StartWorkflowOptions, error),
) {
	op := temporalnexus.NewWorkflowRunOperation(name, workflow, getOptions)
	sb.registerOperation(op)
}

func NewWorkflowRunOperationWithOptions[I, O any](
	sb *ServiceBuilder,
	options temporalnexus.WorkflowRunOperationOptions[I, O],
) {
	op, err := temporalnexus.NewWorkflowRunOperationWithOptions(options)
	if err != nil {
		sb.handleError(err)
		return
	}
	sb.registerOperation(op)
}
