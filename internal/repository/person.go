package repository

import (
	"context"
	"rinha-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(ctx context.Context, person *models.Person) error {
	person.ID = uuid.New().String()
	person.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(person).Error
}

func (r *PersonRepository) GetByID(ctx context.Context, id string) (*models.Person, error) {
	var person models.Person
	err := r.db.WithContext(ctx).First(&person, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *PersonRepository) Search(ctx context.Context, term string) ([]models.Person, error) {
	var people []models.Person
	err := r.db.WithContext(ctx).
		Where("nickname ILIKE ? OR name ILIKE ? OR ? = ANY(stack)",
			"%"+term+"%", "%"+term+"%", term).
		Limit(50).
		Find(&people).Error
	return people, err
}

func (r *PersonRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Person{}).Count(&count).Error
	return count, err
}
