package actions

import (
	"errors"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/rand"
	"strings"
	"time"
	"underkode.ru/secret-santa/application"
	"underkode.ru/secret-santa/store"
	"underkode.ru/secret-santa/utils"
)

var OnPlayOut = SaveLastActionDecorate(func(app *application.ApplicationContext) func(message *tb.Message) {
	return func(message *tb.Message) {
		_, _ = app.Bot.Send(message.Chat, "Enter code of *Secret Santa*:", tb.ModeMarkdownV2)
	}
})

func PlayOut(app *application.ApplicationContext, message *tb.Message, user *store.User) {
	code := strings.Trim(message.Text, " ")
	round := app.RoundStore.FindByCode(code)

	if round == nil {
		app.Bot.Send(
			message.Sender,
			fmt.Sprintf("*Secret Santa* not found by `%s`", message.Text),
			tb.ModeMarkdownV2,
		)

		return
	}

	if round.OwnerUserId != user.Id {
		app.Bot.Send(
			message.Sender,
			fmt.Sprintf("You are not owner of *Secret Santa* `%s`", round.Code),
			tb.ModeMarkdownV2,
		)

		return
	}

	if app.RoundPlayOutStore.ExistsByRoundId(round.Id) {
		app.Bot.Send(
			message.Sender,
			fmt.Sprintf("*Secret Santa* `%s` play out already", round.Code),
			tb.ModeMarkdownV2,
		)

		return
	}

	participants := app.RoundParticipantStore.FindAllByRoundId(round.Id)

	var users []store.User

	for _, participant := range participants {
		users = append(users, *app.UserStore.FindById(participant.UserId))
	}

	playOutResult, err := playOut(len(users))

	if err != nil && err.Error() == "participants count should be greater 1" {
		app.Bot.Send(
			message.Sender,
			fmt.Sprintf("*Secret Santa* `%s` can't play out for one participant", round.Code),
			tb.ModeMarkdownV2,
		)

		return
	}

	var userPairs []store.RoundPlayOutUserPair
	for secretSantaIndex, goodBoyIndex := range playOutResult {
		secretSanta := users[secretSantaIndex]
		goodBoy := users[goodBoyIndex]

		userPairs = append(userPairs, store.RoundPlayOutUserPair{
			SecretSantaUser: &secretSanta,
			KidUser:         &goodBoy,
		})
	}

	_, err = app.RoundPlayOutStore.Create(store.CreateRoundPlayOut{
		Round: round,
		Pairs: userPairs,
	})

	if err != nil {
		println("Error of creation round play out:", err)
	}

	for _, pair := range userPairs {
		_, err := app.Bot.Send(
			&tb.User{ID: utils.ToInt(pair.SecretSantaUser.ExternalId)},
			fmt.Sprintf("*Secret Santa* `%s` was be play out:\n"+
				"You are *Secret Santa* for @%s", round.Code, utils.EscapeMarkdown(pair.KidUser.Username)),
			tb.ModeMarkdownV2,
		)

		if err != nil {
			log.Panic(err)
		}
	}
}

func playOut(participantCount int) ([]int, error) {
	if participantCount < 2 {
		return []int{}, errors.New("participants count should be greater 1")
	}

	var played []int
	preLastIndex := participantCount - 2
	lastIndex := participantCount - 1

	rand.Seed(time.Now().UTC().UnixNano())
	for index := 0; index < participantCount; index++ {
		playing := notPlayed(index, participantCount, played)

		if index == preLastIndex && !contains(played, lastIndex) {
			played = append(played, lastIndex)
			continue
		}

		played = append(played, playing[rand.Intn(len(playing))])
	}

	return played, nil
}

func notPlayed(index int, participantCount int, played []int) []int {
	var playing []int

	for i := 0; i < participantCount; i++ {
		if i == index || contains(played, i) {
			continue
		}
		playing = append(playing, i)
	}

	return playing
}

func contains(slice []int, value int) bool {
	for _, it := range slice {
		if it == value {
			return true
		}
	}

	return false
}
