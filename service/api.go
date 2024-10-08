package service

const (
	MyServiceName   = "nexus-playground"
	CustomTxIDRegEx = "^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$"
)

type Input struct {
	Operation  string
	BusinessID string
	Args       []string
}

type Output struct {
	Message string
}
