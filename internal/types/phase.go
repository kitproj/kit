package types

type Phase string

const (
	CreatingPhase Phase = "creating"
	ExcludedPhase Phase = "excluded"
	BuildingPhase Phase = "building"
	RunningPhase  Phase = "running"
	LivePhase     Phase = "live"
	DeadPhase     Phase = "dead"
	ReadyPhase    Phase = "ready"
	UnreadyPhase  Phase = "unready"
	ExitedPhase   Phase = "exited"
	ErrorPhase    Phase = "error"
)
