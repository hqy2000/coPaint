package model

import "github.com/gorilla/websocket"

type Room struct {
	ID uint `json:"id"`
	Painting Painting `json:"painting"`
	Users []User `json:"users"`
	Touches []Touch `json:"touches"`
}

type Touch struct {
	X uint16 `json:"x"`
	Y uint16 `json:"y"`
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	User User `json:"user"`
}

type User struct {
	ID uint8 `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	InRoom *Room `json:"-"`
	Device *websocket.Conn `json:"-"`
}

type Painting struct {
	ID uint8 `json:"id"`
	Name string `json:"name"`
	Image string `json:"-"`
}