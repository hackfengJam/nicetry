package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"

	"github.com/hackfengJam/nicetry/temporal/helloworld"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		// ID:        "hello_world_workflowID",
		TaskQueue: "hello-world",
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}

	weList := make([]client.WorkflowRun, 0)

	i := 0
	for {
		if i > 1 {
			break
		}
		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, helloworld.Workflow, "Temporal")
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		weList = append(weList, we)
		i++
	}

	// log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	for _, we := range weList {
		// Synchronously wait for the workflow completion.
		var result string
		err = we.Get(context.Background(), &result)
		if err != nil {
			log.Fatalln("Unable get workflow result", err)
		}
		log.Println("Workflow result:", result)
	}

}
