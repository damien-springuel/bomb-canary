package clientstream

type clientEvent struct {
	PlayerJoined                 *playerJoined                 `json:",omitempty"`
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
}

type playerJoined struct {
	Name string
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
}

type missionStarted struct{}

type playerWorkedOnMission struct {
	Player  string
	Success *bool `json:",omitempty"`
}

type missionCompleted struct {
	Success bool
}

type gameEnded struct {
	Winner string
}
