package pkg

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/2hard4me/pkg/models"
	"gorm.io/gorm"
)

func WriteInBatches(db *gorm.DB, batchSize int) {
	ch := make(chan models.Books)
	var wg sync.WaitGroup
	var results []models.Books

	go func() {
		defer close(ch)
		result := db.FindInBatches(&results, batchSize, func(tx *gorm.DB, batch int) error {

			for _, result := range results {
				ch <- result
			}

			if tx.Error != nil {
				fmt.Println(tx.Error)
			}

			return nil
		})

		fmt.Println(result.Error)
		fmt.Printf("Rows affected: %d\n", result.RowsAffected)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for book := range ch {
			AppendToFile(book)
		}
	}()

	wg.Wait()
}

func AppendToFile(book models.Books) {
	f, err := os.OpenFile("text.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err := f.WriteString("id: " + strconv.Itoa(int(book.ID)) + " | " + "author: " + book.Author + " | " + "title: " + book.Title + " | " + "publisher: " + book.Publisher + "\n"); err != nil {
		panic(err)
	}
}
