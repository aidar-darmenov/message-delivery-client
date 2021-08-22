package model

type ClientParams struct {
	Id       string `json:"-"`
	HttpPort int    `json:"http_port"`
	Name     string `json:"name"`
}
