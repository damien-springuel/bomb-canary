package service

type party struct {
	code string
}

func (p party) GetPartyCode() string {
	return p.code
}

type joinParty struct {
	party
	user string
}

type startGame struct {
	party
}

type leaderSelectsMember struct {
	party
	leader         string
	memberToSelect string
}

type leaderDeselectsMember struct {
	party
	leader           string
	memberToDeselect string
}

type leaderConfirmsTeamSelection struct {
	party
	leader string
}

type approveTeam struct {
	party
	player string
}

type rejectTeam struct {
	party
	player string
}

type succeedMission struct {
	party
	player string
}

type failMission struct {
	party
	player string
}