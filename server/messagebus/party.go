package messagebus

type Party struct {
	Code string
}

func (p Party) GetPartyCode() string {
	return p.Code
}
