package main

import (
	"github.com/joaovrivero/rinha-backend/internal/config"
)

func main() {

	db, err := config.NewDB()
	if err != nil {
		panic("failed to connect database")
	}

	router := config.SetupRouter(db)

	r.Run(":8080") // nginx redicionará para a porta 9999
}
