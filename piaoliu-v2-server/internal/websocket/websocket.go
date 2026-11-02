package websocket

import (
	"drift-bottle-server/internal/middleware"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 表示一个 WebSocket 客户端
type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub 管理所有 WebSocket 连接
type Hub struct {
	clients    map[int64]*Client
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
}

var hub = &Hub{
	clients:    make(map[int64]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func init() {
	go hub.run()
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("WebSocket: user %d connected", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket: user %d disconnected", client.UserID)
		}
	}
}

// SendMessage 向指定用户推送消息
func SendMessage(userID int64, message []byte) {
	hub.mu.RLock()
	client, ok := hub.clients[userID]
	hub.mu.RUnlock()

	if ok {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			hub.mu.Lock()
			delete(hub.clients, userID)
			hub.mu.Unlock()
		}
	}
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(c *gin.Context) {
	// 从 URL 参数获取 token
	token := c.Query("token")
	if token == "" {
		c.JSON(401, gin.H{"code": 401, "message": "未提供token"})
		return
	}

	claims, err := middleware.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"code": 401, "message": "token无效"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	hub.register <- client

	// 读写协程
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
		// 心跳响应等处理
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}
