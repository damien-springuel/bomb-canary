package party

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

type mockPartyBroker struct {
	givenName string
}

func (m *mockPartyBroker) CreateParty() string {
	return "testCode"
}

func (m *mockPartyBroker) JoinParty(name string) {
	m.givenName = name
}

type mockSession struct {
	givenName string
}

func (m *mockSession) Create(name string) string {
	m.givenName = name
	return "testSessionId"
}

func jsonReader(obj interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(obj)
	return bytes.NewReader(jsonBytes)
}

func makeCall(req *http.Request, partyBroker *mockPartyBroker) (*mockPartyBroker, *mockSession, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	ginEngine := gin.New()

	if partyBroker == nil {
		partyBroker = &mockPartyBroker{}
	}
	sessions := &mockSession{}
	Register(ginEngine, partyBroker, sessions)

	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	return partyBroker, sessions, w
}

func Test_JoinParty(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/join", jsonReader(joinPartyRequest{Name: "testName"}))
	_, sessions, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	actualCookie := w.Result().Cookies()[0]
	g.Expect(actualCookie.Name).To(Equal("session"))
	g.Expect(actualCookie.Value).To(Equal("testSessionId"))
	g.Expect(actualCookie.MaxAge).To(Equal(int((time.Hour * 5).Seconds())))

	g.Expect(*sessions).To(Equal(mockSession{givenName: "testName"}))
}

func Test_JoinParty_Should400IfNameAbsent(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/join", jsonReader(joinPartyRequest{Name: ""}))
	partyBroker, sessions, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyBroker).To(Equal(mockPartyBroker{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_JoinParty_Should400IfMalformedBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/join", strings.NewReader("garbage"))
	partyBroker, sessions, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyBroker).To(Equal(mockPartyBroker{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}
