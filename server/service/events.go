package service

type playerJoined struct {
	party
	user string
}

type gameStarted struct {
	party
}
