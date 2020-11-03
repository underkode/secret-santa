package application

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"underkode.ru/secret-santa/store"
)

type ApplicationContext struct {
	WorkDirectory         string
	UserStore             *store.UserStore
	RoundStore            *store.RoundStore
	LastActionStore       *store.LastActionStore
	RoundParticipantStore *store.RoundParticipantStore
	Bot                   *tb.Bot
}
