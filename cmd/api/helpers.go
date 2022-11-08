package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/l2x/dota2api"
)

func newApp() (*App, error) {
	token := os.Getenv("TOKEN")
	vk := api.NewVK(token)
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		return nil, err
	}

	ip, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		return nil, err
	}
	dota2, err := dota2api.LoadConfig("./config.ini")
	if err != nil {
		return nil, err
	}
	heroes, err := dota2.GetHeroes()
	if err != nil {
		log.Fatal(err)
	}
	return &App{vk, ip, dota2, heroes}, nil
}

func (app *App) sendMessage(vk *api.VK, obj events.MessageNewObject, message string) {
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

func (app *App) getInfo(param map[string]interface{}, nick string, obj events.MessageNewObject) {
	var matchID int64
	var heroID int
	var hero string
	var pl dota2api.Player
	var match dota2api.MatchDetails
	for {
		var win = "LOSE >(("
		matchHistory, err := app.dota2.GetMatchHistory(param)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Minute * 2)
			continue
		}
		for i := range matchHistory.Result.Matches {
			if matchHistory.Result.Matches[i].LobbyType == 7 {
				match, err = app.dota2.GetMatchDetails(matchHistory.Result.Matches[i].MatchID)
				if err != nil {
					log.Println(err)
					time.Sleep(time.Minute * 2)
					continue
				}
				break
			}
		}

		if match.Result.MatchID != matchID && match.Result.MatchID != 0 {
			players := match.Result.Players
			for _, player := range players {
				if player.AccountID == int(param["account_id"].(int64)) {
					pl = player
					heroID = player.HeroID
					for _, v := range app.heroes {
						if v.ID == heroID {
							hero = v.Name
							break
						}
					}

					break
				}
			}

			if hero == "" {
				log.Println("Empty hero")
				time.Sleep(time.Minute * 2)
				continue
			}
			if (match.Result.RadiantWin && pl.PlayerSlot>>7&1 == 0) || (!match.Result.RadiantWin && pl.PlayerSlot>>7&1 == 1) {
				win = "WIN B-)"
			}
			duration := match.Result.Duration / 60
			info := fmt.Sprintf("%v\n%v %v\n%d-%d-%d | %v min\nhttps://dotabuff.com/matches/%d", nick, strings.Title(strings.Split(hero, "_")[3]), win, pl.Kills, pl.Deaths, pl.Assists, duration, match.Result.MatchID)
			app.sendMessage(app.vk, obj, info)
			matchID = match.Result.MatchID
		}

		time.Sleep(time.Minute * 2)
	}
}

func (app *App) makeScreenshot(site string) (string, error) {
	screenShotService := fmt.Sprintf("http://localhost:7171?url=%s%s", "https://", url.QueryEscape(site))

	resp, err := http.Get(screenShotService)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	filename := fmt.Sprintf("%s.png", strings.Replace(site, "/", "-", -1))
	err = os.WriteFile(filename, body, 0666)
	if err != nil {
		return "", err
	}
	log.Printf("....saved screenshot to file %s", filename)
	return filename, nil
}
