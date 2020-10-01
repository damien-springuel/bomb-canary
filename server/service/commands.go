package service

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

type startGame struct {
	Party
}

type leaderSelectsMember struct {
	Party
	leader         string
	memberToSelect string
}

type leaderDeselectsMember struct {
	Party
	leader           string
	memberToDeselect string
}

type leaderConfirmsTeamSelection struct {
	Party
	leader string
}

type approveTeam struct {
	Party
	player string
}

type rejectTeam struct {
	Party
	player string
}

type succeedMission struct {
	Party
	player string
}

type failMission struct {
	Party
	player string
}
