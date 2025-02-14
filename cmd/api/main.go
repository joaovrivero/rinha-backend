package main

import (
	"github.com/joaovrivero/rinha-backend/internal/config"
	"github.com/joaovrivero/rinha-backend/internal/routes"
)

func main() {

	db, err := config.NewDB()
	if err != nil {
		panic("failed to connect database")
	}

	router := routes.SetupRouter(db)

	router.Run(":8080") // nginx redicionará para a porta 9999
}
