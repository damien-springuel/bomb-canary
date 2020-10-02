package clientstream

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const outKey = "out"

type websocketWriter func(messageType int, data []byte) error

type sessionGetter interface {
	Get(session string) (code string, name string, err error)
}

type clientBroker interface {
	Add(code, name string) (chan []byte, func())
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
		c.Abort()
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
		c.Abort()
		return
	}
	out, closer := s.clientBroker.Add(partyCode, playerName)
	defer closer()

	c.Set(outKey, out)

	c.Next()
}

func getOutFromContext(c *gin.Context) chan []byte {
	out, _ := c.Get(outKey)
	return out.(chan []byte)
}

func (s clientStreamServer) streamEvents(c *gin.Context) {
	write := getWebsocketWriterFromContext(c)
	out := getOutFromContext(c)

	for messageToSend := range out {
		err := write(websocket.TextMessage, messageToSend)
		if err != nil {
			c.Abort()
			return
		}
	}

	_ = write(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""))
}
