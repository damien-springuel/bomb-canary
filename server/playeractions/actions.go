package playeractions

import "github.com/damien-springuel/bomb-canary/server/messagebus"

type messageDispatcher interface {
	Dispatch(message messagebus.Message)
}

type actionService struct {
	messageDispatcher messageDispatcher
}

func NewActionService(messageDispatcher messageDispatcher) actionService {
	return actionService{
		messageDispatcher: messageDispatcher,
	}
}

func (a actionService) StartGame() {
	a.messageDispatcher.Dispatch(messagebus.StartGame{})
}

func (a actionService) LeaderSelectsMember(leader string, member string) {
	a.messageDispatcher.Dispatch(
		messagebus.LeaderSelectsMember{
			Leader:         leader,
			MemberToSelect: member,
		},
	)
}

func (a actionService) LeaderDeselectsMember(leader string, member string) {
	a.messageDispatcher.Dispatch(
		messagebus.LeaderDeselectsMember{
			Leader:           leader,
			MemberToDeselect: member,
		},
	)
}

func (a actionService) LeaderConfirmsTeam(leader string) {
	a.messageDispatcher.Dispatch(
		messagebus.LeaderConfirmsTeamSelection{
			Leader: leader,
		},
	)
}

func (a actionService) ApproveTeam(player string) {
	a.messageDispatcher.Dispatch(
		messagebus.ApproveTeam{
			Player: player,
		},
	)
}

func (a actionService) RejectTeam(player string) {
	a.messageDispatcher.Dispatch(
		messagebus.RejectTeam{
			Player: player,
		},
	)
}

func (a actionService) SucceedMission(player string) {
	a.messageDispatcher.Dispatch(
		messagebus.SucceedMission{
			Player: player,
		},
	)
}

func (a actionService) FailMission(player string) {
	a.messageDispatcher.Dispatch(
		messagebus.FailMission{
			Player: player,
		},
	)
}
