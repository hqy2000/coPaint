package route

import (
	"coPaint/config"
	"coPaint/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"time"
)

func WebSocket(c *gin.Context) {
	ws(c.Writer, c.Request)
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func ws(w http.ResponseWriter, r *http.Request) {
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

	go func() {
		var counter uint16 = 0;
		for range time.Tick(time.Second * 5) {
			if user.Device != nil {
				conn.WriteJSON(model.HeartBeat{
					HeartBeatID: counter,
				})
				counter++
			}
		}
	}()

	for {
		var request model.Request = model.Request{}
		err := conn.ReadJSON(&request)

		if err == nil {
			jsonString, _ := json.Marshal(request.Data)
			switch request.Type {
			case "room":
				roomRequest := model.RoomRequest{}
				json.Unmarshal(jsonString, &roomRequest)
				// fmt.Println(err)
				switch roomRequest.Operation {
				case "join":
					room := config.Rooms[uint16(roomRequest.ID)]
					if room.ID != 0 {
						user.InRoom = room
						room.Users = append(room.Users, user)
						conn.WriteJSON(model.RoomResponse{
							Room: *room,
							Status: "join_success",
						})
					} else {
						conn.WriteJSON(model.Error{
							404, "room_not_found",
						})
					}
					break
				case "create":
					room := model.Room{
						ID: uint(rand.Intn(999999)),
						Users: []model.User{user},
						PaintingID: roomRequest.ID,
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
					us := user.InRoom.Users
					for _, u := range us {
						u.InRoom = nil
					}
					break
				}

			case "event":
				event := model.Event{}
				json.Unmarshal(jsonString, &event)


				touch := event.Touch
				touch.User = user

				if user.InRoom == nil {
					user.Device.WriteJSON(model.Error{
						Code: 404,
						Msg: "not_in_room",
					})
				} else {
					user.InRoom.Touches = append(user.InRoom.Touches, touch)
					for _, u := range user.InRoom.Users {
						u.Device.WriteJSON(touch)
					}
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
