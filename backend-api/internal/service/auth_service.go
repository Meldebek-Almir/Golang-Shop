package service

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"shop/backend-api/internal/models"
	"shop/backend-api/internal/repository"
)

type AuthService interface {
	IsValidUser(identifier any, password string) (models.User, error)
	CreateUser(user *models.User) error
	GetUserById(identifier any) (models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func (a *authService) IsValidUser(identifier any, password string) (models.User, error) {
	t := reflect.TypeOf(identifier)

	switch t {
	case reflect.TypeOf(""):
		user, err := a.userRepo.GetUserByIdentifier(fmt.Sprintf("%v", identifier))
		if err != nil {
			return models.User{}, err
		}
		if user.Password != password {
			return models.User{}, errors.New("invalid password")
		}
		return user, nil

	case reflect.TypeOf(0):
		num, _ := strconv.Atoi(fmt.Sprintf("%v", identifier))
		user, err := a.userRepo.GetUserByUserID(num)
		if err != nil {
			return models.User{}, err
		}
		if user.Password != password {
			return models.User{}, errors.New("invalid password")
		}
		return user, nil
	default:
		return models.User{}, errors.New("invalid type of identifier")

	}
}

func (a *authService) CreateUser(user *models.User) error {
	_, err := a.userRepo.GetUserByIdentifier(user.Email)
	if err == nil {
		return errors.New("user exists")
	}
	_, err = a.userRepo.GetUserByIdentifier(user.Nickname)
	if err == nil {
		return errors.New("user exists")
	}
	return a.userRepo.CreateUser(user)
}

func (a *authService) GetUserById(identifier any) (models.User, error) {
	var id int
	var err error
	switch v := identifier.(type) {
	case string:
		id, err = strconv.Atoi(identifier.(string))
		if err != nil {
			return models.User{}, fmt.Errorf("invalid identifier %q: %w", identifier, err)
		}
	case int:
		id = v
	default:
		return models.User{}, fmt.Errorf("unexpected identifier type %T", identifier)
	}
	user, err := a.userRepo.GetUserByUserID(id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}
	return user, nil
}
