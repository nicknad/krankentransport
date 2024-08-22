package main

import (
	"github.com/nicknad/krankentransport/dataaccess"
	"github.com/nicknad/krankentransport/db"
	"github.com/nicknad/krankentransport/migrations"
)

func main() {
	db.DATABASE_URL = "./krankentransport.db"
	migrations.RunMigrations()

	dataaccess.CreateUser("user", "password", true)
}
