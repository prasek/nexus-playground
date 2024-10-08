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
./cloud-run.sh starter help
```

### sync-op-ok command

```
./cloud-run.sh starter sync-op-ok
```

### async-op-workflow-ok command

```
./cloud-run.sh starter async-op-workflow-ok
```

### async-op-workflow-wait-for-cancel command
- complete by requesting cancelation in the UI (caller workflow or handler workflow):

```
./cloud-run.sh starter async-op-workflow-wait-for-cancel
```

## Signal commands

The signal commands are intended for use together.

### async-op-workflow-wait-for-signal command
```
./cloud-run.sh starter -tx-id <your tx ID> async-op-workflow-wait-for-signal
```

### sync-op-signal command

```
./cloud-run.sh starter -tx-id <your tx ID> sync-op-signal
```

## Error injection

### sync-op-error command

Get the available errors you can inject:
```
./cloud-run.sh starter sync-op-error help
```

### async-op-error command

Get the available errors you can inject:
```
./cloud-run.sh starter async-op-error help
```

### async-op-workflow-error command

Get the available errors you can inject:
```
./cloud-run.sh starter async-op-workflow-error help
```

## Timeout errors

You can force a Nexus Operation to timeout by using a long running command
with a short `-timeout` for example:

```
./cloud-run.sh starter -timeout 5 async-op-workflow-wait-for-cancel
```