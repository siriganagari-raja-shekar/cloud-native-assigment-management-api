package main

import (
	"csye6225-mainproject/conf"
	"csye6225-mainproject/db"
	"csye6225-mainproject/routes"
	"fmt"
)

func main() {

	config := &conf.Configuration{}
	config.Set()

	router := routes.SetupRouter(&db.Store{})

	err := router.Run(":8000")
	if err != nil {
		fmt.Printf("Fatal server error")
	}
}
