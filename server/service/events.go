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

type playerWorkedOnMission struct {
	party
	player  string
	success bool
}