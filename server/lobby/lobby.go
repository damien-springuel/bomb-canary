package lobby

type createPartyRequest struct{}
type createPartyResponse struct {
	Code string `json:"code"`
}

type lobbyServer struct{}

func (l lobbyServer) createParty(r createPartyRequest) (createPartyResponse, error) {
	return createPartyResponse{Code: "My cool code"}, nil
}
