package main

type phase string

func (p phase) isTerminal() bool {
	switch p {
	case errorPhase, killedPhase, exitedPhase:
		return true
	}
	return false
}

const (
	creatingPhase phase = "creating"
	buildingPhase phase = "building"
	pullingPhase  phase = "pulling"
	startingPhase phase = "starting"
	runningPhase  phase = "running"
	livePhase     phase = "live"
	deadPhase     phase = "dead"
	readyPhase    phase = "ready"
	unreadyPhase  phase = "unready"
	exitedPhase   phase = "exited"
	killingPhase  phase = "killing"
	killedPhase   phase = "killed"
	errorPhase    phase = "error"
)
