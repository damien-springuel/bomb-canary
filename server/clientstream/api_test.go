package clientstream

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

type mockClientBroker struct {
	channelToReturn chan []byte
	receivedCode    string
	receivedName    string
	closerCalled    bool
}

func (m *mockClientBroker) Add(code, name string) (chan []byte, func()) {
	m.receivedCode = code
	m.receivedName = name
	return m.channelToReturn, func() {
		m.closerCalled = true
	}
}

func setup(sessionGetter *mockSessionGetter, clientBroker *mockClientBroker, header http.Header) (*websocket.Conn, func()) {
	gin.SetMode(gin.TestMode)
	ginEngine := gin.New()
	Register(ginEngine, sessionGetter, clientBroker)
	s := httptest.NewServer(ginEngine)

	url := "ws" + strings.TrimPrefix(s.URL, "http") + "/events"
	ws, _, _ := websocket.DefaultDialer.Dial(url, header)

	return ws, func() {
		s.Close()
		ws.Close()
	}
}

func Test_StreamEvents(t *testing.T) {
	clientOut := make(chan []byte, 4)
	clientOut <- []byte("m1")
	clientOut <- []byte("m2")
	clientOut <- []byte("m3")
	clientOut <- []byte("m4")
	close(clientOut)
	clientBroker := &mockClientBroker{channelToReturn: clientOut}
	sessionGetter := &mockSessionGetter{}

	header := http.Header{}
	header.Add("Cookie", "session=testSession")
	conn, closer := setup(sessionGetter, clientBroker, header)
	defer closer()

	actualMessages := [][]byte{}
	g := NewWithT(t)
	for {
		_, nextMessage, err := conn.ReadMessage()
		if err != nil {
			closeError, ok := err.(*websocket.CloseError)
			g.Expect(ok).To(BeTrue())
			g.Expect(closeError.Code).To(Equal(1000))
			break
		}
		actualMessages = append(actualMessages, nextMessage)
	}

	g.Expect(actualMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
	}))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(clientBroker.receivedCode).To(Equal("testCode"))
	g.Expect(clientBroker.receivedName).To(Equal("testName"))
	g.Expect(clientBroker.closerCalled).To(BeTrue())
}

func Test_StreamEvents_CloseConnectionWith4401IfNoSessionCookie(t *testing.T) {
	clientBroker := &mockClientBroker{}
	sessionGetter := &mockSessionGetter{}

	conn, closer := setup(sessionGetter, clientBroker, http.Header{})
	defer closer()

	_, _, err := conn.ReadMessage()
	g := NewWithT(t)
	closeError, ok := err.(*websocket.CloseError)
	g.Expect(ok).To(BeTrue())
	g.Expect(closeError.Code).To(Equal(4401))

	g.Expect(sessionGetter.receivedSession).To(BeEmpty())
	g.Expect(clientBroker.receivedCode).To(BeEmpty())
	g.Expect(clientBroker.receivedName).To(BeEmpty())
}

func Test_StreamEvents_CloseConnectionWith4403IfSessionIsInvalid(t *testing.T) {
	clientBroker := &mockClientBroker{}
	sessionGetter := &mockSessionGetter{}
	sessionGetter.getError = fmt.Errorf("invalid session")

	header := http.Header{}
	header.Add("Cookie", "session=testSession")
	conn, closer := setup(sessionGetter, clientBroker, header)
	defer closer()

	_, _, err := conn.ReadMessage()

	g := NewWithT(t)
	closeError, ok := err.(*websocket.CloseError)
	g.Expect(ok).To(BeTrue())
	g.Expect(closeError.Code).To(Equal(4403))

	g.Expect(sessionGetter.receivedSession).To(Equal("testSession"))
	g.Expect(clientBroker.receivedCode).To(BeEmpty())
	g.Expect(clientBroker.receivedName).To(BeEmpty())
}
