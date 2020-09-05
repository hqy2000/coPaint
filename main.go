package main

import (
	"coPaint/config"
	"coPaint/model"
	"coPaint/route"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	userId := rand.Intn(10); //strconv.Atoi(r.Header.Get("userId"))
	user := config.Users[userId]
	user.Device = conn
	conn.SetCloseHandler(func(code int, text string) error {
		user.Device = nil
		return nil
	})

	//go func() {
	//	var counter uint16 = 0;
	//	for range time.Tick(time.Second * 5) {
	//		if user
	//		conn.WriteJSON(model.HeartBeat{
	//			HeartBeatID: counter,
	//		})
	//		counter++
	//	}
	//}()

	for {
		var request model.Request = model.Request{}
		err := conn.ReadJSON(&request)

		if err == nil {
			jsonString, _ := json.Marshal(request.Data)
			switch request.Type {
			case "room":
				roomRequest := model.RoomRequest{}
				err := json.Unmarshal(jsonString, &roomRequest)
				fmt.Println(err)
				switch roomRequest.Operation {
				case "join":
					room := config.Rooms[uint16(roomRequest.RoomID)]
					if room.ID != 0 {
						user.InRoom = room
						room.Users = append(room.Users, user)
						conn.WriteJSON(model.RoomResponse{
							Room: *room,
							Status: "join_success",
						})
					}
					break
				case "create":
					room := model.Room{
						ID: uint(rand.Intn(999999)),
						Users: []model.User{user},
						Painting: config.Paintings[0],
						Touches: []model.Touch{},
					}
					user.InRoom = &room
					config.Rooms[uint16(room.ID)] = user.InRoom
					conn.WriteJSON(model.RoomResponse{
						Room: room,
						Status: "create_success",
					})
					break
				case "exit":
					user.InRoom.ID = 0
					break
				}

			case "event":
				event := model.Event{}
				err := json.Unmarshal(jsonString, &event)
				fmt.Println(err)
				fmt.Println(event.Touch.B)

				touch := event.Touch
				touch.User = user
				user.InRoom.Touches = append(user.InRoom.Touches, touch)

				for _, user := range user.InRoom.Users {
					fmt.Println(user.Username)
					user.Device.WriteJSON(touch)
				}
				break
			default:
				break
			}
		} else {
			fmt.Println(err)
		}

		//fmt.Println("Unknown operation sent by the client.")
	}
}


func main() {
	r := gin.Default()
	r.GET("/", route.Default)
	r.GET("/paintings/list", route.List)
	r.GET("/ws", func(c *gin.Context) {
		serveWs(c.Writer, c.Request)
	})
	r.Run()
}