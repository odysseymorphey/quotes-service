package models

type Quote struct {
	Id     string `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}
