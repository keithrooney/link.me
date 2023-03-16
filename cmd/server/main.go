package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type postgresDataSource struct {
	Username string
	Password string
	Host     string
	Port     string
}

func main() {
	// database := internal.GetDatabase()
	// database.Exec("SELECT * FROM users")
	viper.AddConfigPath("/home/krooney/workspace/anchorly/cmd/server")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	dataSource := &postgresDataSource{}
	if err := viper.UnmarshalKey("datasource", dataSource); err != nil {
		panic(err)
	}
	fmt.Println(dataSource)
	// database = NewDatabase(dataSource)
}
