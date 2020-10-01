package lobby

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

type mockPartyService struct {
	givenCode string
	givenName string
}

func (m *mockPartyService) CreateParty() string {
	return "testCode"
}

func (m *mockPartyService) JoinParty(code string, name string) {
	m.givenCode = code
	m.givenName = name
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

func makeCall(req *http.Request) (*mockPartyService, *mockSession, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	ginEngine := gin.New()

	partyService := &mockPartyService{}
	sessions := &mockSession{}
	Register(ginEngine, partyService, sessions)

	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	return partyService, sessions, w
}

func Test_CreateParty(t *testing.T) {
	req, _ := http.NewRequest("GET", "/lobby/create-party", jsonReader(createPartyRequest{Name: "testName"}))
	partyService, sessions, w := makeCall(req)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(200))

	actualResponse := createPartyResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &actualResponse)
	g.Expect(actualResponse).To(Equal(createPartyResponse{Code: "testCode"}))

	actualCookie := w.Result().Cookies()[0]
	g.Expect(actualCookie.Name).To(Equal("session"))
	g.Expect(actualCookie.Value).To(Equal("testSessionId"))
	g.Expect(actualCookie.MaxAge).To(Equal(int((time.Hour * 3).Seconds())))

	g.Expect(*partyService).To(Equal(mockPartyService{givenCode: "testCode", givenName: "testName"}))
	g.Expect(*sessions).To(Equal(mockSession{givenCode: "testCode", givenName: "testName"}))
}

func Test_CreateParty_ShouldReturn400IfNameIsAbsent(t *testing.T) {
	req, _ := http.NewRequest("GET", "/lobby/create-party", jsonReader(createPartyRequest{Name: ""}))
	partyService, sessions, w := makeCall(req)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))
	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyService).To(Equal(mockPartyService{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_CreateParty_ShouldReturn400IfMalformedBody(t *testing.T) {
	req, _ := http.NewRequest("GET", "/lobby/create-party", strings.NewReader("garbage"))
	partyService, sessions, w := makeCall(req)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyService).To(Equal(mockPartyService{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_JoinParty(t *testing.T) {
	req, _ := http.NewRequest("GET", "/lobby/join-party", jsonReader(joinPartyRequest{Code: "joinCode", Name: "testName"}))
	_, sessions, w := makeCall(req)

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
	req, _ := http.NewRequest("GET", "/lobby/join-party", jsonReader(joinPartyRequest{Code: "", Name: "testName"}))
	partyService, sessions, w := makeCall(req)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyService).To(Equal(mockPartyService{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_JoinParty_Should400IfNameAbsent(t *testing.T) {
	req, _ := http.NewRequest("GET", "/lobby/join-party", jsonReader(joinPartyRequest{Code: "joinCode", Name: ""}))
	partyService, sessions, w := makeCall(req)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyService).To(Equal(mockPartyService{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}

func Test_JoinParty_Should400IfMalformedBody(t *testing.T) {
	req, _ := http.NewRequest("GET", "/lobby/join-party", strings.NewReader("garbage"))
	partyService, sessions, w := makeCall(req)

	g := NewWithT(t)
	g.Expect(w.Code).To(Equal(400))

	g.Expect(w.Result().Cookies()).To(BeEmpty())

	g.Expect(*partyService).To(Equal(mockPartyService{}))
	g.Expect(*sessions).To(Equal(mockSession{}))
}
