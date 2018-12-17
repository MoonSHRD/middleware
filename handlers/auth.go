package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"middleware/models"
    "time"
)

type centConfInstance struct {
    Instance *models.CentrifugoConfig
}

var centConf = centConfInstance{}

func SetConf(conf *models.CentrifugoConfig) {
    centConf.Instance=conf
}

func VerifySig(from, sigHex string, msg []byte) bool {
	fromAddr := common.HexToAddress(from)

	sig := hexutil.MustDecode(sigHex)
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	if sig[64] != 27 && sig[64] != 28 {
		return false
	}
	sig[64] -= 27

	pubKey, err := crypto.SigToPub(signHash(msg), sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return fromAddr == recoveredAddr
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func GenerateNonce() string {
	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	nonce := string(fmt.Sprintf("%x", sha256.Sum256(data)))
	//db.GetRedis().Set(address+"_nonce",nonce,time.Minute*5)
	return nonce
}

func GenerateAuthRequest(nonce string) (int, []byte) {
	b, e := json.Marshal(models.Request{"auth", models.RequestData{nonce}})
	if e != nil {
		panic(e)
	}
	return 1, b
}

func GenerateNewJWTAnswer(address string) (int, []byte) {
	b, e := json.Marshal(models.OutData{"newJWT", generateJWT(address)})
	if e != nil {
		log.Println(e)
	}
	return 1, b
}

func HandleAuth(conn *websocket.Conn, in *models.InData, nonce string) bool {
	suc := VerifySig(in.Address, in.Data, []byte(nonce))
	log.Println(suc)
	if suc {
		b, e := json.Marshal(models.OutData{"auth_success", generateJWT(in.Address)})
		if e != nil {
			log.Println(e)
		}
		conn.WriteMessage(1, b)
	} else {
		conn.Close()
	}
	return suc
}

func generateJWT(address string) string {
    var mySigningKey = []byte(centConf.Instance.Secret)
    token := jwt.New(jwt.SigningMethodHS256)
    claims := make(jwt.MapClaims)
    claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
    //claims["iat"] = time.Now().Unix()
    claims["sub"] = address
    token.Claims = claims
    //fmt.Printf("Token for user %v expires %v", claims["user"], claims["exp"])
    tokenString, _ := token.SignedString(mySigningKey)
    return tokenString
}


func NewToken(conn *websocket.Conn, in *models.InData) {
    conn.WriteMessage(GenerateNewJWTAnswer(in.Address))
}