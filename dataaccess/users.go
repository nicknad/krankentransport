package dataaccess

import (
	"github.com/nicknad/krankentransport/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uint
	Login        string
	PasswordHash string
	Admin        bool
}

func GetUsers() ([]User, error) {
	db := db.GetDB()
	results, err := db.Query("SELECT id, login, admin FROM users")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var users []User
	for results.Next() {
		var user User

		if err := results.Scan(&user.Id, &user.Login, &user.Admin); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUser(login string) (User, error) {
	db := db.GetDB()
	stmt, err := db.Prepare("SELECT id, login, passwordhash, admin FROM users WHERE login = ?")
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	var foundUser User
	err = stmt.QueryRow(login).Scan(&foundUser.Id, &foundUser.Login, &foundUser.PasswordHash, &foundUser.Admin)
	if err != nil {
		return User{}, err
	}

	return foundUser, nil
}

func GetUserById(id int) (User, error) {
	db := db.GetDB()
	stmt, err := db.Prepare("SELECT id, login, passwordhash, admin FROM users WHERE id = ?")
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	var foundUser User
	err = stmt.QueryRow(id).Scan(&foundUser.Id, &foundUser.Login, &foundUser.PasswordHash, &foundUser.Admin)
	if err != nil {
		return User{}, err
	}

	return foundUser, nil
}

func CreateUser(login string, password string, admin bool) (User, error) {
	db := db.GetDB()
	stmt, err := db.Prepare("INSERT INTO users (login, passwordhash, admin) VALUES (?, ?, ?) RETURNING id, login, admin")
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return User{}, err
	}

	var createdUser User
	err = stmt.QueryRow(login, string(hash), admin).Scan(&createdUser.Id, &createdUser.Login, &createdUser.Admin)
	if err != nil {
		return User{}, err
	}

	return createdUser, nil
}

func DeleteUser(login string) error {
	db := db.GetDB()
	stmt, err := db.Prepare("DELETE FROM users WHERE login = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(login)
	if err != nil {
		return err
	}

	return nil
}
