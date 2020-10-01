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

type GameHub struct {
	codeGenerator       codeGenerator
	messageDispatcher   messageDispatcher
	allegianceGenerator gamerules.AllegianceGenerator
	gamesByPartyCode    map[string]gamerules.Game
}

func New(codeGenerator codeGenerator, messageDispatcher messageDispatcher, allegianceGenerator gamerules.AllegianceGenerator) GameHub {
	return GameHub{
		codeGenerator:       codeGenerator,
		messageDispatcher:   messageDispatcher,
		allegianceGenerator: allegianceGenerator,
		gamesByPartyCode:    make(map[string]gamerules.Game),
	}
}

func (s GameHub) CreateParty() string {
	newCode := s.codeGenerator.GenerateCode()
	s.gamesByPartyCode[newCode] = gamerules.NewGame()
	return newCode
}

func (s GameHub) HandleMessage(m messagebus.Message) {
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
	updatedGame, messagesToDispatch := handler(s.gamesByPartyCode[m.GetPartyCode()], m)

	s.gamesByPartyCode[m.GetPartyCode()] = updatedGame
	for _, messageToDispatch := range messagesToDispatch {
		s.messageDispatcher.Dispatch(messageToDispatch)
	}
}

func (s GameHub) getGameForPartyCode(code string) gamerules.Game {
	return s.gamesByPartyCode[code]
}

func (s GameHub) handleJoinPartyCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
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

func (s GameHub) handleStartGameCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
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

func (s GameHub) handleLeaderSelectsMember(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
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

func (s GameHub) handleLeaderDeselectsMember(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
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

func (s GameHub) handleLeaderConfirmsTeamSelection(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
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

func (s GameHub) handleApproveTeam(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
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

func (s GameHub) handleRejectTeam(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
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

func commonVoteOutgoingMessages(updatedGame gamerules.Game, code string) []messagebus.Message {
	commonVoteMessages := []messagebus.Message{}
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

func (s GameHub) handleSucceedMission(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
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

func (s GameHub) handleFailMission(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	failMissionCommand := message.(FailMission)
	updatedGame, err := currentGame.FailMissionBy(failMissionCommand.player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		PlayerWorkedOnMission{
			Party:   Party{Code: message.GetPartyCode()},
			Player:  failMissionCommand.player,
			Success: false,
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
			MissionCompleted{
				Party:   Party{Code: code},
				success: lastGameSuccess,
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
				success: lastGameSuccess,
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