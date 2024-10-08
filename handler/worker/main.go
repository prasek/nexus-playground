package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/nexus-playground/handler"
	"github.com/temporalio/nexus-playground/options"
	"github.com/temporalio/nexus-playground/utils"
)

const (
	taskQueue = "my-handler-task-queue"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	args, err := options.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(args.ClientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, taskQueue, worker.Options{})

	nexusService, err := handler.MyService.GetNexusService()
	if err != nil {
		log.Fatalln("Unable to register operations", err)
	}

	log.Print(handler.MyService.Describe())

	w.RegisterNexusService(nexusService)
	utils.RegisterWorkflowStruct(w, &handler.Workflow{})
	w.RegisterActivity(&handler.Activities{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
