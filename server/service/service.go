package service

import (
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type handler func(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message)

type codeGenerator interface {
	generateCode() string
}

type messageDispatcher interface {
	dispatchMessage(m messagebus.Message)
}

type service struct {
	codeGenerator       codeGenerator
	messageDispatcher   messageDispatcher
	allegianceGenerator gamerules.AllegianceGenerator
	gamesByPartyCode    map[string]gamerules.Game
}

func newService(codeGenerator codeGenerator, messageDispatcher messageDispatcher, allegianceGenerator gamerules.AllegianceGenerator) service {
	return service{
		codeGenerator:       codeGenerator,
		messageDispatcher:   messageDispatcher,
		allegianceGenerator: allegianceGenerator,
		gamesByPartyCode:    make(map[string]gamerules.Game),
	}
}

func (s service) createParty() string {
	newCode := s.codeGenerator.generateCode()
	s.gamesByPartyCode[newCode] = gamerules.NewGame()
	return newCode
}

func (s service) handleMessage(m messagebus.Message) {
	var handler handler
	switch m.(type) {
	case joinParty:
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
		s.messageDispatcher.dispatchMessage(messageToDispatch)
	}
}

func (s service) getGameForPartyCode(code string) gamerules.Game {
	return s.gamesByPartyCode[code]
}

func (s service) handleJoinPartyCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	joinPartyCommand := message.(joinParty)
	updatedGame, err := currentGame.AddPlayer(joinPartyCommand.user)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			playerJoined{
				party: party{code: message.GetPartyCode()},
				user:  joinPartyCommand.user,
			},
		)
	}
	return
}

func (s service) handleStartGameCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	updatedGame, err := currentGame.Start(s.allegianceGenerator)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			leaderStartedToSelectMembers{
				party:  party{code: message.GetPartyCode()},
				leader: updatedGame.Leader(),
			},
		)
	}
	return
}

func (s service) handleLeaderSelectsMember(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderSelectsMemberCommand := message.(leaderSelectsMember)

	if leaderSelectsMemberCommand.leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderSelectsMember(leaderSelectsMemberCommand.memberToSelect)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			leaderSelectedMember{
				party:          party{code: message.GetPartyCode()},
				selectedMember: leaderSelectsMemberCommand.memberToSelect,
			},
		)
	}
	return
}

func (s service) handleLeaderDeselectsMember(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderDeselectsMemberCommand := message.(leaderDeselectsMember)

	if leaderDeselectsMemberCommand.leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderDeselectsMember(leaderDeselectsMemberCommand.memberToDeselect)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			leaderDeselectedMember{
				party:            party{code: message.GetPartyCode()},
				deselectedMember: leaderDeselectsMemberCommand.memberToDeselect,
			},
		)
	}
	return
}

func (s service) handleLeaderConfirmsTeamSelection(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	leaderConfirmsTeamSelectionCommand := message.(leaderConfirmsTeamSelection)

	if leaderConfirmsTeamSelectionCommand.leader != currentGame.Leader() {
		updatedGame = currentGame
		return
	}

	updatedGame, err := currentGame.LeaderConfirmsTeamSelection()

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			leaderConfirmedSelection{
				party: party{code: message.GetPartyCode()},
			},
		)
	}
	return
}

func (s service) handleApproveTeam(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	approveTeamCommand := message.(approveTeam)
	updatedGame, err := currentGame.ApproveTeamBy(approveTeamCommand.player)

	if err != nil {
		updatedGame = currentGame
		return
	}

	messagesToDispatch = append(messagesToDispatch,
		playerVotedOnTeam{
			party:    party{code: message.GetPartyCode()},
			player:   approveTeamCommand.player,
			approved: true,
		},
	)

	if updatedGame.State() == gamerules.SelectingTeam {
		messagesToDispatch = append(messagesToDispatch,
			allPlayerVotedOnTeam{
				party:        party{code: message.GetPartyCode()},
				approved:     false,
				voteFailures: updatedGame.VoteFailures(),
			},
		)
	} else if updatedGame.State() == gamerules.ConductingMission {
		messagesToDispatch = append(messagesToDispatch,
			allPlayerVotedOnTeam{
				party:    party{code: message.GetPartyCode()},
				approved: true,
			},
		)
	}

	return
}

func (s service) handleRejectTeam(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	rejectTeamCommand := message.(rejectTeam)
	updatedGame, err := currentGame.RejectTeamBy(rejectTeamCommand.player)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			playerVotedOnTeam{
				party:    party{code: message.GetPartyCode()},
				player:   rejectTeamCommand.player,
				approved: false,
			},
		)
	}
	return
}

func (s service) handleSucceedMission(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	succeedMissionCommand := message.(succeedMission)
	updatedGame, err := currentGame.SucceedMissionBy(succeedMissionCommand.player)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			playerWorkedOnMission{
				party:   party{code: message.GetPartyCode()},
				player:  succeedMissionCommand.player,
				success: true,
			},
		)
	}
	return
}

func (s service) handleFailMission(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messagesToDispatch []messagebus.Message) {
	failMissionCommand := message.(failMission)
	updatedGame, err := currentGame.FailMissionBy(failMissionCommand.player)

	if err == nil {
		messagesToDispatch = append(messagesToDispatch,
			playerWorkedOnMission{
				party:   party{code: message.GetPartyCode()},
				player:  failMissionCommand.player,
				success: false,
			},
		)
	}
	return
}
