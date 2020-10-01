package service

type playerJoined struct {
	Party
	user string
}

type leaderStartedToSelectMembers struct {
	Party
	leader string
}

type leaderSelectedMember struct {
	Party
	selectedMember string
}

type leaderDeselectedMember struct {
	Party
	deselectedMember string
}

type leaderConfirmedSelection struct {
	Party
}

type playerVotedOnTeam struct {
	Party
	player   string
	approved bool
}

type allPlayerVotedOnTeam struct {
	Party
	approved     bool
	voteFailures int
}

type missionStarted struct {
	Party
}

type missionCompleted struct {
	Party
	success bool
}

type playerWorkedOnMission struct {
	Party
	player  string
	success bool
}

type allegiance string

const (
	Spy        allegiance = "Spy"
	Resistance allegiance = "Resistance"
)

type gameEnded struct {
	Party
	winner allegiance
}
