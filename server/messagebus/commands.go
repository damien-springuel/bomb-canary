package messagebus

type Command struct {
	Party
}

func (c Command) Type() Type {
	return CommandMessage
}

type JoinParty struct {
	Command
	Player string
}

type StartGame struct {
	Command
}

type LeaderSelectsMember struct {
	Command
	Leader         string
	MemberToSelect string
}

type LeaderDeselectsMember struct {
	Command
	Leader           string
	MemberToDeselect string
}

type LeaderConfirmsTeamSelection struct {
	Command
	Leader string
}

type ApproveTeam struct {
	Command
	Player string
}

type RejectTeam struct {
	Command
	Player string
}

type SucceedMission struct {
	Command
	Player string
}

type FailMission struct {
	Command
	Player string
}
