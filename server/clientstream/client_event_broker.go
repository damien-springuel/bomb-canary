package clientstream

import (
	"encoding/json"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type eventSender interface {
	Send(message []byte)
	SendToPlayer(name string, message []byte)
	SendToAllButPlayer(name string, message []byte)
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

func (c clientEventBroker) send(event clientEvent) {
	c.eventSender.Send(toJsonBytes(event))
}

func (c clientEventBroker) sendToPlayer(name string, event clientEvent) {
	c.eventSender.SendToPlayer(name, toJsonBytes(event))
}

func (c clientEventBroker) sendToAllButPlayer(name string, event clientEvent) {
	c.eventSender.SendToAllButPlayer(name, toJsonBytes(event))
}

func (c clientEventBroker) Consume(m messagebus.Message) {
	switch m := m.(type) {

	case messagebus.PlayerConnected:
		c.send(clientEvent{PlayerConnected: &playerConnected{Name: m.Player}})

	case messagebus.PlayerDisconnected:
		c.send(clientEvent{PlayerDisconnected: &playerDisconnected{Name: m.Player}})

	case messagebus.PlayerJoined:
		c.send(clientEvent{PlayerJoined: &playerJoined{Name: m.Player}})

	case messagebus.GameStarted:
		requirements := make([]missionRequirement, len(m.MissionRequirements))
		for i := range requirements {
			requirements[i] = missionRequirement{
				NbPeopleOnMission:        m.MissionRequirements[i].NbPeopleOnMission,
				NbFailuresRequiredToFail: m.MissionRequirements[i].NbFailuresRequiredToFail,
			}
		}
		c.send(clientEvent{GameStarted: &gameStarted{MissionRequirements: requirements}})

	case messagebus.AllegianceRevealed:
		spies := make(map[string]struct{})
		for name, allegiance := range m.AllegianceByPlayer {
			if allegiance == messagebus.Resistance {
				c.sendToPlayer(name, clientEvent{SpiesRevealed: &spiesRevealed{}})
			} else {
				spies[name] = struct{}{}
			}
		}
		for spy := range spies {
			c.sendToPlayer(spy, clientEvent{SpiesRevealed: &spiesRevealed{Spies: spies}})
		}

	case messagebus.LeaderStartedToSelectMembers:
		c.send(clientEvent{LeaderStartedToSelectMembers: &leaderStartedToSelectMembers{Leader: m.Leader}})

	case messagebus.LeaderSelectedMember:
		c.send(clientEvent{LeaderSelectedMember: &leaderSelectedMember{SelectedMember: m.SelectedMember}})

	case messagebus.LeaderDeselectedMember:
		c.send(clientEvent{LeaderDeselectedMember: &leaderDeselectedMember{DeselectedMember: m.DeselectedMember}})

	case messagebus.LeaderConfirmedSelection:
		c.send(clientEvent{LeaderConfirmedSelection: &leaderConfirmedSelection{}})

	case messagebus.PlayerVotedOnTeam:
		c.sendToPlayer(m.Player, clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: m.Player, Approved: boolP(m.Approved)}})
		c.sendToAllButPlayer(m.Player, clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: m.Player}})

	case messagebus.AllPlayerVotedOnTeam:
		c.send(clientEvent{AllPlayerVotedOnTeam: &allPlayerVotedOnTeam{Approved: m.Approved, VoteFailures: m.VoteFailures, PlayerVotes: m.PlayerVotes}})

	case messagebus.MissionStarted:
		c.send(clientEvent{MissionStarted: &missionStarted{}})

	case messagebus.PlayerWorkedOnMission:
		c.sendToPlayer(m.Player, clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: m.Player, Success: boolP(m.Success)}})
		c.sendToAllButPlayer(m.Player, clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: m.Player}})

	case messagebus.MissionCompleted:
		c.send(clientEvent{MissionCompleted: &missionCompleted{Success: m.Success, NbFails: m.Outcomes[false]}})

	case messagebus.GameEnded:
		c.send(clientEvent{GameEnded: &gameEnded{Winner: string(m.Winner), Spies: m.Spies}})
	}
}

func boolP(b bool) *bool {
	return &b
}
