package migrations

import "github.com/nicknad/krankentransport/db"

func Up(key string, upFun func()) {
	db := db.GetDB()
	stmt, err := db.Prepare("SELECT key FROM migrations WHERE key = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var foundKey string
	err = stmt.QueryRow(key).Scan(&foundKey)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			upFun()
			stmt, err = db.Prepare("INSERT INTO migrations (key) VALUES (?)")
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(key)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

}

func RunMigrations() {
	db := db.GetDB()
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS migrations (key TEXT PRIMARY KEY, createdAt INTEGER DEFAULT (strftime('%s', 'now')));")
	if err != nil {
		panic(err)
	}
	stmt.Exec()

	Up("create-tables", func() {
		stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, login TEXT UNIQUE, passwordhash TEXT, admin BOOLEAN DEFAULT 0,  createdAt INTEGER DEFAULT (strftime('%s', 'now')));")
		if err != nil {
			panic(err)
		}
		stmt.Exec()

		// add a unique constraint to the email column if it does not exists
		stmt, err = db.Prepare("CREATE UNIQUE INDEX IF NOT EXISTS user_unique ON users (login);")
		if err != nil {
			panic(err)
		}
		stmt.Exec()

		stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS krankenfahrten (id INTEGER PRIMARY KEY, description TEXT, createdAt INTEGER DEFAULT (strftime('%s', 'now')), acceptedByLogin TEXT, acceptedAt INTEGER, finished BOOLEAN DEFAULT 0);")
		if err != nil {
			panic(err)
		}
		stmt.Exec()
		stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS sessions (token TEXT PRIMARY KEY, login TEXT, expireAt INTEGER);")
		if err != nil {
			panic(err)
		}
		stmt.Exec()
	})
}
