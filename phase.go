package main

type Phase string

const (
	creatingPhase Phase = "creating"
	excludedPhase Phase = "excluded"
	buildingPhase Phase = "building"
	runningPhase  Phase = "running"
	livePhase     Phase = "live"
	deadPhase     Phase = "dead"
	readyPhase    Phase = "ready"
	unreadyPhase  Phase = "unready"
	exitedPhase   Phase = "exited"
	errorPhase    Phase = "error"
)
