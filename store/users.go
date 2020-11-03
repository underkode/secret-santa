package store

import (
	"github.com/google/uuid"
	"underkode.ru/secret-santa/utils"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	ExternalId string    `json:"externalId"`
	Username   string    `json:"username"`
	ChatId     string    `json:"chatId"`
	Enable     bool      `json:"enable"`
}

type UserStore struct {
	filename string
	items    []User
}

func NewUserStore(filename string) (*UserStore, error) {
	var users []User

	utils.CreateIfNotExists(filename)

	err := utils.UnmarshalFileJson(filename, &users)

	if err != nil {
		return nil, err
	}

	return &UserStore{
		filename: filename,
		items:    users,
	}, nil
}

type PutUser struct {
	ExternalId string
	Username   string
	ChatId     string
}

func (self *UserStore) Put(arg PutUser) (*User, error) {
	user := self.FindByExternalId(arg.ExternalId)

	if user != nil {
		user.Username = arg.Username
		user.Enable = true
		user.ChatId = arg.ChatId
	} else {
		self.items = append(self.items, User{
			Id:         uuid.New(),
			ExternalId: arg.ExternalId,
			ChatId:     arg.ChatId,
			Username:   arg.Username,
			Enable:     true,
		})
	}

	err := self.save()

	if err != nil {
		return nil, err
	}

	return &self.items[len(self.items)-1], nil
}

func (self *UserStore) find(checkFunc func(user User) bool) *User {
	for index, it := range self.items {
		if checkFunc(it) {
			return &self.items[index]
		}
	}

	return nil
}

func (self *UserStore) FindByExternalId(externalId string) *User {
	return self.find(func(it User) bool { return it.ExternalId == externalId })
}

func (self *UserStore) save() error {
	if err := utils.MarshalFileJson(self.filename, &self.items); err != nil {
		return err
	}

	return nil
}

func (self *UserStore) FindById(id uuid.UUID) *User {
	return self.find(func(user User) bool {
		return user.Id == id
	})
}
