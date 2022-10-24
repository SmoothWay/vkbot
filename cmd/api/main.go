package main

import (
	"context"
	"fmt"
	"log"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/l2x/dota2api"
)

func main() {
	token := "27e4657d1480eb5785cda6f8ca59167674483360417a792e92831c01499aafd99c2197dbe8ce239b62c3a"
	vk := api.NewVK(token)
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	ip, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	ip.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
		if obj.Message.Text == "zhanbot" {
			dota2, err := dota2api.LoadConfig("./config.ini")
			if err != nil {
				log.Fatal(err)
			}
			steamIds := []int64{
				76561199227027508,
			}
			accountId := dota2.GetAccountId(steamIds[0])
			param := map[string]interface{}{
				"account_id":        accountId,
				"matches_requested": "1",
			}
			matchHistory, err := dota2.GetMatchHistory(param)
			if err != nil {
				log.Fatal(err)
			}
			// log.Println(matchHistory)
			DBLink := fmt.Sprintf("https://dotabuff.com/matches/%d", matchHistory.Result.Matches[0].MatchID)
			// matchHistoryObject, _ := json.Marshal(matchHistory)
			b := params.NewMessagesSendBuilder()
			b.Message(DBLink)
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)

			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		}

	})

	log.Println("Start Long Poll")
	if err := ip.Run(); err != nil {
		log.Fatal(err)
	}
}
