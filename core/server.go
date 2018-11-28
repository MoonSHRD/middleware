package core

import (
    "encoding/json"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/websocket"
    "log"
    "moonshard_middleware/handlers"
    "time"
    
    //"moonshard_middleware/helpers"
    "moonshard_middleware/models"
    "net/http"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type Server struct {
    Port string
    Config *models.CentrifugoConfig
}


func (s Server) generateJWT(address string) string {
    var mySigningKey = []byte(s.Config.Secret)
    token := jwt.New(jwt.SigningMethodHS256)
    claims := make(jwt.MapClaims)
    claims["exp"] = time.Now().Add(time.Minute*5).Unix()
    //claims["iat"] = time.Now().Unix()
    claims["sub"] = address
    token.Claims = claims
    //fmt.Printf("Token for user %v expires %v", claims["user"], claims["exp"])
    tokenString, _ := token.SignedString(mySigningKey)
    return tokenString
}

func (s Server) handleRequest(conn *websocket.Conn,in *models.InData) {
    log.Println("sending answer")
    conn.WriteMessage(handlers.GenerateNewJWT(s.generateJWT(in.Address)))
}

func (s Server) handleConnect(conn *websocket.Conn) {
    nonce:=handlers.GenerateNonce()
    conn.WriteMessage(handlers.GenerateAuthRequest(nonce))
    authorized:=false
    
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
            authorized=handlers.HandleAuth(conn,&in,nonce,s.generateJWT(in.Address))
            break
        default:
            if !authorized {
                conn.Close()
                return
            } else {
                log.Println("handling request")
                s.handleRequest(conn,&in)
            }
            break
        }
        
        //if in.Type == "auth" {
        //    suc := handlers.VerifySig(in.Address, in.Data, []byte(nonce))
        //    log.Println(suc)
        //    if suc {
        //        authorized =true
        //        conn.WriteMessage(handlers.GenerateAuthSuccessRequest(in.Address))
        //    } else {
        //        conn.Close()
        //        return
        //    }
        //} else {
        //    log.Println("handling request")
        //    if !authorized {
        //        conn.Close()
        //        return
        //    } else {
        //        log.Println("handling request")
        //        handleRequest(conn,&in)
        //    }
        //}
    }
}

func (s Server) Serve() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Println("user connecting")
        conn, _ := upgrader.Upgrade(w, r, nil)
        go s.handleConnect(conn)
    })
    
    log.Println("Server starting at "+s.Port+" port...")
    http.ListenAndServe(":"+s.Port, nil)
}
