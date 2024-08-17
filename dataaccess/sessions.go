package dataaccess

import (
	"time"

	"github.com/nicknad/krankentransport/db"
)

type Session struct {
	Token      string
	Login      string
	ExpiryDate time.Time
}

func GetSessions() ([]Session, error) {
	db := db.GetDB()
	results, err := db.Query("SELECT token, login FROM sessions")

	if err != nil {
		return nil, err
	}
	defer results.Close()

	var sessions []Session
	for results.Next() {
		var session Session

		if err := results.Scan(&session.Token, &session.Login); err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func GetSession(session string) (Session, error) {
	db := db.GetDB()
	stmt, err := db.Prepare("SELECT token, login FROM sessions where token = ?")
	if err != nil {
		return Session{}, err
	}
	defer stmt.Close()
	var foundSession Session
	err = stmt.QueryRow(session).Scan(&foundSession.Token, &foundSession.Login)
	if err != nil {
		return Session{}, err
	}

	return foundSession, nil
}

func CreateSession(token string, login string) (Session, error) {
	db := db.GetDB()
	stmt, err := db.Prepare("INSERT INTO sessions (token, login) VALUES (?, ?) RETURNING token, login")

	if err != nil {
		return Session{}, err
	}

	var createdSession Session
	err = stmt.QueryRow(token, login).Scan(&createdSession.Token, &createdSession.Login)
	if err != nil {
		return Session{}, err
	}

	return createdSession, nil
}

func ClearSession() error {
	db := db.GetDB()
	stmt, err := db.Prepare("DELETE FROM sessions")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
