package store

import (
	"github.com/google/uuid"
	"underkode.ru/secret-santa/utils"
)

type RoundParticipant struct {
	RoundId uuid.UUID `json:"roundId"`
	UserId  uuid.UUID `json:"userId"`
}

type RoundParticipantStore struct {
	filename string
	items    []RoundParticipant
}

func NewRoundParticipantStore(filename string) (*RoundParticipantStore, error) {
	var items []RoundParticipant

	utils.CreateIfNotExists(filename)

	err := utils.UnmarshalFileJson(filename, &items)

	if err != nil {
		return nil, err
	}

	return &RoundParticipantStore{
		filename: filename,
		items:    items,
	}, nil
}

func (self *RoundParticipantStore) find(checkFunc func(round RoundParticipant) bool) *RoundParticipant {
	for index, it := range self.items {
		if checkFunc(it) {
			return &self.items[index]
		}
	}

	return nil
}

func (self *RoundParticipantStore) save() error {
	if err := utils.MarshalFileJson(self.filename, &self.items); err != nil {
		return err
	}

	return nil
}

type JoinRound struct {
	Round           *Round
	ParticipantUser *User
}

func (self *RoundParticipantStore) Join(arg JoinRound) (*RoundParticipant, error) {
	item := self.find(func(participant RoundParticipant) bool {
		return participant.UserId == arg.ParticipantUser.Id && participant.RoundId == arg.Round.Id
	})

	if item != nil {
		return item, nil
	}

	item = &RoundParticipant{
		RoundId: arg.Round.Id,
		UserId:  arg.ParticipantUser.Id,
	}
	self.items = append(self.items, *item)

	err := self.save()

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (self *RoundParticipantStore) findAll(checkFunc func(round RoundParticipant) bool) []RoundParticipant {
	var items []RoundParticipant

	for _, it := range self.items {
		if checkFunc(it) {
			items = append(items, it)
		}
	}

	return items
}

func (self *RoundParticipantStore) FindAllByRoundId(roundId uuid.UUID) []RoundParticipant {
	return self.findAll(func(round RoundParticipant) bool {
		return round.RoundId == roundId
	})
}
