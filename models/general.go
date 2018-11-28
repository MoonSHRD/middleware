package models

type InData struct {
    Address string `json:"address"`
    Type    string `json:"type"`
    Data    string `json:"data"`
}

type OutData struct {
    Type    string `json:"type"`
    Data    string `json:"data"`
}

type Request struct {
    Type string `json:"type"`
    Data RequestData `json:"data"`
}

type RequestData struct {
    Hash string `json:"hash"`
}

type CentrifugoConfig struct {
    Secret string `json:"secret"`
}