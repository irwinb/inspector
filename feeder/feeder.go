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

var connections = list.New()
var connectionsMutex sync.Mutex

var messageQueue = make(chan models.Request)

func InitializeFeeder() {
	http.Handle(config.FeederEndpoint, websocket.Handler(SockServer))
	go startFeeder()
}

func SockServer(w *websocket.Conn) {
	ele := connections.PushFront(w)

	var message string
	for {
		if err := websocket.Message.Receive(w, &message); err != nil {
			connectionsMutex.Lock()
			defer connectionsMutex.Unlock()

			connections.Remove(ele)
			return
		}
		// Ignore messages.
	}
}

func startFeeder() {
	for {
		req := <-messageQueue
		go Feed(&req)
	}
}

func Feed(r *models.Request) {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatal("Error converto to JSON on Feed: ", err)
		return
	}

	removeThese := make([]*list.Element, 0, 0)
	for conn := connections.Front(); conn != nil; conn = conn.Next() {
		w := conn.Value.(*websocket.Conn)
		_, err := w.Write(data)
		if err != nil {
			log.Fatal("Writing JSON failed: ", err)
			removeThese = append(removeThese, conn)
		}
	}

	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()
	for _, v := range removeThese {
		connections.Remove(v)
	}
}
