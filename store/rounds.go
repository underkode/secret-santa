package store

import (
	"github.com/google/uuid"
	"time"
	"underkode.ru/secret-santa/randomstring"
	"underkode.ru/secret-santa/utils"
)

type Round struct {
	Id          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	OwnerUserId uuid.UUID `json:"ownerUserId"`
	Year        int       `json:"year"`
}

type RoundStore struct {
	filename string
	items    []Round
}

func NewRoundStore(filename string) (*RoundStore, error) {
	var items []Round

	utils.CreateIfNotExists(filename)

	err := utils.UnmarshalFileJson(filename, &items)

	if err != nil {
		return nil, err
	}

	return &RoundStore{
		filename: filename,
		items:    items,
	}, nil
}

func (self *RoundStore) find(checkFunc func(round Round) bool) *Round {
	for index, it := range self.items {
		if checkFunc(it) {
			return &self.items[index]
		}
	}

	return nil
}

func (self *RoundStore) FindByCode(code string) *Round {
	return self.find(func(round Round) bool {
		return round.Code == code
	})
}

func (self *RoundStore) save() error {
	if err := utils.MarshalFileJson(self.filename, &self.items); err != nil {
		return err
	}

	return nil
}

type GenerateRound struct {
	OwnerUser *User
}

func (self *RoundStore) Generate(arg GenerateRound) (*Round, error) {
	item := Round{
		Id:          uuid.New(),
		OwnerUserId: arg.OwnerUser.Id,
		Code:        randomstring.RandomString(8),
		Year:        time.Now().Year(),
	}
	self.items = append(self.items, item)

	err := self.save()

	if err != nil {
		return nil, err
	}

	return &item, nil
}
