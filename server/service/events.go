package service

type playerJoined struct {
	party
	user string
}

type leaderStartedToSelectMembers struct {
	party
	leader string
}

type leaderSelectedMember struct {
	party
	selectedMember string
}

type leaderDeselectedMember struct {
	party
	deselectedMember string
}

type leaderConfirmedSelection struct {
	party
}

type playerVotedOnTeam struct {
	party
	player   string
	approved bool
}

type allPlayerVotedOnTeam struct {
	party
	approved     bool
	voteFailures int
}

type missionStarted struct {
	party
}

type missionCompleted struct {
	party
	success bool
}

type playerWorkedOnMission struct {
	party
	player  string
	success bool
}

type allegiance string

const (
	Spy        allegiance = "Spy"
	Resistance allegiance = "Resistance"
)

type gameEnded struct {
	party
	winner allegiance
}
