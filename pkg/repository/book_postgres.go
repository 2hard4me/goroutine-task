package repository

import (
	"sync"
	"github.com/2hard4me/pkg/models"
	"gorm.io/gorm"
)

const batchSize = 100

type BookPostgres struct {
	db *gorm.DB
}

func NewBookPostgres(db *gorm.DB) *BookPostgres {
	return &BookPostgres{
		db: db,
	}
}

func (r *BookPostgres) GetBatch(id int) ([]models.Books, error) {
	var batch []models.Books
	ch := make(chan models.Books)
	var wg sync.WaitGroup
	var results []models.Books

	go func() {
		defer close(ch)
		result := r.db.FindInBatches(&results, batchSize, func(tx *gorm.DB, batch int) error {
			if batch == id {
				for _, result := range results {
					ch <- result
				}
			}

			if tx.Error != nil {
				return tx.Error
			}

			return nil
		})
		if result.Error != nil {
			panic(result.Error)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for book := range ch {
			batch = append(batch, book)
		}
	}()

	wg.Wait()

	return batch, nil
}

