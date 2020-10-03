package messagebus

type Event struct {
	Party
}

func (e Event) Type() Type {
	return EventMessage
}

type PlayerJoined struct {
	Event
	User string
}

type LeaderStartedToSelectMembers struct {
	Event
	Leader string
}

type LeaderSelectedMember struct {
	Event
	SelectedMember string
}

type LeaderDeselectedMember struct {
	Event
	DeselectedMember string
}

type LeaderConfirmedSelection struct {
	Event
}

type PlayerVotedOnTeam struct {
	Event
	Player   string
	Approved bool
}

type AllPlayerVotedOnTeam struct {
	Event
	Approved     bool
	VoteFailures int
}

type MissionStarted struct {
	Event
}

type PlayerWorkedOnMission struct {
	Event
	Player  string
	Success bool
}

type MissionCompleted struct {
	Event
	Success bool
}

type Allegiance string

const (
	Spy        Allegiance = "Spy"
	Resistance Allegiance = "Resistance"
)

type GameEnded struct {
	Event
	Winner Allegiance
}
