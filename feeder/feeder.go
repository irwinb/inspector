package feeder

import (
	"code.google.com/p/go.net/websocket"
	"container/list"
	"encoding/json"
	"github.com/irwinb/inspector/config"
	"github.com/irwinb/inspector/models"
	"log"
	"net/http"
	"sync"
)

type Connection struct {
	socket *websocket.Conn
	id     uint
}

var connections = list.New()
var connectionsMutex sync.Mutex

var messageQueue = make(chan interface{})
var lastId uint = 0

func InitializeFeeder() {
	http.Handle(config.FeederEndpoint, websocket.Handler(SockServer))
	go startFeeder()
}

func SockServer(w *websocket.Conn) {
	ele := addConnection(w)

	var message string
	for {
		if err := websocket.Message.Receive(w, &message); err != nil {
			removeConnection(ele)
			return
		}
		// Ignore messages.
	}
}

func startFeeder() {
	for {
		msg := <-messageQueue
		go feed(&msg)
	}
}

func FeedRequest(id int, p *models.Project, r *models.Request) {
	message := models.RequestMessage{
		TransactionId: id,
		MessageType:   "request",
		Project:       p,
		Request:       r}
	messageQueue <- message
}

func FeedResponse(id int, p *models.Project, r *models.Response) {
	message := models.ResponseMessage{
		TransactionId: id,
		MessageType:   "response",
		Project:       p,
		Response:      r}
	messageQueue <- message
}

func feed(r interface{}) {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatal("Error converto to JSON on Feed: ", err)
		return
	}

	connectionsMutex.Lock()
	removeThese := make([]*list.Element, 0, 0)
	for ele := connections.Front(); ele != nil; ele = ele.Next() {
		conn := ele.Value.(*Connection)
		_, err := conn.socket.Write(data)
		if err != nil {
			log.Fatal("Writing JSON failed: ", err)
			removeThese = append(removeThese, ele)
		}
	}
	connectionsMutex.Unlock()

	for _, v := range removeThese {
		removeConnection(v)
	}
}

func addConnection(w *websocket.Conn) *list.Element {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	lastId++
	conn := Connection{
		socket: w,
		id:     lastId,
	}
	log.Println("New connection ", conn)
	return connections.PushFront(&conn)
}

func removeConnection(ele *list.Element) {
	log.Println("Removing connection ", ele.Value)
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()
	connections.Remove(ele)
}
