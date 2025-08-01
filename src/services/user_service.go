package services

import (
	"SimpleChat/src/models"
	"errors"
)

type UserServices struct {
	users map[int]*models.User
}

func NewUserService() *UserServices {
	return &UserServices{
		users: map[int]*models.User{
			1: &models.User{ID: 1},
			2: &models.User{ID: 2},
			4: &models.User{ID: 2},
		},
	}
}

func (us *UserServices) GetUserById(id int) (*models.User, error) {
	user := us.users[id]

	if user == nil {
		return nil, errors.New("Invalid id")
	}

	return user, nil
}
