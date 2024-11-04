package gamehub

import (
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type handler func(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message)

type messageDispatcher interface {
	Dispatch(m messagebus.Message)
}

type gameHub struct {
	messageDispatcher   messageDispatcher
	allegianceGenerator gamerules.AllegianceGenerator
	game                gamerules.Game
}

func New(messageDispatcher messageDispatcher, allegianceGenerator gamerules.AllegianceGenerator) *gameHub {
	return &gameHub{
		messageDispatcher:   messageDispatcher,
		allegianceGenerator: allegianceGenerator,
		game:                gamerules.NewGame(),
	}
}

func (s *gameHub) Consume(m messagebus.Message) {
	var handler handler
	switch m.(type) {
	case messagebus.JoinParty:
		handler = s.handleJoinPartyCommand
	case messagebus.StartGame:
		handler = s.handleStartGameCommand
	case messagebus.LeaderSelectsMember:
		handler = s.handleLeaderSelectsMember
	case messagebus.LeaderDeselectsMember:
		handler = s.handleLeaderDeselectsMember
	case messagebus.LeaderConfirmsTeamSelection:
		handler = s.handleLeaderConfirmsTeamSelection
	case messagebus.ApproveTeam:
		handler = s.handleApproveTeam
	case messagebus.RejectTeam:
		handler = s.handleRejectTeam
	case messagebus.SucceedMission:
		handler = s.handleSucceedMission
	case messagebus.FailMission:
		handler = s.handleFailMission
	default:
		return
	}

	updatedGame, messagesToDispatch := handler(s.game, m)

	s.game = updatedGame
	for _, messageToDispatch := range messagesToDispatch {
		s.messageDispatcher.Dispatch(messageToDispatch)
	}
}

func (s gameHub) handleJoinPartyCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	joinPartyCommand := message.(messagebus.JoinParty)
	updatedGame, err := currentGame.AddPlayer(joinPartyCommand.Player)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			messagebus.PlayerJoined{
				Player: joinPartyCommand.Player,
			},
		)
	}
	return
}

func (s gameHub) handleStartGameCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	updatedGame, playerAllegiancesByName, missionRequirementsByMission, err := currentGame.Start(s.allegianceGenerator)

	if err == nil {

		missionRequirements := make([]messagebus.MissionRequirement, len(missionRequirementsByMission))
		for i := range missionRequirements {
			requirement := missionRequirementsByMission[gamerules.Mission(i+1)]
			missionRequirements[i] = messagebus.MissionRequirement{
				NbPeopleOnMission:        requirement.NbOfPeopleToGo,
				NbFailuresRequiredToFail: requirement.NbFailuresRequiredToFailMission,
			}
		}

		messagesToDispatch = append(messagesToDispatch,
			messagebus.GameStarted{
				MissionRequirements: missionRequirements,
			},
		)

		allegiances := make(map[string]messagebus.Allegiance)
		for name, allegiance := range playerAllegiancesByName {
			allegiances[name] = messagebus.Allegiance(allegiance)
		}

		messagesToDispatch = append(messagesToDispatch,
			messagebus.AllegianceRevealed{
				AllegianceByPlayer: allegiances,
			},
		)
		messagesToDispatch = append(messagesToDispatch,
			messagebus.LeaderStartedToSelectMembers{
				Leader: updatedGame.Leader(),
			},
		)
	}
	return
}

func (s gameHub) handleLeaderSelectsMember(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderSelectsMemberCommand := message.(messagebus.LeaderSelectsMember)

	if leaderSelectsMemberCommand.Leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderSelectsMember(leaderSelectsMemberCommand.MemberToSelect)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			messagebus.LeaderSelectedMember{
				SelectedMember: leaderSelectsMemberCommand.MemberToSelect,
			},
		)
	}
	return
}

func (s gameHub) handleLeaderDeselectsMember(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderDeselectsMemberCommand := message.(messagebus.LeaderDeselectsMember)

	if leaderDeselectsMemberCommand.Leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderDeselectsMember(leaderDeselectsMemberCommand.MemberToDeselect)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			messagebus.LeaderDeselectedMember{
				DeselectedMember: leaderDeselectsMemberCommand.MemberToDeselect,
			},
		)
	}
	return
}

func (s gameHub) handleLeaderConfirmsTeamSelection(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderConfirmsTeamSelectionCommand := message.(messagebus.LeaderConfirmsTeamSelection)

	if leaderConfirmsTeamSelectionCommand.Leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderConfirmsTeamSelection()

	if err == nil {
		messagesToDispatch = append(messagesToDispatch, messagebus.LeaderConfirmedSelection{})
	}
	return
}

func (s gameHub) handleApproveTeam(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	approveTeamCommand := message.(messagebus.ApproveTeam)
	updatedGame, resultingVote, err := currentGame.ApproveTeamBy(approveTeamCommand.Player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		messagebus.PlayerVotedOnTeam{
			Player:   approveTeamCommand.Player,
			Approved: true,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonVoteOutgoingMessages(updatedGame, resultingVote)...)

	return
}

func (s gameHub) handleRejectTeam(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	rejectTeamCommand := message.(messagebus.RejectTeam)
	updatedGame, resultingVotes, err := currentGame.RejectTeamBy(rejectTeamCommand.Player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		messagebus.PlayerVotedOnTeam{
			Player:   rejectTeamCommand.Player,
			Approved: false,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonVoteOutgoingMessages(updatedGame, resultingVotes)...)

	return
}

func commonVoteOutgoingMessages(updatedGame gamerules.Game, resultingVote map[string]bool) []messagebus.Message {
	commonVoteMessages := []messagebus.Message{}
	if updatedGame.State() == gamerules.SelectingTeam {
		commonVoteMessages = append(commonVoteMessages,
			messagebus.AllPlayerVotedOnTeam{
				Approved:     false,
				VoteFailures: updatedGame.VoteFailures(),
				PlayerVotes:  resultingVote,
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			messagebus.LeaderStartedToSelectMembers{
				Leader: updatedGame.Leader(),
			},
		)
	} else if updatedGame.State() == gamerules.ConductingMission {
		commonVoteMessages = append(commonVoteMessages,
			messagebus.AllPlayerVotedOnTeam{
				Approved:    true,
				PlayerVotes: resultingVote,
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			messagebus.MissionStarted{},
		)
	} else if updatedGame.State() == gamerules.GameOver {
		commonVoteMessages = append(commonVoteMessages,
			messagebus.AllPlayerVotedOnTeam{
				Approved:     false,
				VoteFailures: updatedGame.VoteFailures(),
				PlayerVotes:  resultingVote,
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			messagebus.GameEnded{
				Winner: messagebus.Spy,
				Spies:  updatedGame.Spies(),
			},
		)
	}

	return commonVoteMessages
}

func (s gameHub) handleSucceedMission(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	succeedMissionCommand := message.(messagebus.SucceedMission)
	updatedGame, outcomes, err := currentGame.SucceedMissionBy(succeedMissionCommand.Player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		messagebus.PlayerWorkedOnMission{
			Player:  succeedMissionCommand.Player,
			Success: true,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonMissionOutgoingMessages(updatedGame, outcomes)...)

	return
}

func (s gameHub) handleFailMission(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	failMissionCommand := message.(messagebus.FailMission)
	updatedGame, outcomes, err := currentGame.FailMissionBy(failMissionCommand.Player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		messagebus.PlayerWorkedOnMission{
			Player:  failMissionCommand.Player,
			Success: false,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonMissionOutgoingMessages(updatedGame, outcomes)...)

	return
}

func commonMissionOutgoingMessages(updatedGame gamerules.Game, outcomes map[string]bool) []messagebus.Message {
	commonMissionMessages := []messagebus.Message{}

	if updatedGame.State() == gamerules.SelectingTeam {
		lastMissionSuccess := updatedGame.GetMissionResults()[updatedGame.CurrentMission()-1]
		talliedOutcomes := tallyOutcomes(outcomes)
		commonMissionMessages = append(commonMissionMessages,
			messagebus.MissionCompleted{
				Success:  lastMissionSuccess,
				Outcomes: talliedOutcomes,
			},
		)
		commonMissionMessages = append(commonMissionMessages,
			messagebus.LeaderStartedToSelectMembers{
				Leader: updatedGame.Leader(),
			},
		)
	} else if updatedGame.State() == gamerules.GameOver {
		lastMissionSuccess := updatedGame.GetMissionResults()[updatedGame.CurrentMission()]
		talliedOutcomes := tallyOutcomes(outcomes)
		commonMissionMessages = append(commonMissionMessages,
			messagebus.MissionCompleted{
				Success:  lastMissionSuccess,
				Outcomes: talliedOutcomes,
			},
		)

		commonMissionMessages = append(commonMissionMessages,
			messagebus.GameEnded{
				Winner: messagebus.Allegiance(updatedGame.Winner()),
				Spies:  updatedGame.Spies(),
			},
		)
	}

	return commonMissionMessages
}

func tallyOutcomes(outcomes map[string]bool) map[bool]int {
	results := make(map[bool]int)
	for _, outcome := range outcomes {
		results[outcome] += 1
	}
	return results
}
