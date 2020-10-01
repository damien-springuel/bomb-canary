package service

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

type Service struct {
	codeGenerator       codeGenerator
	messageDispatcher   messageDispatcher
	allegianceGenerator gamerules.AllegianceGenerator
	gamesByPartyCode    map[string]gamerules.Game
}

func New(codeGenerator codeGenerator, messageDispatcher messageDispatcher, allegianceGenerator gamerules.AllegianceGenerator) Service {
	return Service{
		codeGenerator:       codeGenerator,
		messageDispatcher:   messageDispatcher,
		allegianceGenerator: allegianceGenerator,
		gamesByPartyCode:    make(map[string]gamerules.Game),
	}
}

func (s Service) CreateParty() string {
	newCode := s.codeGenerator.GenerateCode()
	s.gamesByPartyCode[newCode] = gamerules.NewGame()
	return newCode
}

func (s Service) HandleMessage(m messagebus.Message) {
	var handler handler
	switch m.(type) {
	case JoinParty:
		handler = s.handleJoinPartyCommand
	case startGame:
		handler = s.handleStartGameCommand
	case leaderSelectsMember:
		handler = s.handleLeaderSelectsMember
	case leaderDeselectsMember:
		handler = s.handleLeaderDeselectsMember
	case leaderConfirmsTeamSelection:
		handler = s.handleLeaderConfirmsTeamSelection
	case approveTeam:
		handler = s.handleApproveTeam
	case rejectTeam:
		handler = s.handleRejectTeam
	case succeedMission:
		handler = s.handleSucceedMission
	case failMission:
		handler = s.handleFailMission
	default:
		return
	}
	updatedGame, messagesToDispatch := handler(s.gamesByPartyCode[m.GetPartyCode()], m)

	s.gamesByPartyCode[m.GetPartyCode()] = updatedGame
	for _, messageToDispatch := range messagesToDispatch {
		s.messageDispatcher.Dispatch(messageToDispatch)
	}
}

func (s Service) getGameForPartyCode(code string) gamerules.Game {
	return s.gamesByPartyCode[code]
}

func (s Service) handleJoinPartyCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	joinPartyCommand := message.(JoinParty)
	updatedGame, err := currentGame.AddPlayer(joinPartyCommand.User)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			playerJoined{
				Party: Party{Code: message.GetPartyCode()},
				user:  joinPartyCommand.User,
			},
		)
	}
	return
}

func (s Service) handleStartGameCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	updatedGame, err := currentGame.Start(s.allegianceGenerator)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			leaderStartedToSelectMembers{
				Party:  Party{Code: message.GetPartyCode()},
				leader: updatedGame.Leader(),
			},
		)
	}
	return
}

func (s Service) handleLeaderSelectsMember(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderSelectsMemberCommand := message.(leaderSelectsMember)

	if leaderSelectsMemberCommand.leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderSelectsMember(leaderSelectsMemberCommand.memberToSelect)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			leaderSelectedMember{
				Party:          Party{Code: message.GetPartyCode()},
				selectedMember: leaderSelectsMemberCommand.memberToSelect,
			},
		)
	}
	return
}

func (s Service) handleLeaderDeselectsMember(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderDeselectsMemberCommand := message.(leaderDeselectsMember)

	if leaderDeselectsMemberCommand.leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderDeselectsMember(leaderDeselectsMemberCommand.memberToDeselect)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			leaderDeselectedMember{
				Party:            Party{Code: message.GetPartyCode()},
				deselectedMember: leaderDeselectsMemberCommand.memberToDeselect,
			},
		)
	}
	return
}

func (s Service) handleLeaderConfirmsTeamSelection(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderConfirmsTeamSelectionCommand := message.(leaderConfirmsTeamSelection)

	if leaderConfirmsTeamSelectionCommand.leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderConfirmsTeamSelection()

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			leaderConfirmedSelection{
				Party: Party{Code: message.GetPartyCode()},
			},
		)
	}
	return
}

func (s Service) handleApproveTeam(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	approveTeamCommand := message.(approveTeam)
	updatedGame, err := currentGame.ApproveTeamBy(approveTeamCommand.player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		playerVotedOnTeam{
			Party:    Party{Code: message.GetPartyCode()},
			player:   approveTeamCommand.player,
			approved: true,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonVoteOutgoingMessages(updatedGame, message.GetPartyCode())...)

	return
}

func (s Service) handleRejectTeam(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	rejectTeamCommand := message.(rejectTeam)
	updatedGame, err := currentGame.RejectTeamBy(rejectTeamCommand.player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		playerVotedOnTeam{
			Party:    Party{Code: message.GetPartyCode()},
			player:   rejectTeamCommand.player,
			approved: false,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonVoteOutgoingMessages(updatedGame, message.GetPartyCode())...)

	return
}

func commonVoteOutgoingMessages(updatedGame gamerules.Game, code string) []messagebus.Message {
	commonVoteMessages := []messagebus.Message{}
	if updatedGame.State() == gamerules.SelectingTeam {
		commonVoteMessages = append(commonVoteMessages,
			allPlayerVotedOnTeam{
				Party:        Party{Code: code},
				approved:     false,
				voteFailures: updatedGame.VoteFailures(),
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			leaderStartedToSelectMembers{
				Party:  Party{Code: code},
				leader: updatedGame.Leader(),
			},
		)
	} else if updatedGame.State() == gamerules.ConductingMission {
		commonVoteMessages = append(commonVoteMessages,
			allPlayerVotedOnTeam{
				Party:    Party{Code: code},
				approved: true,
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			missionStarted{
				Party: Party{Code: code},
			},
		)
	} else if updatedGame.State() == gamerules.GameOver {
		commonVoteMessages = append(commonVoteMessages,
			allPlayerVotedOnTeam{
				Party:        Party{Code: code},
				approved:     false,
				voteFailures: updatedGame.VoteFailures(),
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			gameEnded{
				Party:  Party{Code: code},
				winner: Spy,
			},
		)
	}

	return commonVoteMessages
}

func (s Service) handleSucceedMission(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	succeedMissionCommand := message.(succeedMission)
	updatedGame, err := currentGame.SucceedMissionBy(succeedMissionCommand.player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		playerWorkedOnMission{
			Party:   Party{Code: message.GetPartyCode()},
			player:  succeedMissionCommand.player,
			success: true,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonMissionOutgoingMessages(updatedGame, message.GetPartyCode())...)

	return
}

func (s Service) handleFailMission(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	failMissionCommand := message.(failMission)
	updatedGame, err := currentGame.FailMissionBy(failMissionCommand.player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		playerWorkedOnMission{
			Party:   Party{Code: message.GetPartyCode()},
			player:  failMissionCommand.player,
			success: false,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonMissionOutgoingMessages(updatedGame, message.GetPartyCode())...)

	return
}

func commonMissionOutgoingMessages(updatedGame gamerules.Game, code string) []messagebus.Message {
	commonMissionMessages := []messagebus.Message{}

	if updatedGame.State() == gamerules.SelectingTeam {
		lastGameSuccess := updatedGame.GetMissionResults()[updatedGame.CurrentMission()-1]
		commonMissionMessages = append(commonMissionMessages,
			missionCompleted{
				Party:   Party{Code: code},
				success: lastGameSuccess,
			},
		)
		commonMissionMessages = append(commonMissionMessages,
			leaderStartedToSelectMembers{
				Party:  Party{Code: code},
				leader: updatedGame.Leader(),
			},
		)
	} else if updatedGame.State() == gamerules.GameOver {
		lastGameSuccess := updatedGame.GetMissionResults()[updatedGame.CurrentMission()-1]
		commonMissionMessages = append(commonMissionMessages,
			missionCompleted{
				Party:   Party{Code: code},
				success: lastGameSuccess,
			},
		)

		winner := updatedGame.Winner()
		var messageWinner allegiance
		if winner == gamerules.Resistance {
			messageWinner = Resistance
		} else if winner == gamerules.Spy {
			messageWinner = Spy
		}
		commonMissionMessages = append(commonMissionMessages,
			gameEnded{
				Party:  Party{Code: code},
				winner: messageWinner,
			},
		)
	}

	return commonMissionMessages
}
