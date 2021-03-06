package fvb

import (
	"fmt"
	"github.com/Haski007/fav-videos/internal/fvb/resource"
	"github.com/Haski007/fav-videos/pkg/graceshut"
	"github.com/Haski007/go-errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"time"
)

func StartBot(bot *resource.FVBService) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.Bot.GetUpdatesChan(u)
	if err != nil {
		errors.Println(err)
		return
	}
	defer func() {
		if errR := recover(); errR != nil {
			_, err := bot.Bot.Send(
				tgbotapi.NewMessage(
					bot.Cfg.CreatorID,
					fmt.Sprintf("[Main panic] err: %+v\n", errR)))
			if err != nil {
				logrus.Fatalf("[defer panic] err: %s", err)
			}
		}
	}()

	go bot.HandleRoutes(updates)
	go tiktokLoop(bot)

	graceshut.Loop()
}

func tiktokLoop(bot *resource.FVBService) {
	for {
		bot.CheckNewVideos()
		fmt.Println("Loop passed!")
		time.Sleep(5 * time.Second)
	}
}
