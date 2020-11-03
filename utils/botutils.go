package utils

import tb "gopkg.in/tucnak/telebot.v2"

func HasMessageFromPrivateChat(upd *tb.Update) bool {
	return upd.Message != nil && isPrivateChat(upd)
}

func isPrivateChat(upd *tb.Update) bool {
	return upd.Message.Chat.Type == tb.ChatPrivate
}
