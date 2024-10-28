package gamehub

import (
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type handler func(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message)

type codeGenerator interface {
	GenerateCode() string
}

type messageDispatcher interface {
	Dispatch(m messagebus.Message)
}

type gameHub struct {
	codeGenerator       codeGenerator
	messageDispatcher   messageDispatcher
	allegianceGenerator gamerules.AllegianceGenerator
	games               games
}

func New(codeGenerator codeGenerator, messageDispatcher messageDispatcher, allegianceGenerator gamerules.AllegianceGenerator) gameHub {
	return gameHub{
		codeGenerator:       codeGenerator,
		messageDispatcher:   messageDispatcher,
		allegianceGenerator: allegianceGenerator,
		games:               newGames(),
	}
}

func (s gameHub) CreateParty() string {
	newCode := s.codeGenerator.GenerateCode()
	s.games.create(newCode)
	s.messageDispatcher.Dispatch(messagebus.PartyCreated{Event: messagebus.Event{Party: messagebus.Party{Code: newCode}}})
	return newCode
}

func (s gameHub) DoesPartyExist(code string) bool {
	_, exists := s.games.get(code)
	return exists
}

func (s gameHub) Consume(m messagebus.Message) {
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

	code := m.GetPartyCode()
	game, exists := s.games.get(code)
	if !exists {
		return
	}
	updatedGame, messagesToDispatch := handler(game, m)

	s.games.set(code, updatedGame)
	for _, messageToDispatch := range messagesToDispatch {
		s.messageDispatcher.Dispatch(messageToDispatch)
	}
}

func (s gameHub) getGameForPartyCode(code string) gamerules.Game {
	game, _ := s.games.get(code)
	return game
}

func (s gameHub) handleJoinPartyCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	joinPartyCommand := message.(messagebus.JoinParty)
	updatedGame, err := currentGame.AddPlayer(joinPartyCommand.Player)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			messagebus.PlayerJoined{
				Event:  messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
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
				Event:               messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
				MissionRequirements: missionRequirements,
			},
		)

		allegiances := make(map[string]messagebus.Allegiance)
		for name, allegiance := range playerAllegiancesByName {
			allegiances[name] = messagebus.Allegiance(allegiance)
		}

		messagesToDispatch = append(messagesToDispatch,
			messagebus.AllegianceRevealed{
				Event:              messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
				AllegianceByPlayer: allegiances,
			},
		)
		messagesToDispatch = append(messagesToDispatch,
			messagebus.LeaderStartedToSelectMembers{
				Event:  messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
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
				Event:          messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
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
				Event:            messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
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
		messagesToDispatch = append(messagesToDispatch,
			messagebus.LeaderConfirmedSelection{
				Event: messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
			},
		)
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
			Event:    messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
			Player:   approveTeamCommand.Player,
			Approved: true,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonVoteOutgoingMessages(updatedGame, message.GetPartyCode(), resultingVote)...)

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
			Event:    messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
			Player:   rejectTeamCommand.Player,
			Approved: false,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonVoteOutgoingMessages(updatedGame, message.GetPartyCode(), resultingVotes)...)

	return
}

func commonVoteOutgoingMessages(updatedGame gamerules.Game, code string, resultingVote map[string]bool) []messagebus.Message {
	commonVoteMessages := []messagebus.Message{}
	if updatedGame.State() == gamerules.SelectingTeam {
		commonVoteMessages = append(commonVoteMessages,
			messagebus.AllPlayerVotedOnTeam{
				Event:        messagebus.Event{Party: messagebus.Party{Code: code}},
				Approved:     false,
				VoteFailures: updatedGame.VoteFailures(),
				PlayerVotes:  resultingVote,
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			messagebus.LeaderStartedToSelectMembers{
				Event:  messagebus.Event{Party: messagebus.Party{Code: code}},
				Leader: updatedGame.Leader(),
			},
		)
	} else if updatedGame.State() == gamerules.ConductingMission {
		commonVoteMessages = append(commonVoteMessages,
			messagebus.AllPlayerVotedOnTeam{
				Event:       messagebus.Event{Party: messagebus.Party{Code: code}},
				Approved:    true,
				PlayerVotes: resultingVote,
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			messagebus.MissionStarted{
				Event: messagebus.Event{Party: messagebus.Party{Code: code}},
			},
		)
	} else if updatedGame.State() == gamerules.GameOver {
		commonVoteMessages = append(commonVoteMessages,
			messagebus.AllPlayerVotedOnTeam{
				Event:        messagebus.Event{Party: messagebus.Party{Code: code}},
				Approved:     false,
				VoteFailures: updatedGame.VoteFailures(),
				PlayerVotes:  resultingVote,
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			messagebus.GameEnded{
				Event:  messagebus.Event{Party: messagebus.Party{Code: code}},
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
			Event:   messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
			Player:  succeedMissionCommand.Player,
			Success: true,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonMissionOutgoingMessages(updatedGame, message.GetPartyCode(), outcomes)...)

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
			Event:   messagebus.Event{Party: messagebus.Party{Code: message.GetPartyCode()}},
			Player:  failMissionCommand.Player,
			Success: false,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonMissionOutgoingMessages(updatedGame, message.GetPartyCode(), outcomes)...)

	return
}

func commonMissionOutgoingMessages(updatedGame gamerules.Game, code string, outcomes map[string]bool) []messagebus.Message {
	commonMissionMessages := []messagebus.Message{}

	if updatedGame.State() == gamerules.SelectingTeam {
		lastMissionSuccess := updatedGame.GetMissionResults()[updatedGame.CurrentMission()-1]
		talliedOutcomes := tallyOutcomes(outcomes)
		commonMissionMessages = append(commonMissionMessages,
			messagebus.MissionCompleted{
				Event:    messagebus.Event{Party: messagebus.Party{Code: code}},
				Success:  lastMissionSuccess,
				Outcomes: talliedOutcomes,
			},
		)
		commonMissionMessages = append(commonMissionMessages,
			messagebus.LeaderStartedToSelectMembers{
				Event:  messagebus.Event{Party: messagebus.Party{Code: code}},
				Leader: updatedGame.Leader(),
			},
		)
	} else if updatedGame.State() == gamerules.GameOver {
		lastMissionSuccess := updatedGame.GetMissionResults()[updatedGame.CurrentMission()]
		talliedOutcomes := tallyOutcomes(outcomes)
		commonMissionMessages = append(commonMissionMessages,
			messagebus.MissionCompleted{
				Event:    messagebus.Event{Party: messagebus.Party{Code: code}},
				Success:  lastMissionSuccess,
				Outcomes: talliedOutcomes,
			},
		)

		commonMissionMessages = append(commonMissionMessages,
			messagebus.GameEnded{
				Event:  messagebus.Event{Party: messagebus.Party{Code: code}},
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
