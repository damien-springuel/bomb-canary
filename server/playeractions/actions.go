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
