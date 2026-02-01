package types

type TaskStatus int

const (
	StatusSuccess TaskStatus = iota
	StatusFailed
	StatusSkipped
)

type TaskResult struct {
	Name   string
	Status TaskStatus
	Error  error
}

type OnTaskComplete func(TaskResult)
