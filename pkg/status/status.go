package status

type Stage int
type State int

const (
	Configure Stage = iota
	Build
)

const (
	Failed State = iota
	InProgress
	Success
)

type BuildStatus struct{
	BuildStage 	Stage
	BuildStatus State
	Message 	string
}
