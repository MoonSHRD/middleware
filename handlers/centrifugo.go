package handlers

import (
    "github.com/dgrijalva/jwt-go"
    "time"
)

func GenerateJWT(address string) string {
    var mySigningKey = []byte("secret")
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

//func GenerateJWTSend()(int,[]byte){
//    token:=GenerateJWT()
//    b,e:=json.Marshal(models.OutData{"JWTtoken",token})
//    if e!=nil {
//        log.Println(e)
//    }
//    return 1,b
//}
