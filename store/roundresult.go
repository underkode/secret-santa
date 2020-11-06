package store

import (
	"errors"
	"github.com/google/uuid"
	"underkode.ru/secret-santa/utils"
)

type RoundPlayOutPair struct {
	SecretSantaUserId uuid.UUID `json:"secretSantaUserId"`
	KidUserId         uuid.UUID `json:"kidUserId"`
}

type RoundPlayOut struct {
	RoundId uuid.UUID          `json:"roundId"`
	Pairs   []RoundPlayOutPair `json:"pairs"`
}

type RoundPlayOutStore struct {
	filename string
	items    []RoundPlayOut
}

func NewRoundPlayOutStore(filename string) (*RoundPlayOutStore, error) {
	var items []RoundPlayOut

	utils.CreateIfNotExists(filename)

	err := utils.UnmarshalFileJson(filename, &items)

	if err != nil {
		return nil, err
	}

	return &RoundPlayOutStore{
		filename: filename,
		items:    items,
	}, nil
}

func (self *RoundPlayOutStore) find(checkFunc func(round RoundPlayOut) bool) *RoundPlayOut {
	for index, it := range self.items {
		if checkFunc(it) {
			return &self.items[index]
		}
	}

	return nil
}

func (self *RoundPlayOutStore) save() error {
	if err := utils.MarshalFileJson(self.filename, &self.items); err != nil {
		return err
	}

	return nil
}

type RoundPlayOutUserPair struct {
	SecretSantaUser *User
	KidUser         *User
}

type CreateRoundPlayOut struct {
	Round *Round
	Pairs []RoundPlayOutUserPair
}

func (self *RoundPlayOutStore) Create(arg CreateRoundPlayOut) (*RoundPlayOut, error) {
	item := self.findByRoundId(arg.Round.Id)

	if item != nil {
		return item, errors.New("play out already exists")
	}

	var pairs []RoundPlayOutPair
	for _, pair := range arg.Pairs {
		pairs = append(pairs, RoundPlayOutPair{
			SecretSantaUserId: pair.SecretSantaUser.Id,
			KidUserId:         pair.KidUser.Id,
		})
	}

	item = &RoundPlayOut{
		RoundId: arg.Round.Id,
		Pairs:   pairs,
	}
	self.items = append(self.items, *item)

	err := self.save()

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (self *RoundPlayOutStore) findAll(checkFunc func(round RoundPlayOut) bool) []RoundPlayOut {
	var items []RoundPlayOut

	for _, it := range self.items {
		if checkFunc(it) {
			items = append(items, it)
		}
	}

	return items
}

func (self *RoundPlayOutStore) findByRoundId(roundId uuid.UUID) *RoundPlayOut {
	return self.find(func(round RoundPlayOut) bool {
		return round.RoundId == roundId
	})
}

func (self *RoundPlayOutStore) ExistsByRoundId(roundId uuid.UUID) bool {
	return self.findByRoundId(roundId) != nil
}
