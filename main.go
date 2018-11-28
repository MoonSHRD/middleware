package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"middleware/core"
	"middleware/models"
	"os"
)

func main() {
	confPath := flag.String("c", "./config.json", "centrifugo config.json path")
	port := flag.String("p", "8080", "mw port")
	flag.Parse()

	jsonFile, err := os.Open(*confPath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result models.CentrifugoConfig
	json.Unmarshal(byteValue, &result)

	server := core.Server{Config: &result, Port: *port}
	server.Serve()
}
