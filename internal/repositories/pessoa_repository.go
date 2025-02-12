package repositories

import (
	"rinha-backend/internal/models"

	"gorm.io/gorm"
)

type PessoaRepository struct {
	db *gorm.DB
}

func NewPessoaRepository(db *gorm.DB) *PessoaRepository {
	return &PessoaRepository{db: db}
}

func (r *PessoaRepository) Create(pessoa *models.Pessoa) error {
	return r.db.Create(pessoa).Error
}

func (r *PessoaRepository) FindById(id uint) (*models.Pessoa, error) {
	var pessoa models.Pessoa
	result := r.db.Where("id = ?", id).First(&pessoa)
	return &pessoa, result.Error
}
