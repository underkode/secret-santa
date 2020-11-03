package actions

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	"underkode.ru/secret-santa/application"
	"underkode.ru/secret-santa/store"
	"underkode.ru/secret-santa/utils"
)

func OnText(app application.ApplicationContext) func(message *tb.Message) {
	return func(message *tb.Message) {
		user := app.UserStore.FindByExternalId(utils.ToString(message.Sender.ID))

		if user == nil {
			fmt.Printf("User not found by %d", message.Sender.ID)
			return
		}

		lastAction := app.LastActionStore.FindByUserId(user.Id)

		switch lastAction.Action {
		case Join:
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

			_, err := app.RoundParticipantStore.Join(store.JoinRound{
				Round:           round,
				ParticipantUser: user,
			})

			if err == nil {
				app.Bot.Send(
					message.Sender,
					fmt.Sprintf("You joined to *Secret Santa* `%s`", round.Code),
					tb.ModeMarkdownV2,
				)
			} else {
				app.Bot.Send(
					message.Sender,
					fmt.Sprintf("Joining error to *Secret Santa* `%s`", round.Code),
					tb.ModeMarkdownV2,
				)
			}
		case PlayOut:
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

			participants := app.RoundParticipantStore.FindAllByRoundId(round.Id)

			var users []store.User

			for _, participant := range participants {
				users = append(users, *app.UserStore.FindById(participant.UserId))
			}

			playOutResult := playOut(len(users))

			for secretSantaIndex, goodBoyIndex := range playOutResult {
				secretSanta := users[secretSantaIndex]
				goodBoy := users[goodBoyIndex]

				_, err := app.Bot.Send(
					&tb.User{ID: utils.ToInt(secretSanta.ExternalId)},
					fmt.Sprintf("*Secret Santa* `%s` was be play out:\n"+
						"You are *Secret Santa* for @%s", round.Code, utils.EscapeMarkdown(goodBoy.Username)),
					tb.ModeMarkdownV2,
				)

				if err != nil {
					log.Panic(err)
				}
			}
		default:
			_, err := app.Bot.Reply(message, message.Text)

			if err != nil {
				return
			}
		}
	}
}
