package handler

import (
	"encoding/json"
	"huddle-ws-server/rd"
	"huddle-ws-server/types"
	"huddle-ws-server/ws"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func StartRedisListener() {
	onlineStatusCh := rd.Subscribe("user_online_status")
	events := rd.Subscribe("broadcast")

	go processChannel(onlineStatusCh, handleOnlineStatus)
	go processChannel(events, handleMessage)
}

func processChannel(subscription <-chan *redis.Message, handler func(msg interface{})) {
	for msg := range subscription {

		err := json.Unmarshal([]byte(msg.Payload), &msg)
		if err != nil {
			continue
		}

		handler(msg.Payload)
	}
}

func handleMessage(payload interface{}) {
	var broadcastPayload types.Message
	if err := json.Unmarshal([]byte(payload.(string)), &broadcastPayload); err != nil {
		return
	}

	ws.WsManager.BroadcastMessage(broadcastPayload)
}

func handleOnlineStatus(payload interface{}) {
	var userOnlineStatus struct {
		UserId string `json:"userId"`
		Status string `json:"status"`
	}

	payloadBytes, _ := json.Marshal(payload)
	if err := json.Unmarshal(payloadBytes, &userOnlineStatus); err != nil {
		return
	}

	if userOnlineStatus.UserId == "" || userOnlineStatus.Status == "" {
		return
	}

	ws.WsManager.BroadcastUserStatus(uuid.MustParse(userOnlineStatus.UserId), userOnlineStatus.Status)
}
