package main

import (
	"fmt"

	"github.com/nicknad/krankentransport/dataaccess"
	"github.com/nicknad/krankentransport/db"
	"github.com/nicknad/krankentransport/migrations"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	db.DATABASE_URL = viper.GetString("database.url")
	fmt.Println(db.DATABASE_URL)
	migrations.RunMigrations()

	dataaccess.CreateUser("user", "password", true)
}
