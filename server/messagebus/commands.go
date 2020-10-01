package messagebus

type Party struct {
	Code string
}

func (p Party) GetPartyCode() string {
	return p.Code
}

type JoinParty struct {
	Party
	User string
}

type StartGame struct {
	Party
}

type LeaderSelectsMember struct {
	Party
	Leader         string
	MemberToSelect string
}

type LeaderDeselectsMember struct {
	Party
	Leader           string
	MemberToDeselect string
}

type LeaderConfirmsTeamSelection struct {
	Party
	Leader string
}

type ApproveTeam struct {
	Party
	Player string
}

type RejectTeam struct {
	Party
	Player string
}

type SucceedMission struct {
	Party
	Player string
}

type FailMission struct {
	Party
	Player string
}
