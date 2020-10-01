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

func (a actionService) StartGame(code string) {
	a.messageDispatcher.Dispatch(messagebus.StartGame{Command: messagebus.Command{Party: messagebus.Party{Code: code}}})
}

func (a actionService) LeaderSelectsMember(code string, leader string, member string) {
	a.messageDispatcher.Dispatch(
		messagebus.LeaderSelectsMember{
			Command:        messagebus.Command{Party: messagebus.Party{Code: code}},
			Leader:         leader,
			MemberToSelect: member,
		},
	)
}

func (a actionService) LeaderDeselectsMember(code string, leader string, member string) {
	a.messageDispatcher.Dispatch(
		messagebus.LeaderDeselectsMember{
			Command:          messagebus.Command{Party: messagebus.Party{Code: code}},
			Leader:           leader,
			MemberToDeselect: member,
		},
	)
}

func (a actionService) LeaderConfirmsTeam(code string, leader string) {
	a.messageDispatcher.Dispatch(
		messagebus.LeaderConfirmsTeamSelection{
			Command: messagebus.Command{Party: messagebus.Party{Code: code}},
			Leader:  leader,
		},
	)
}

func (a actionService) ApproveTeam(code string, player string) {
	a.messageDispatcher.Dispatch(
		messagebus.ApproveTeam{
			Command: messagebus.Command{Party: messagebus.Party{Code: code}},
			Player:  player,
		},
	)
}

func (a actionService) RejectTeam(code string, player string) {
	a.messageDispatcher.Dispatch(
		messagebus.RejectTeam{
			Command: messagebus.Command{Party: messagebus.Party{Code: code}},
			Player:  player,
		},
	)
}

func (a actionService) SucceedMission(code string, player string) {
	a.messageDispatcher.Dispatch(
		messagebus.SucceedMission{
			Command: messagebus.Command{Party: messagebus.Party{Code: code}},
			Player:  player,
		},
	)
}

func (a actionService) FailMission(code string, player string) {
	a.messageDispatcher.Dispatch(
		messagebus.FailMission{
			Command: messagebus.Command{Party: messagebus.Party{Code: code}},
			Player:  player,
		},
	)
}
