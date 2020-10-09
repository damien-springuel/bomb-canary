package clientstream

import (
	"encoding/json"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type eventSender interface {
	Send(code string, message []byte)
	SendToPlayer(code, name string, message []byte)
	SendToAllButPlayer(code, name string, message []byte)
}

type clientEventBroker struct {
	eventSender eventSender
}

func NewClientEventBroker(eventSender eventSender) clientEventBroker {
	return clientEventBroker{
		eventSender: eventSender,
	}
}

func toJsonBytes(event clientEvent) []byte {
	eventBytes, _ := json.Marshal(event)
	return eventBytes
}

func (c clientEventBroker) send(code string, event clientEvent) {
	c.eventSender.Send(code, toJsonBytes(event))
}

func (c clientEventBroker) sendToPlayer(code, name string, event clientEvent) {
	c.eventSender.SendToPlayer(code, name, toJsonBytes(event))
}

func (c clientEventBroker) sendToAllButPlayer(code, name string, event clientEvent) {
	c.eventSender.SendToAllButPlayer(code, name, toJsonBytes(event))
}

func (c clientEventBroker) Consume(m messagebus.Message) {
	code := m.GetPartyCode()
	switch m := m.(type) {
	case messagebus.PlayerConnected:
		c.send(code, clientEvent{PlayerConnected: &playerConnected{Name: m.Player}})

	case messagebus.PlayerDisconnected:
		c.send(code, clientEvent{PlayerDisconnected: &playerDisconnected{Name: m.Player}})

	case messagebus.PlayerJoined:
		c.send(code, clientEvent{PlayerJoined: &playerJoined{Name: m.Player}})

	case messagebus.AllegianceRevealed:
		spies := make(map[string]struct{})
		for name, allegiance := range m.AllegianceByPlayer {
			if allegiance == messagebus.Resistance {
				c.sendToPlayer(code, name, clientEvent{SpiesRevealed: &spiesRevealed{}})
			} else {
				spies[name] = struct{}{}
			}
		}
		for spy := range spies {
			c.sendToPlayer(code, spy, clientEvent{SpiesRevealed: &spiesRevealed{Spies: spies}})
		}

	case messagebus.LeaderStartedToSelectMembers:
		c.send(code, clientEvent{LeaderStartedToSelectMembers: &leaderStartedToSelectMembers{Leader: m.Leader}})

	case messagebus.LeaderSelectedMember:
		c.send(code, clientEvent{LeaderSelectedMember: &leaderSelectedMember{SelectedMember: m.SelectedMember}})

	case messagebus.LeaderDeselectedMember:
		c.send(code, clientEvent{LeaderDeselectedMember: &leaderDeselectedMember{DeselectedMember: m.DeselectedMember}})

	case messagebus.LeaderConfirmedSelection:
		c.send(code, clientEvent{LeaderConfirmedSelection: &leaderConfirmedSelection{}})

	case messagebus.PlayerVotedOnTeam:
		c.sendToPlayer(code, m.Player, clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: m.Player, Approved: boolP(m.Approved)}})
		c.sendToAllButPlayer(code, m.Player, clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: m.Player}})

	case messagebus.AllPlayerVotedOnTeam:
		c.send(code, clientEvent{AllPlayerVotedOnTeam: &allPlayerVotedOnTeam{Approved: m.Approved, VoteFailures: m.VoteFailures, PlayerVotes: m.PlayerVotes}})

	case messagebus.MissionStarted:
		c.send(code, clientEvent{MissionStarted: &missionStarted{}})

	case messagebus.PlayerWorkedOnMission:
		c.sendToPlayer(code, m.Player, clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: m.Player, Success: boolP(m.Success)}})
		c.sendToAllButPlayer(code, m.Player, clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: m.Player}})

	case messagebus.MissionCompleted:
		c.send(code, clientEvent{MissionCompleted: &missionCompleted{Success: m.Success, NbFails: m.Outcomes[false]}})

	case messagebus.GameEnded:
		c.send(code, clientEvent{GameEnded: &gameEnded{Winner: string(m.Winner)}})
	}
}

func boolP(b bool) *bool {
	return &b
}
