package clientstream

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	outKey             = "out"
	connClosedKey      = "connClosed"
	websocketWriterKey = "websocketWriter"
)

type websocketWriter func(messageType int, data []byte) error

type sessionGetter interface {
	Get(session string) (name string, err error)
}

type clientBroker interface {
	Add(name string) (chan []byte, func())
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
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // temporary while developing frontend
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Abort()
		return
	}
	connClosed := make(chan bool)
	c.Set(connClosedKey, connClosed)
	go func() {
		defer conn.Close()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
		close(connClosed)
	}()

	c.Set(websocketWriterKey, websocketWriter(conn.WriteMessage))
	c.Next()
}

func getWebsocketWriterFromContext(c *gin.Context) websocketWriter {
	value, _ := c.Get(websocketWriterKey)
	return value.(websocketWriter)
}

func getConnClosedFromContext(c *gin.Context) chan bool {
	value, _ := c.Get(connClosedKey)
	return value.(chan bool)
}

func (s clientStreamServer) checkSession(c *gin.Context) {
	writer := getWebsocketWriterFromContext(c)
	session, err := c.Cookie("session")
	if err != nil {
		_ = writer(websocket.CloseMessage, websocket.FormatCloseMessage(4401, "no session cookie"))
		c.Abort()
		return
	}

	playerName, err := s.sessionGetter.Get(session)
	if err != nil {
		_ = writer(websocket.CloseMessage, websocket.FormatCloseMessage(4403, "invalid session"))
		c.Abort()
		return
	}
	out, closeClientStream := s.clientBroker.Add(playerName)
	go func() {
		connClosed := getConnClosedFromContext(c)
		<-connClosed
		closeClientStream()
	}()

	c.Set(outKey, out)

	c.Next()
	closeClientStream()
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
