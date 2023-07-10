package repository

import (
	"github.com/2hard4me/pkg/models"
	"gorm.io/gorm"
)

type Book interface {
	GetBatch(id int) ([]models.Books, error)
}

type Repository struct {
	Book
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Book: NewBookPostgres(db),
	}
}