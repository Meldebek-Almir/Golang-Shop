package repository

import (
	"database/sql"
	"time"

	"shop/backend-api/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByIdentifier(identifier string) (models.User, error)
	GetUserByUserID(id int) (models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users(nickname, age, gender, first_name, last_name, email, password, role, created, updated) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := u.db.Exec(query, user.Nickname, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, user.Password, user.Role, time.Now(), time.Now())
	return err
}

func (u *userRepository) GetUserByIdentifier(identifier string) (models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1 OR nickname = $2`
	err := u.db.QueryRow(query, identifier, identifier).Scan(&user.ID, &user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.Created, &user.Updated)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userRepository) GetUserByUserID(id int) (models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE user_id = $1`
	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.Created, &user.Updated)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
