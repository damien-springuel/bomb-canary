package clientstream

type clientEvent struct {
	PlayerConnected              *playerConnected              `json:",omitempty"`
	PlayerDisconnected           *playerDisconnected           `json:",omitempty"`
	PlayerJoined                 *playerJoined                 `json:",omitempty"`
	SpiesRevealed                *spiesRevealed                `json:",omitempty"`
	LeaderStartedToSelectMembers *leaderStartedToSelectMembers `json:",omitempty"`
	LeaderSelectedMember         *leaderSelectedMember         `json:",omitempty"`
	LeaderDeselectedMember       *leaderDeselectedMember       `json:",omitempty"`
	LeaderConfirmedSelection     *leaderConfirmedSelection     `json:",omitempty"`
	PlayerVotedOnTeam            *playerVotedOnTeam            `json:",omitempty"`
	AllPlayerVotedOnTeam         *allPlayerVotedOnTeam         `json:",omitempty"`
	MissionStarted               *missionStarted               `json:",omitempty"`
	PlayerWorkedOnMission        *playerWorkedOnMission        `json:",omitempty"`
	MissionCompleted             *missionCompleted             `json:",omitempty"`
	GameEnded                    *gameEnded                    `json:",omitempty"`
	EventsReplayEnded            *eventsReplayEnded            `json:",omitempty"`
}

type playerJoined struct {
	Name string
}

type playerConnected struct {
	Name string
}

type playerDisconnected struct {
	Name string
}

type spiesRevealed struct {
	Spies map[string]struct{} `json:",omitempty"`
}

type leaderStartedToSelectMembers struct {
	Leader string
}

type leaderSelectedMember struct {
	SelectedMember string
}

type leaderDeselectedMember struct {
	DeselectedMember string
}

type leaderConfirmedSelection struct{}

type playerVotedOnTeam struct {
	Player   string
	Approved *bool `json:",omitempty"`
}

type allPlayerVotedOnTeam struct {
	Approved     bool
	VoteFailures int
	PlayerVotes  map[string]bool
}

type missionStarted struct{}

type playerWorkedOnMission struct {
	Player  string
	Success *bool `json:",omitempty"`
}

type missionCompleted struct {
	Success bool
	NbFails int
}

type gameEnded struct {
	Winner string
}

type eventsReplayEnded struct{}
