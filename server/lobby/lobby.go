package lobby

import (
	"context"

	"github.com/damien-springuel/bomb-canary/server/generated"
)

type LobbyServer struct {
	generated.UnimplementedLobbyServer
}

func (l LobbyServer) CreateParty(context.Context, *generated.CreatePartyRequest) (*generated.CreatePartyResponse, error) {
	return &generated.CreatePartyResponse{
		PartyCode: "MyCoolCode",
	}, nil
}
