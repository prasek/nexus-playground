# Nexus Playground

See [prasek/nexus-playground](https://github.com/prasek/nexus-playground) for initial setup.

Also note that [commands](#commands) are provided to exercise the `nexus-playground` service.

[Service: nexus-playground](https://github.com/prasek/nexus-playground/tree/main/handler/service.go)
 - operation: help
 - operation: sync-op-ok
 - operation: async-op-workflow-ok
 - operation: async-op-workflow-wait-for-signal
 - operation: sync-op-signal
 - operation: sync-op-wait-for-hour
 - operation: async-op-workflow-wait-for-cancel
 - operation: sync-op-error
 - operation: async-op-error
 - operation: async-op-workflow-error

Underlying Handler Workflows:
 - workflow: Error
 - workflow: OK
 - workflow: WaitForCancel
 - workflow: WaitForSignal

Common I/O types for all operaitons

- [Input](#input)
- [Output](#output)

#### Input
```go
type Input struct {
	Operation string //Nexus Opeation name
	BusinessID  string // really just a timestamp
	Args        []string // additional args if needed, can also be used as `tag`
}
```

#### Output
```go
type Output struct {
	Message string
}
```

## Commands

The following commands should be run in the `nexus-playground` project directory.

All commands support the following flags:
- `-timeout` - schedule-to-close timeout in seconds, default is 1 day; use to force a timeout with long running ops
- `-tx-id` - optional for most commands, used for handler wf id, default: timestamp

### help command

Prints available commands
```
./local-run.sh starter help
```

### sync-op-ok command

```
./local-run.sh starter sync-op-ok
```

### async-op-workflow-ok command

```
./local-run.sh starter async-op-workflow-ok
```

### async-op-workflow-wait-for-cancel command

caller workflow -> nexus op -> handler workflow -> long running activity (with heartbeats waiting for cancellation)

#### Cancel caller workflow from UI - handler workflow doesn't `WaitForCancellation`

Run the command below and then request cancelation of the **caller workflow** in the UI
```
./cloud-run.sh starter async-op-workflow-wait-for-cancel
```
- caller workflow will `RequestCancelNexusOperation`
- caller workflow waits for nexus op cancellation to be completed, which waits for the underlying handler workflow
- by default the underlying handler workflow will not `WaitForCancellation` for Activity execution, so cancellation will happen quickly as it won't wait for the activity cancellation.
- handler workflow will recieve and return a `(*internal.CanceledError)` from `activityFut.Get()`, so the handler workflow will be marked `Canceled`.
- caller workflow will recieve and return a `(*internal.CanceledError)` from `nexusFut.Get()`, so the caller workflow will be marked `Canceled`.

#### Cancel handler workflow from UI - handler workflow doesn't `WaitForCancellation`

Run the command below and then request cancelation of the **handler workflow** in the UI
```
./cloud-run.sh starter async-op-workflow-wait-for-cancel
```
- by default the underlying handler workflow will not `WaitForCancellation` for Activity execution, so cancellation will happen quickly as it won't wait for the activity cancellation.
- handler workflow will recieve and return a `(*internal.CanceledError)` from `activityFut.Get()`, so the handler workflow will be marked `Canceled`.
- caller workflow will recieve and return a `(*internal.CanceledError)` from `nexusFut.Get()`, so the caller workflow will be marked `Canceled`.
- this demonstrates handler -> caller canceled error propagation.

#### Handler workflow using `WaitForCancellation` for a long running activity

Using the `-handler-wait-for-cancellation` flag will set `WaitForCancellation: true` in the handler workflow activity options.
```
./cloud-run.sh starter -handler-wait-for-cancellation async-op-workflow-wait-for-cancel
```
- complete by requesting cancelation in the UI (caller workflow or handler workflow)
- caller workflow will wait for cancellation to be processed in the underlying handler activity, which returns success/completed.
- caller workflow and handler workflow will show `Completed`, not cancelled, with a result of "canceled by Done"

#### Cancel a Nexus Operation in a caller workflow using `workflow.WithCancel()`

The caller workflow may cancel with `workflow.WithCancel()`.
- Use `-caller-cancel <N seconds>` for the caller to cancel the Nexus Operation using, which by default it a try/cancel with an immediate exit.
- Use `-caller-wait-for-cancellation` to do a subsequent `fut.Get()` to wait for the operation to become `Canceled` or `Completed`.
```
./cloud-run.sh starter -caller-cancel 2 -caller-wait-for-cancellation async-op-workflow-wait-for-cancel
```
- handler workflow will recieve and return a `(*internal.CanceledError)` from `activityFut.Get()`, so the handler workflow will be marked `Canceled`.
- caller workflow will recieve and return a `(*internal.CanceledError)` from `nexusFut.Get()`, so the caller workflow will be marked `Canceled`.
- note: since the handler workflow is not using `WaitForCancellation: true` for it's long running activity, the activity will keep running until it detects "Error workflow execution already completed".

#### Try/Cancel a Nexus Operation in a caller workflow using `workflow.WithCancel()`

This is not recommended as cancelation isn't guaranteed if the parent exists before delivery is done.
For example, run the following command and see the handler workflow keeps running:
```
./cloud-run.sh starter -caller-cancel 2 async-op-workflow-wait-for-cancel
```
- caller workflow uses the `workflow.WithCancel()` func, which adds a `NexusOperationCancelRequested` event to the caller's workflow history and then returns immediately
- the nexus op cancellation request may be processed in a delayed fashion (or not at all) in some cases the underlying activity may take 10 minutes or longer to get a "canceled by Done" in the heartbeat.
- the handler workflow callback will also be unable to be delivered with a "request failed with: 404 Not Found" since the caller workflow has already completed.

#### Cancel Nexus Operation in a caller workflow with `workflow.WithCancel()` and both caller and handler `WaitForCancellation: true`

Waiting for cancellation for the entire chain is demonstrated with:
```
./cloud-run.sh starter -caller-cancel 2 -caller-wait-for-cancellation -handler-wait-for-cancellation async-op-workflow-wait-for-cancel
```
- caller workflow and handler workflow will show `Completed`, not cancelled, with a result of "canceled by Done" from underlying handler activity

## Signal commands

The signal commands are intended for use together.

### async-op-workflow-wait-for-signal command
```
./local-run.sh starter -tx-id <your tx ID> async-op-workflow-wait-for-signal
```

### sync-op-signal command

```
./local-run.sh starter -tx-id <your tx ID> sync-op-signal
```

## Error injection

### sync-op-error command

Get the available errors you can inject:
```
./local-run.sh starter sync-op-error help
```

### async-op-error command

Get the available errors you can inject:
```
./local-run.sh starter async-op-error help
```

### async-op-workflow-error command

Get the available errors you can inject:
```
./local-run.sh starter async-op-workflow-error help
```

## Timeout errors

You can force a Nexus Operation to timeout by using a long running command
with a short `-timeout` for example:

```
./local-run.sh starter -timeout 5 async-op-workflow-wait-for-cancel
```