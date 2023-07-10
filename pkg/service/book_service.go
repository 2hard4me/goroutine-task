package service

import (
	"github.com/2hard4me/pkg/models"
	"github.com/2hard4me/pkg/repository"
)

type BookService struct {
	repo repository.Book
}

func NewBookService(repo repository.Book) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s *BookService) GetBatch(id int) ([]models.Books, error) {
	return s.repo.GetBatch(id)
}