package actions

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"math/rand"
	"time"
	"underkode.ru/secret-santa/application"
)

var OnPlayOut = SaveLastActionDecorate(func(app application.ApplicationContext) func(message *tb.Message) {
	return func(message *tb.Message) {
		_, _ = app.Bot.Send(message.Chat, "Enter code of *Secret Santa*:", tb.ModeMarkdownV2)
	}
})

func playOut(participantCount int) []int {
	if participantCount < 2 {
		return []int{}
	}

	var played []int

	rand.Seed(time.Now().UTC().UnixNano())
	for index := 0; index < participantCount; index++ {
		var playing []int
		for i := 0; i < participantCount; i++ {
			if i == index || contains(played, i) {
				continue
			}
			playing = append(playing, i)
		}

		played = append(played, playing[rand.Intn(len(playing))])
	}

	return played
}

func contains(slice []int, value int) bool {
	for _, it := range slice {
		if it == value {
			return true
		}
	}

	return false
}
