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
