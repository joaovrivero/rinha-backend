package models

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid"`
	Nickname  string         `json:"apelido" gorm:"uniqueIndex;not null;size:32"`
	Name      string         `json:"nome" gorm:"not null;size:100"`
	Birthdate string         `json:"nascimento" gorm:"not null"`
	Stack     []string       `json:"stack" gorm:"type:text[]"`
	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
