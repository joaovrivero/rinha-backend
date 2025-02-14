package handlers

import (
	"net/http"
	"rinha-backend/internal/models"
	"rinha-backend/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PersonHandler struct {
	repo *repository.PersonRepository
}

func NewPersonHandler(repo *repository.PersonRepository) *PersonHandler {
	return &PersonHandler{repo: repo}
}

func (h *PersonHandler) Create(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validações
	if person.Nickname == "" || len(person.Nickname) > 32 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "apelido inválido"})
		return
	}
	if person.Name == "" || len(person.Name) > 100 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "nome inválido"})
		return
	}
	if _, err := time.Parse("2006-01-02", person.Birthdate); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "data de nascimento inválida"})
		return
	}
	if person.Stack != nil {
		for _, tech := range person.Stack {
			if len(tech) > 32 {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "item da stack inválido"})
				return
			}
		}
	}

	err := h.repo.Create(c.Request.Context(), &person)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "apelido já existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", "/pessoas/"+person.ID)
	c.Status(http.StatusCreated)
}

func (h *PersonHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	person, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "pessoa não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *PersonHandler) Search(c *gin.Context) {
	term := c.Query("t")
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "termo de busca obrigatório"})
		return
	}

	people, err := h.repo.Search(c.Request.Context(), term)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, people)
}

func (h *PersonHandler) Count(c *gin.Context) {
	count, err := h.repo.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "%d", count)
}
