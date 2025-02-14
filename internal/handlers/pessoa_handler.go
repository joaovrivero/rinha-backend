package handlers

import (
	"net/http"
	"rinha-backend/internal/models"
	"rinha-backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/joaovrivero/rinha-backend/internal/models"
)

type PessoaHandler struct {
	repo *repositories.PessoaRepository
}

func NewPessoaHandler(repo *repositories.PessoaRepository) *PessoaHandler {
	return &PessoaHandler{repo: repo}
}

func (h *PessoaHandler) CreatePessoa(c *gin.Context) {
	var req models.CreatePessoaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pessoa := models.Pessoa{
		ID:         uuid.New().String(),
		Apelido:    req.Apelido,
		Nome:       req.Nome,
		Nascimento: req.Nascimento,
		Stack:      req.Stack,
	}

	if err := h.repo.Create(&pessoa); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "apelido já existe"})
		return
	}

	if len(pessoa.Apelido) > 32 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "apelido muito longo"})
		return
	}

	c.Header("Location", "/pessoas/"+pessoa.ID)
	c.Status(http.StatusCreated)
}
