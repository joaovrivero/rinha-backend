package server

import (
	"fmt"
	"rinha-backend/internal/config"
	"rinha-backend/internal/handlers"
	"rinha-backend/internal/models"
	"rinha-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
	db     *gorm.DB
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:    cfg,
		router: gin.Default(),
	}
}

func (s *Server) Run() error {
	// Conecta ao banco de dados
	db, err := gorm.Open(postgres.Open(s.cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	s.db = db

	// Migração
	if err := s.db.AutoMigrate(&models.Person{}); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Inicializa repositórios e handlers
	personRepo := repository.NewPersonRepository(s.db)
	personHandler := handlers.NewPersonHandler(personRepo)

	// Rotas
	s.router.POST("/pessoas", personHandler.Create)
	s.router.GET("/pessoas/:id", personHandler.GetByID)
	s.router.GET("/pessoas", personHandler.Search)
	s.router.GET("/contagem-pessoas", personHandler.Count)

	// Inicia o servidor
	return s.router.Run(":" + s.cfg.Port)
}
