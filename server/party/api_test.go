package party

import (
	"bytes"
	"encoding/json"
	"errors"
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
	givenCode         string
	givenName         string
	joinPartyResponse error
}

func (m *mockPartyBroker) CreateParty() string {
	return "testCode"
}

func (m *mockPartyBroker) JoinParty(code string, name string) error {
	m.givenCode = code
	m.givenName = name
	return m.joinPartyResponse
}

type mockSession struct {
	givenCode string
	givenName string
}

func (m *mockSession) Create(code string, name string) string {
	m.givenCode = code
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

func Test_CreateParty(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/create", jsonReader(createPartyRequest{Name: "testName"}))
	partyBroker, sessions, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))

	actualResponse := createPartyResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &actualResponse)
	g.Expect(actualResponse).To(Equal(createPartyResponse{Code: "testCode"}))

	actualCookie := w.Result().Cookies()[0]
	g.Expect(actualCookie.Name).To(Equal("session"))
	g.Expect(actualCookie.Value).To(Equal("testSessionId"))
	g.Expect(actualCookie.MaxAge).To(Equal(int((time.Hour * 3).Seconds())))

	g.Expect(*partyBroker).To(Equal(mockPartyBroker{givenCode: "testCode", givenName: "testName"}))
	g.Expect(*sessions).To(Equal(mockSession{givenCode: "testCode", givenName: "testName"}))
}

func Test_CreateParty_ShouldReturn400IfNameIsAbsent(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/create", jsonReader(createPartyRequest{Name: ""}))
	partyBroker, sessions, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))
	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyBroker).To(Equal(mockPartyBroker{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_CreateParty_ShouldReturn400IfMalformedBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/create", strings.NewReader("garbage"))
	partyBroker, sessions, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyBroker).To(Equal(mockPartyBroker{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_JoinParty(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/join", jsonReader(joinPartyRequest{Code: "joinCode", Name: "testName"}))
	_, sessions, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Body.String()).To(Equal("{}"))

	actualCookie := w.Result().Cookies()[0]
	g.Expect(actualCookie.Name).To(Equal("session"))
	g.Expect(actualCookie.Value).To(Equal("testSessionId"))
	g.Expect(actualCookie.MaxAge).To(Equal(int((time.Hour * 3).Seconds())))

	g.Expect(*sessions).To(Equal(mockSession{givenCode: "joinCode", givenName: "testName"}))
}

func Test_JoinParty_Should400IfCodeAbsent(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/join", jsonReader(joinPartyRequest{Code: "", Name: "testName"}))
	partyBroker, sessions, w := makeCall(req, nil)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyBroker).To(Equal(mockPartyBroker{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_JoinParty_Should400IfNameAbsent(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/join", jsonReader(joinPartyRequest{Code: "joinCode", Name: ""}))
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

func Test_JoinParty_Should404IfCantJoinBecausePartyDoesntExist(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/join", jsonReader(joinPartyRequest{Code: "code", Name: "name"}))
	partyBroker := &mockPartyBroker{}
	partyBroker.joinPartyResponse = errPartyDoesntExists
	_, sessions, w := makeCall(req, partyBroker)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(404))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyBroker).To(Equal(mockPartyBroker{givenCode: "code", givenName: "name", joinPartyResponse: errPartyDoesntExists}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_JoinParty_Should500IfJoinReturnsOtherError(t *testing.T) {
	req, _ := http.NewRequest("POST", "/party/join", jsonReader(joinPartyRequest{Code: "code", Name: "name"}))
	partyBroker := &mockPartyBroker{}
	partyBroker.joinPartyResponse = errors.New("random error")
	_, sessions, w := makeCall(req, partyBroker)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(500))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyBroker).To(Equal(mockPartyBroker{givenCode: "code", givenName: "name", joinPartyResponse: errors.New("random error")}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}
