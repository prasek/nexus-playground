package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"

	"github.com/temporalio/nexus-playground/caller"
	"github.com/temporalio/nexus-playground/options"
	"github.com/temporalio/nexus-playground/service"
)

func main() {

	//TODO: add -timeout option for ScheduleToClose timeout in caller workflow
	//and pass to caller workflow as input
	args, err := options.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}

	c, err := client.Dial(args.ClientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	if len(args.Args) < 1 {

		log.Fatalln("Error: not enough args\nusage: starter nexus-op-name ...\nexample: starter sync-op-ok\nhelp: starter help")
	}
	timestamp := time.Now().Format("20060102150405")
	nexusOpName := args.Args[0]
	otherArgs := args.Args[1:]
	businessID := args.BusinessID
	if len(businessID) == 0 {
		match, err := regexp.MatchString("^.*signal.*$", nexusOpName)
		if err != nil {
			log.Fatalln("Error: -tx-id validation error", err)
		}
		if match {
			log.Fatalln("Error: -tx-id is required for signal commands")
		}

		businessID = timestamp
	} else {
		match, err := regexp.MatchString(service.CustomTxIDRegEx, businessID)
		if err != nil {
			log.Fatalln("Error: -tx-id validation error", err)
		}
		if !match {
			log.Fatalln("Error: -tx-id does not match the required regex: ^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$")
		}
	}
	callerWorkflowID := fmt.Sprintf("call-%s-%s", nexusOpName, timestamp)

	ctx := context.Background()
	workflowOptions := client.StartWorkflowOptions{
		ID:                                       callerWorkflowID,
		TaskQueue:                                caller.TaskQueue,
		WorkflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
		WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}

	input := caller.CallerWorkflowInput{
		Endpoint: *args.Endpoint,
		Service:  service.MyServiceName,
		Input: service.Input{
			Operation:  nexusOpName,
			BusinessID: businessID, // should use a real biz id, but for testing let's use this
			Args:       otherArgs,
		},
		Timeout:     args.Timeout, //seconds
		Concurrency: args.Concurrency,
		BadInput:    args.BadInput,
	}

	//fmt.Printf("\nOp name:\n- %s\nother args:\n- %s\n\n", nexusOpName, strings.Join(otherArgs, ","))

	wr, err := c.ExecuteWorkflow(ctx, workflowOptions, caller.CallerWorkflow, input)
	if err != nil {
		log.Fatalln("Caller execute error:", err)
	}
	log.Println("Caller workflow started", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())

	var result string
	err = wr.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Caller workflow error:", err)
	}
	log.Println("Caller workflow result:", result)
}
