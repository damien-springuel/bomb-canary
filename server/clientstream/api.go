package clientstream

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type websocketWriter func(messageType int, data []byte) error

type sessionGetter interface {
	Get(session string) (code string, name string, err error)
}

type clientBroker interface {
	Add(code, name string) chan string
}

type clientStreamServer struct {
	sessionGetter sessionGetter
	clientBroker  clientBroker
}

func Register(engine *gin.Engine, sessionGetter sessionGetter, clientBroker clientBroker) {

	clientStream := clientStreamServer{
		sessionGetter: sessionGetter,
		clientBroker:  clientBroker,
	}

	engine.GET("/events",
		createWebsocketConnection,
		clientStream.checkSession,
		clientStream.streamEvents,
	)
}

func createWebsocketConnection(c *gin.Context) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	c.Set("websocketWriter", websocketWriter(conn.WriteMessage))
	c.Next()
}

func getWebsocketWriterFromContext(c *gin.Context) websocketWriter {
	value, _ := c.Get("websocketWriter")
	return value.(websocketWriter)
}

func (s clientStreamServer) checkSession(c *gin.Context) {
	writer := getWebsocketWriterFromContext(c)
	session, err := c.Cookie("session")
	if err != nil {
		_ = writer(websocket.CloseMessage, websocket.FormatCloseMessage(4401, "no session cookie"))
		c.Abort()
		return
	}

	partyCode, playerName, err := s.sessionGetter.Get(session)
	if err != nil {
		_ = writer(websocket.CloseMessage, websocket.FormatCloseMessage(4403, "invalid session"))
		return
	}

	c.Set("code", partyCode)
	c.Set("name", playerName)

	c.Next()
}

func (s clientStreamServer) streamEvents(c *gin.Context) {
	write := getWebsocketWriterFromContext(c)
	code := c.GetString("code")
	name := c.GetString("name")
	out := s.clientBroker.Add(code, name)

	for {
		messageToSend, ok := <-out
		if !ok {
			break
		}
		_ = write(websocket.TextMessage, []byte(messageToSend))
	}
	err := write(websocket.CloseMessage, websocket.FormatCloseMessage(1000, fmt.Sprintf("hello %s from party %s", name, code)))
	fmt.Printf("write - err: %#v\n", err)
}
