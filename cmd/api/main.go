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
	var matchID int64 = 6824792266
	var heroID int
	var win = "LOSE >(("
	var hero dota2api.Hero
	var zhanbot dota2api.Player
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
	heroes, err := dota2.GetHeroes()
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
						log.Println(err)
						time.Sleep(time.Minute * 2)
						continue
					}
					log.Println(matchHistory.Result.Matches[0].MatchID)
					match, err := dota2.GetMatchDetails(matchHistory.Result.Matches[0].MatchID)
					if err != nil {
						log.Println(err)
						time.Sleep(time.Minute * 2)
						continue
					}
					if match.Result.MatchID != matchID {
						players := match.Result.Players
						for _, player := range players {
							if player.AccountID == int(accountId) {
								zhanbot = player
								heroID = player.HeroID
								hero = heroes[heroID-1]
								break
							}
						}
						if hero == (dota2api.Hero{}) {
							log.Println("Empty hero")
							time.Sleep(time.Minute * 2)
							continue
						}
						if match.Result.RadiantWin && zhanbot.PlayerSlot < 6 {
							win = "WIN B-)"
						}
						duration := match.Result.Duration / 60
						info := fmt.Sprintf("%v %v\n%d-%d-%d | %v min\nhttps://dotabuff.com/matches/%d", hero.Name[14:], win, zhanbot.Kills, zhanbot.Deaths, zhanbot.Assists, duration, match.Result.MatchID)
						sendMessage(vk, obj, info)
						matchID = match.Result.MatchID
					}
					time.Sleep(time.Minute * 2)
				}
			}()
		} else if obj.Message.Text == "lm" {
			matchHistory, err := dota2.GetMatchHistory(param)
			if err != nil {
				log.Println(err)
				return
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
		log.Println(err)
		return
	}
}

// func formatOutput() {
// 	duration := match.Result.Duration / 60
// 	info := fmt.Sprintf("%v %v\n%d-%d-%d | %v min\nhttps://dotabuff.com/matches/%d", hero.Name[14:], win, zhanbot.Kills, zhanbot.Deaths, zhanbot.Assists, duration, match.Result.MatchID)
// 	sendMessage(vk, obj, info)
// }
