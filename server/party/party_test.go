package party

import (
	"testing"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type mockDispatcher struct {
	receivedMessage messagebus.Message
}

func (m *mockDispatcher) Dispatch(message messagebus.Message) {
	m.receivedMessage = message
}

type mockPartyCreator struct{}

func (m *mockPartyCreator) CreateParty() string {
	return "testCode"
}

type mockPartyChecker struct {
	response bool
}

func (m *mockPartyChecker) DoesPartyExist(code string) bool {
	return m.response
}

func Test_ServiceCreateParty(t *testing.T) {
	dispatcher := &mockDispatcher{}
	creator := &mockPartyCreator{}
	checker := &mockPartyChecker{}
	service := NewPartyService(creator, checker, dispatcher)

	g := NewWithT(t)
	g.Expect(service.CreateParty()).To(Equal("testCode"))
}

func Test_ServiceJoinParty_PartyExists(t *testing.T) {
	dispatcher := &mockDispatcher{}
	creator := &mockPartyCreator{}
	checker := &mockPartyChecker{response: true}
	service := NewPartyService(creator, checker, dispatcher)

	g := NewWithT(t)
	g.Expect(service.JoinParty("code", "name")).To(BeNil())
	g.Expect(dispatcher.receivedMessage).To(Equal(messagebus.JoinParty{Party: messagebus.Party{Code: "code"}, User: "name"}))
}

func Test_ServiceJoinParty_PartyDoesntExist(t *testing.T) {
	dispatcher := &mockDispatcher{}
	creator := &mockPartyCreator{}
	checker := &mockPartyChecker{response: false}
	service := NewPartyService(creator, checker, dispatcher)

	g := NewWithT(t)
	g.Expect(service.JoinParty("code", "name")).To(MatchError(errPartyDoesntExists))
	g.Expect(dispatcher.receivedMessage).To(BeNil())
}
