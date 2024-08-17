package dataaccess

import (
	"database/sql"
	"time"

	"github.com/nicknad/krankentransport/db"
)

type Krankenfahrt struct {
	Id              int
	Description     string
	CreatedAt       time.Time
	AcceptedByLogin *string
	AcceptedAt      *time.Time
	Finished        bool
}

func GetKrankenfahrt(id int) (Krankenfahrt, error) {
	db := db.GetDB()
	stmt, err := db.Prepare("SELECT id, description, createdAt, acceptedByLogin, acceptedAt, finished FROM krankenfahrten WHERE id = ?")
	if err != nil {
		return Krankenfahrt{}, err
	}
	defer stmt.Close()

	var createdAt int64
	var acceptedAt sql.NullInt64 // Handles NULL values
	var acceptedByLogin sql.NullString
	var foundFahrt Krankenfahrt
	err = stmt.QueryRow(id).Scan(&foundFahrt.Id, &foundFahrt.Description, &createdAt, &acceptedByLogin, &acceptedAt, &foundFahrt.Finished)
	if err != nil {
		return Krankenfahrt{}, err
	}

	foundFahrt.CreatedAt = time.Unix(createdAt, 0)
	if acceptedAt.Valid {
		acceptedAtTime := time.Unix(acceptedAt.Int64, 0)
		foundFahrt.AcceptedAt = &acceptedAtTime
	} else {
		foundFahrt.AcceptedAt = nil
	}

	if acceptedByLogin.Valid {
		foundFahrt.AcceptedByLogin = &acceptedByLogin.String
	} else {
		foundFahrt.AcceptedByLogin = nil
	}

	return foundFahrt, nil
}

func UndoAcceptKrankenfahrt(fahrt Krankenfahrt) error {
	db := db.GetDB()

	stmt, err := db.Prepare("Update krankenfahrten SET acceptedAt = NULL, acceptedByLogin = NULL WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(fahrt.Id)

	if err != nil {
		return err
	}

	return nil
}

func UpdateKrankenfahrt(fahrt Krankenfahrt) error {
	db := db.GetDB()
	stmt, err := db.Prepare("Update krankenfahrten SET acceptedAt = ?, acceptedByLogin = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(fahrt.AcceptedAt.Unix(), fahrt.AcceptedByLogin, fahrt.Id)

	if err != nil {
		return err
	}

	return nil
}

func GetKrankenfahrten() ([]Krankenfahrt, error) {
	db := db.GetDB()
	results, err := db.Query("SELECT id, description, createdAt, acceptedByLogin, acceptedAt, finished FROM krankenfahrten")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var fahrten []Krankenfahrt
	for results.Next() {
		var fahrt Krankenfahrt
		var createdAt int64
		var acceptedAt sql.NullInt64 // Handles NULL values
		var acceptedByLogin sql.NullString

		if err := results.Scan(&fahrt.Id, &fahrt.Description, &createdAt, &acceptedByLogin, &acceptedAt, &fahrt.Finished); err != nil {
			return nil, err
		}

		fahrt.CreatedAt = time.Unix(createdAt, 0)
		if acceptedAt.Valid {
			acceptedAtTime := time.Unix(acceptedAt.Int64, 0)
			fahrt.AcceptedAt = &acceptedAtTime
		} else {
			fahrt.AcceptedAt = nil
		}

		if acceptedByLogin.Valid {
			fahrt.AcceptedByLogin = &acceptedByLogin.String
		} else {
			fahrt.AcceptedByLogin = nil
		}

		fahrten = append(fahrten, fahrt)
	}

	return fahrten, nil
}

func DeleteKrankenfahrt(id int) error {
	db := db.GetDB()
	stmt, err := db.Prepare("DELETE FROM krankenfahrten WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func CreateKrankenfahrt(desc string) (Krankenfahrt, error) {
	db := db.GetDB()
	stmt, err := db.Prepare("INSERT INTO krankenfahrten (description, createdAt, finished) VALUES (?, ?, ?) RETURNING id, description, finished")
	if err != nil {
		return Krankenfahrt{}, err
	}
	defer stmt.Close()
	unixTime := time.Now().Unix()
	var createdFahrt Krankenfahrt
	err = stmt.QueryRow(desc, unixTime, false).Scan(&createdFahrt.Id, &createdFahrt.Description, &createdFahrt.Finished)
	if err != nil {
		return Krankenfahrt{}, err
	}

	return createdFahrt, nil
}
