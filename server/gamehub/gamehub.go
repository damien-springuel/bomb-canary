package gamehub

import (
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	. "github.com/damien-springuel/bomb-canary/server/messagebus"
)

type handler func(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message)

type codeGenerator interface {
	GenerateCode() string
}

type messageDispatcher interface {
	Dispatch(m Message)
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
	return newCode
}

func (s gameHub) DoesPartyExist(code string) bool {
	_, exists := s.games.get(code)
	return exists
}

func (s gameHub) Consume(m Message) {
	var handler handler
	switch m.(type) {
	case JoinParty:
		handler = s.handleJoinPartyCommand
	case StartGame:
		handler = s.handleStartGameCommand
	case LeaderSelectsMember:
		handler = s.handleLeaderSelectsMember
	case LeaderDeselectsMember:
		handler = s.handleLeaderDeselectsMember
	case LeaderConfirmsTeamSelection:
		handler = s.handleLeaderConfirmsTeamSelection
	case ApproveTeam:
		handler = s.handleApproveTeam
	case RejectTeam:
		handler = s.handleRejectTeam
	case SucceedMission:
		handler = s.handleSucceedMission
	case FailMission:
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

func (s gameHub) handleJoinPartyCommand(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	joinPartyCommand := message.(JoinParty)
	updatedGame, err := currentGame.AddPlayer(joinPartyCommand.User)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			PlayerJoined{
				Party: Party{Code: message.GetPartyCode()},
				User:  joinPartyCommand.User,
			},
		)
	}
	return
}

func (s gameHub) handleStartGameCommand(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	updatedGame, err := currentGame.Start(s.allegianceGenerator)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			LeaderStartedToSelectMembers{
				Party:  Party{Code: message.GetPartyCode()},
				Leader: updatedGame.Leader(),
			},
		)
	}
	return
}

func (s gameHub) handleLeaderSelectsMember(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	leaderSelectsMemberCommand := message.(LeaderSelectsMember)

	if leaderSelectsMemberCommand.Leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderSelectsMember(leaderSelectsMemberCommand.MemberToSelect)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			LeaderSelectedMember{
				Party:          Party{Code: message.GetPartyCode()},
				SelectedMember: leaderSelectsMemberCommand.MemberToSelect,
			},
		)
	}
	return
}

func (s gameHub) handleLeaderDeselectsMember(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	leaderDeselectsMemberCommand := message.(LeaderDeselectsMember)

	if leaderDeselectsMemberCommand.Leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderDeselectsMember(leaderDeselectsMemberCommand.MemberToDeselect)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			LeaderDeselectedMember{
				Party:            Party{Code: message.GetPartyCode()},
				DeselectedMember: leaderDeselectsMemberCommand.MemberToDeselect,
			},
		)
	}
	return
}

func (s gameHub) handleLeaderConfirmsTeamSelection(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	leaderConfirmsTeamSelectionCommand := message.(LeaderConfirmsTeamSelection)

	if leaderConfirmsTeamSelectionCommand.Leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderConfirmsTeamSelection()

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			LeaderConfirmedSelection{
				Party: Party{Code: message.GetPartyCode()},
			},
		)
	}
	return
}

func (s gameHub) handleApproveTeam(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	approveTeamCommand := message.(ApproveTeam)
	updatedGame, err := currentGame.ApproveTeamBy(approveTeamCommand.Player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		PlayerVotedOnTeam{
			Party:    Party{Code: message.GetPartyCode()},
			Player:   approveTeamCommand.Player,
			Approved: true,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonVoteOutgoingMessages(updatedGame, message.GetPartyCode())...)

	return
}

func (s gameHub) handleRejectTeam(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	rejectTeamCommand := message.(RejectTeam)
	updatedGame, err := currentGame.RejectTeamBy(rejectTeamCommand.Player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		PlayerVotedOnTeam{
			Party:    Party{Code: message.GetPartyCode()},
			Player:   rejectTeamCommand.Player,
			Approved: false,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonVoteOutgoingMessages(updatedGame, message.GetPartyCode())...)

	return
}

func commonVoteOutgoingMessages(updatedGame gamerules.Game, code string) []Message {
	commonVoteMessages := []Message{}
	if updatedGame.State() == gamerules.SelectingTeam {
		commonVoteMessages = append(commonVoteMessages,
			AllPlayerVotedOnTeam{
				Party:        Party{Code: code},
				Approved:     false,
				VoteFailures: updatedGame.VoteFailures(),
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			LeaderStartedToSelectMembers{
				Party:  Party{Code: code},
				Leader: updatedGame.Leader(),
			},
		)
	} else if updatedGame.State() == gamerules.ConductingMission {
		commonVoteMessages = append(commonVoteMessages,
			AllPlayerVotedOnTeam{
				Party:    Party{Code: code},
				Approved: true,
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			MissionStarted{
				Party: Party{Code: code},
			},
		)
	} else if updatedGame.State() == gamerules.GameOver {
		commonVoteMessages = append(commonVoteMessages,
			AllPlayerVotedOnTeam{
				Party:        Party{Code: code},
				Approved:     false,
				VoteFailures: updatedGame.VoteFailures(),
			},
		)
		commonVoteMessages = append(commonVoteMessages,
			GameEnded{
				Party:  Party{Code: code},
				Winner: Spy,
			},
		)
	}

	return commonVoteMessages
}

func (s gameHub) handleSucceedMission(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	succeedMissionCommand := message.(SucceedMission)
	updatedGame, err := currentGame.SucceedMissionBy(succeedMissionCommand.Player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		PlayerWorkedOnMission{
			Party:   Party{Code: message.GetPartyCode()},
			Player:  succeedMissionCommand.Player,
			Success: true,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonMissionOutgoingMessages(updatedGame, message.GetPartyCode())...)

	return
}

func (s gameHub) handleFailMission(currentGame gamerules.Game, message Message) (updatedGame gamerules.Game, messagesToDispatch []Message) {
	failMissionCommand := message.(FailMission)
	updatedGame, err := currentGame.FailMissionBy(failMissionCommand.Player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		PlayerWorkedOnMission{
			Party:   Party{Code: message.GetPartyCode()},
			Player:  failMissionCommand.Player,
			Success: false,
		},
	)

	messagesToDispatch = append(messagesToDispatch, commonMissionOutgoingMessages(updatedGame, message.GetPartyCode())...)

	return
}

func commonMissionOutgoingMessages(updatedGame gamerules.Game, code string) []Message {
	commonMissionMessages := []Message{}

	if updatedGame.State() == gamerules.SelectingTeam {
		lastGameSuccess := updatedGame.GetMissionResults()[updatedGame.CurrentMission()-1]
		commonMissionMessages = append(commonMissionMessages,
			MissionCompleted{
				Party:   Party{Code: code},
				Success: lastGameSuccess,
			},
		)
		commonMissionMessages = append(commonMissionMessages,
			LeaderStartedToSelectMembers{
				Party:  Party{Code: code},
				Leader: updatedGame.Leader(),
			},
		)
	} else if updatedGame.State() == gamerules.GameOver {
		lastGameSuccess := updatedGame.GetMissionResults()[updatedGame.CurrentMission()-1]
		commonMissionMessages = append(commonMissionMessages,
			MissionCompleted{
				Party:   Party{Code: code},
				Success: lastGameSuccess,
			},
		)

		winner := updatedGame.Winner()
		var messageWinner Allegiance
		if winner == gamerules.Resistance {
			messageWinner = Resistance
		} else if winner == gamerules.Spy {
			messageWinner = Spy
		}
		commonMissionMessages = append(commonMissionMessages,
			GameEnded{
				Party:  Party{Code: code},
				Winner: messageWinner,
			},
		)
	}

	return commonMissionMessages
}
