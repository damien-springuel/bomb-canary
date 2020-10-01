package gamehub

type PlayerJoined struct {
	Party
	User string
}

type LeaderStartedToSelectMembers struct {
	Party
	Leader string
}

type LeaderSelectedMember struct {
	Party
	SelectedMember string
}

type LeaderDeselectedMember struct {
	Party
	DeselectedMember string
}

type LeaderConfirmedSelection struct {
	Party
}

type PlayerVotedOnTeam struct {
	Party
	Player   string
	Approved bool
}

type AllPlayerVotedOnTeam struct {
	Party
	Approved     bool
	VoteFailures int
}

type MissionStarted struct {
	Party
}

type MissionCompleted struct {
	Party
	success bool
}

type PlayerWorkedOnMission struct {
	Party
	Player  string
	Success bool
}

type Allegiance string

const (
	Spy        Allegiance = "Spy"
	Resistance Allegiance = "Resistance"
)

type GameEnded struct {
	Party
	Winner Allegiance
}
