package service

import (
	"github.com/2hard4me/pkg/models"
	"github.com/2hard4me/pkg/repository"
)

type Book interface {
	GetBatch(id int) ([]models.Books, error)
}

type Service struct {
	Book
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Book: NewBookService(repos.Book),
	}
}