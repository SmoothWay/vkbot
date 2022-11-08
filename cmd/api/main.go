package main

import (
	"context"
	"log"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/l2x/dota2api"
)

type App struct {
	vk     *api.VK
	ip     *longpoll.LongPoll
	dota2  dota2api.Dota2
	heroes []dota2api.Hero
}

func main() {
	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}
	inChat := make(map[int]int)
	steamIds := []int64{
		76561199227027508, // isBack
		76561198194205802, // 1thousand
		76561198117214277, // unsiz
		76561198194088028, // no voice
	}
	nickNames := []string{"Gama", "Gama(main)", "Kaba", "t1mon"}
	app.ip.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		_, ok := inChat[obj.Message.PeerID]

		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
		if obj.Message.Text == "start" && !ok {
			for i := range steamIds {
				accountID := app.dota2.GetAccountId(steamIds[i])
				param := map[string]interface{}{
					"account_id":        accountID,
					"game_mode":         7,
					"matches_requested": "5",
				}
				inChat[obj.Message.PeerID] = 1
				go app.getInfo(param, nickNames[i], obj)
			}

		}
	})
	log.Println("Start Long Poll")
	if err := app.ip.Run(); err != nil {
		log.Fatal(err)
	}
}
