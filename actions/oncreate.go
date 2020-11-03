package actions

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"underkode.ru/secret-santa/application"
	"underkode.ru/secret-santa/store"
	"underkode.ru/secret-santa/utils"
)

var OnCreate = SaveLastActionDecorate(func(app application.ApplicationContext) func(message *tb.Message) {
	return func(message *tb.Message) {
		user := app.UserStore.FindByExternalId(utils.ToString(message.Sender.ID))

		if user == nil {
			fmt.Printf("User not found by %d", message.Sender.ID)
			return
		}

		round, err := app.RoundStore.Generate(store.GenerateRound{
			OwnerUser: user,
		})

		if err != nil {
			return
		}

		_, _ = app.Bot.Send(
			message.Chat,
			fmt.Sprintf("*Secret Santa* `%s` for _%d_ year", round.Code, round.Year),
			tb.ModeMarkdownV2,
		)
	}
})
