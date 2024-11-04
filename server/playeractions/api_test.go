package playeractions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

type mockSessionGetter struct {
	receivedSession string
	getError        error
}

func (m *mockSessionGetter) Get(session string) (name string, err error) {
	m.receivedSession = session
	return "testName", m.getError
}

type mockActionBroker struct {
	gameStarted              bool
	receivedLeader           string
	receivedSelectedMember   string
	receivedDeselectedMember string
	teamConfirmed            bool
	receivedPlayerApprove    string
	receivedPlayerReject     string
	receivedPlayerSucceed    string
	receivedPlayerFail       string
}

func (m *mockActionBroker) StartGame() {
	m.gameStarted = true
}

func (m *mockActionBroker) LeaderSelectsMember(leader string, member string) {
	m.receivedLeader = leader
	m.receivedSelectedMember = member
}

func (m *mockActionBroker) LeaderDeselectsMember(leader string, member string) {
	m.receivedLeader = leader
	m.receivedDeselectedMember = member
}

func (m *mockActionBroker) LeaderConfirmsTeam(leader string) {
	m.receivedLeader = leader
	m.teamConfirmed = true
}

func (m *mockActionBroker) ApproveTeam(player string) {
	m.receivedPlayerApprove = player
}

func (m *mockActionBroker) RejectTeam(player string) {
	m.receivedPlayerReject = player
}

func (m *mockActionBroker) SucceedMission(player string) {
	m.receivedPlayerSucceed = player
}

func (m *mockActionBroker) FailMission(player string) {
	m.receivedPlayerFail = player
}

func jsonReader(obj interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(obj)
	return bytes.NewReader(jsonBytes)
}

func makeCall(req *http.Request, sessionGetter *mockSessionGetter) (*mockSessionGetter, *mockActionBroker, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	ginEngine := gin.New()

	if sessionGetter == nil {
		sessionGetter = &mockSessionGetter{}
	}

	actionBroker := &mockActionBroker{}
	Register(ginEngine, sessionGetter, actionBroker)

	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	return sessionGetter, actionBroker, w
}

func Test_CheckSessionMiddleware_Return401IfNoSessionCookie(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/start-game", nil)
	req.AddCookie(&http.Cookie{Name: "other", Value: "value"})
	_, _, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(401))
}

func Test_CheckSessionMiddleware_Return403IfGettingSessionReturnsError(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/start-game", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter := &mockSessionGetter{}
	sessionGetter.getError = fmt.Errorf("get error")
	_, _, w := makeCall(req, sessionGetter)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(403))
	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
}

func Test_StartGame(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/start-game", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.gameStarted).To(BeTrue())
}

func Test_LeaderSelectsMember(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/leader-selects-member", jsonReader(leaderSelectionRequest{Member: "aMember"}))
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedLeader).To(Equal("testName"))
	g.Expect(actionBroker.receivedSelectedMember).To(Equal("aMember"))
}

func Test_LeaderSelectsMember_Returns400IfMemberIsMissing(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/leader-selects-member", jsonReader(leaderSelectionRequest{Member: ""}))
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedLeader).To(Equal(""))
	g.Expect(actionBroker.receivedSelectedMember).To(Equal(""))
}

func Test_LeaderSelectsMember_Returns400IfBodyMalformed(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/leader-selects-member", strings.NewReader("garbage"))
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedLeader).To(Equal(""))
	g.Expect(actionBroker.receivedSelectedMember).To(Equal(""))
}

func Test_LeaderDeselectsMember(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/leader-deselects-member", jsonReader(leaderSelectionRequest{Member: "aMember"}))
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedLeader).To(Equal("testName"))
	g.Expect(actionBroker.receivedDeselectedMember).To(Equal("aMember"))
}

func Test_LeaderDeselectsMember_Returns400IfMemberIsMissing(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/leader-deselects-member", jsonReader(leaderSelectionRequest{Member: ""}))
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedLeader).To(Equal(""))
	g.Expect(actionBroker.receivedDeselectedMember).To(Equal(""))
}

func Test_LeaderDeselectsMember_Returns400IfBodyMalformed(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/leader-deselects-member", strings.NewReader("garbage"))
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedLeader).To(Equal(""))
	g.Expect(actionBroker.receivedDeselectedMember).To(Equal(""))
}

func Test_LeaderConfirmsTeam(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/leader-confirms-team", strings.NewReader("garbage"))
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedLeader).To(Equal("testName"))
	g.Expect(actionBroker.teamConfirmed).To(BeTrue())
}

func Test_ApproveTeam(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/approve-team", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedPlayerApprove).To(Equal("testName"))
}

func Test_RejectTeam(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/reject-team", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedPlayerReject).To(Equal("testName"))
}

func Test_SucceedMission(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/succeed-mission", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedPlayerSucceed).To(Equal("testName"))
}

func Test_FailMission(t *testing.T) {
	req, _ := http.NewRequest("POST", "/actions/fail-mission", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedPlayerFail).To(Equal("testName"))
}
