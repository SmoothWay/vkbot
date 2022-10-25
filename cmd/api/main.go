package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/l2x/dota2api"
)

func main() {
	var matchID int64 = 6821785352
	token := os.Getenv("TOKEN")
	vk := api.NewVK(token)
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	ip, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}
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

	go ip.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
		if obj.Message.Text == "start" {

			go func() {
				for {
					matchHistory, err := dota2.GetMatchHistory(param)
					if err != nil {
						log.Fatal(err)
					}
					log.Println(matchHistory.Result.Matches[0].MatchID)
					if matchHistory.Result.Matches[0].MatchID != matchID {
						DBLink := fmt.Sprintf("https://dotabuff.com/matches/%d", matchHistory.Result.Matches[0].MatchID)
						// matchHistoryObject, _ := json.Marshal(matchHistory)
						sendMessage(vk, obj, DBLink)
						matchID = matchHistory.Result.Matches[0].MatchID
						time.Sleep(time.Minute * 9)
					}
					time.Sleep(time.Minute)
				}
			}()
		} else if obj.Message.Text == "lm" {
			matchHistory, err := dota2.GetMatchHistory(param)
			if err != nil {
				log.Fatal(err)
			}
			message := fmt.Sprintf("https://dotabuff.com/matches/%d", matchHistory.Result.Matches[0].MatchID)
			sendMessage(vk, obj, message)
		}
	})

	log.Println("Start Long Poll")
	if err := ip.Run(); err != nil {
		log.Fatal(err)
	}
}

func sendMessage(vk *api.VK, obj events.MessageNewObject, message string) {
	b := params.NewMessagesSendBuilder()
	b.Message(message)
	b.RandomID(0)
	if obj.Message.PeerID <= 100000000 {
		b.ChatID(obj.Message.PeerID)
	} else {
		b.PeerID(obj.Message.PeerID)
	}
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
}
