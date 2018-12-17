package handlers

import (
    "github.com/gorilla/websocket"
    "middleware/models"
)

func SendMessage(conn *websocket.Conn, in *models.InData) {
//    url := "http://127.0.0.1:/"
//    fmt.Println("URL:>", url)
//
//    s := napping.Session{}
//    h := &http.Header{}
//    h.Set("X-Custom-Header", "myvalue")
//    s.Header = h
//
//    var jsonStr = []byte(`
//{
//    "title": "Buy cheese and bread for breakfast."
//}`)
//
//    var data map[string]json.RawMessage
//    err := json.Unmarshal(jsonStr, &data)
//    if err != nil {
//        fmt.Println(err)
//    }
//
//    resp, err := s.Post(url, &data, nil, nil)
//    if err != nil {
//        log.Fatal(err)
//    }
//    fmt.Println("response Status:", resp.Status())
//    fmt.Println("response Headers:", resp.HttpResponse().Header)
//    fmt.Println("response Body:", resp.RawText())
}