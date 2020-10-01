package playeractions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

type mockSessionGetter struct {
	receivedSession string
	getError        error
}

func (m *mockSessionGetter) Get(session string) (code string, name string, err error) {
	m.receivedSession = session
	return "testCode", "testName", m.getError
}

type mockActionBroker struct {
	receivedCode string
}

func (m *mockActionBroker) StartGame(code string) {
	m.receivedCode = code
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
	req, _ := http.NewRequest("GET", "/actions/start-game", nil)
	req.AddCookie(&http.Cookie{Name: "other", Value: "value"})
	_, _, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(401))
}

func Test_CheckSessionMiddleware_Return403IfGettingSessionReturnsError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/actions/start-game", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter := &mockSessionGetter{}
	sessionGetter.getError = fmt.Errorf("get error")
	_, _, w := makeCall(req, sessionGetter)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(403))
	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
}

func Test_StartGame(t *testing.T) {
	req, _ := http.NewRequest("GET", "/actions/start-game", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testSession"})
	sessionGetter, actionBroker, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(actionBroker.receivedCode).To(Equal("testCode"))
}
