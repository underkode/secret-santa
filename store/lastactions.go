package store

import (
	"github.com/google/uuid"
	"underkode.ru/secret-santa/utils"
)

type LastAction struct {
	UserId uuid.UUID `json:"userId"`
	Action string    `json:"action"`
}

type LastActionStore struct {
	filename string
	items    []LastAction
}

func NewLastActionStore(filename string) (*LastActionStore, error) {
	var items []LastAction

	utils.CreateIfNotExists(filename)

	err := utils.UnmarshalFileJson(filename, &items)

	if err != nil {
		return nil, err
	}

	return &LastActionStore{
		filename: filename,
		items:    items,
	}, nil
}

type PutLastAction struct {
	Action    string
	OwnerUser *User
}

func (self *LastActionStore) Put(arg PutLastAction) (*LastAction, error) {
	lastAction := self.FindByUserId(arg.OwnerUser.Id)

	if lastAction != nil {
		lastAction.Action = arg.Action
	} else {
		self.items = append(self.items, LastAction{
			Action: arg.Action,
			UserId: arg.OwnerUser.Id,
		})
	}

	err := self.save()

	if err != nil {
		return nil, err
	}

	return &self.items[len(self.items)-1], nil
}

func (self *LastActionStore) find(checkFunc func(user LastAction) bool) *LastAction {
	for index, it := range self.items {
		if checkFunc(it) {
			return &self.items[index]
		}
	}

	return nil
}

func (self *LastActionStore) FindByUserId(externalId uuid.UUID) *LastAction {
	return self.find(func(it LastAction) bool { return it.UserId == externalId })
}

func (self *LastActionStore) save() error {
	if err := utils.MarshalFileJson(self.filename, &self.items); err != nil {
		return err
	}

	return nil
}
