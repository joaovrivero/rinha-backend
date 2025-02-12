package routes

import (
	"rinha-backend/internal/handlers"
	"rinha-backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	pessoaRepository := repositories.NewPessoaRepository(db)
	pessoaHandler := handlers.NewPessoaHandler(pessoaRepository)

	// Rotas
	r.POST("/pessoas", pessoaHandler.CriarPessoa)
	r.GET("/pessoas/:id", pessoaHandler.BuscarPessoaPorId)
	r.GET("/pessoas", pessoaHandler.BuscarPessoas)
	r.GET("/contagem-pessoas", pessoaHandler.ContagemPessoas)

	return r
}
