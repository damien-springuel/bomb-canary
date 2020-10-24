package messagebus

type Event struct {
	Party
}

func (e Event) Type() Type {
	return EventMessage
}

type PartyCreated struct {
	Event
}

type PlayerConnected struct {
	Event
	Player string
}

type PlayerDisconnected struct {
	Event
	Player string
}

type PlayerJoined struct {
	Event
	Player string
}

type AllegianceRevealed struct {
	Event
	AllegianceByPlayer map[string]Allegiance
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
	PlayerVotes  map[string]bool
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
	Success  bool
	Outcomes map[bool]int
}

type Allegiance string

const (
	Spy        Allegiance = "spy"
	Resistance Allegiance = "resistance"
)

type GameEnded struct {
	Event
	Winner Allegiance
}
