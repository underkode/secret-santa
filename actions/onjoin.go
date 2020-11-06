package actions

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"underkode.ru/secret-santa/application"
	"underkode.ru/secret-santa/store"
)

var OnJoin = SaveLastActionDecorate(func(app *application.ApplicationContext) func(message *tb.Message) {
	return func(message *tb.Message) {
		_, _ = app.Bot.Send(message.Chat, "Enter code of *Secret Santa*:", tb.ModeMarkdownV2)
	}
})

func Join(app *application.ApplicationContext, message *tb.Message, user *store.User) {
	code := strings.Trim(message.Text, " ")
	round := app.RoundStore.FindByCode(code)

	if round == nil {
		_, _ = app.Bot.Send(
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
		_, _ = app.Bot.Send(
			message.Sender,
			fmt.Sprintf("You joined to *Secret Santa* `%s`", round.Code),
			tb.ModeMarkdownV2,
		)
	} else {
		_, _ = app.Bot.Send(
			message.Sender,
			fmt.Sprintf("Joining error to *Secret Santa* `%s`", round.Code),
			tb.ModeMarkdownV2,
		)
	}
}
