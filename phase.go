package main

type phase string

const (
	creatingPhase phase = "creating"
	buildingPhase phase = "building"
	runningPhase  phase = "running"
	livePhase     phase = "live"
	deadPhase     phase = "dead"
	readyPhase    phase = "ready"
	unreadyPhase  phase = "unready"
	exitedPhase   phase = "exited"
	errorPhase    phase = "error"
)
