package actions

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"underkode.ru/secret-santa/application"
)

var OnJoin = SaveLastActionDecorate(func(app application.ApplicationContext) func(message *tb.Message) {
	return func(message *tb.Message) {
		_, _ = app.Bot.Send(message.Chat, "Enter code of *Secret Santa*:", tb.ModeMarkdownV2)
	}
})
