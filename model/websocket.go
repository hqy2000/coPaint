package model

type RoomResponse struct {
	Room Room `json:"room"`
	Status string `json:"status"`
}

type Request struct {
	Type string `json:"type"` // "room" or "event"
	Data interface{} `json:"data"`
}

type RoomRequest struct {
	Operation string `json:"operation"`
	ID int `json:"Id"`
}

type Event struct {
	Touch Touch `json:"touch"`
	Room Room `json:"room"`
}

type HeartBeat struct {
	HeartBeatID uint16 `json:"heartBeatId"`
}

type Error struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}