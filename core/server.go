package core

import (
    "encoding/json"
    "github.com/gorilla/websocket"
    "log"
    "middleware/handlers"
    //"moonshard_middleware/helpers"
    "middleware/models"
    "net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	Host   string
	Port   string
	Config *models.CentrifugoConfig
}

var r_handlers = map[string]func(*websocket.Conn, *models.InData) {
    "new_token": handlers.NewToken,
    "message": handlers.SendMessage,
}

func (s Server) handleRequest(conn *websocket.Conn, in *models.InData) {
    fHandler,err:=r_handlers[in.Type]
    if err == false {
        log.Println("wrong request")
        conn.Close()
        return
    }
    fHandler(conn,in)
}

func (s Server) handleConnect(conn *websocket.Conn) {
	nonce := handlers.GenerateNonce()
	conn.WriteMessage(handlers.GenerateAuthRequest(nonce))
	authorized := false

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType != 1 {
			conn.Close()
			return
		}

		var in models.InData
		json.Unmarshal(msg, &in)

		switch in.Type {
		case "auth":
			authorized = handlers.HandleAuth(conn, &in, nonce)
			break
		default:
			if !authorized {
				conn.Close()
				return
			} else {
				log.Println("handling request")
				s.handleRequest(conn, &in)
			}
			break
		}
	}
}

func (s Server) Serve() {
    handlers.SetConf(s.Config)
    
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("user connecting")
		conn, _ := upgrader.Upgrade(w, r, nil)
		go s.handleConnect(conn)
	})

	log.Println("Server starting at " + s.Host+":"+s.Port)
	http.ListenAndServe(s.Host+":"+s.Port, nil)
}
