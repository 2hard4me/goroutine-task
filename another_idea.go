package goroutinetask


// Need to make new struct with open file, append to file and close file methods.

import (
	"fmt"
	//"io/ioutil"
	"log"
	"os"

	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/2hard4me/pkg/models"
	"github.com/2hard4me/pkg/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const batchSize = 10

type FileStr struct {
	f *os.File
}

func (F *FileStr) FileOpen() error {
	file, err := os.OpenFile("text.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	F.f = file
	return err
}

func (F *FileStr) FileClose() {
	F.f.Close()
}

func (F *FileStr) FileAppend(data string) error {
	_, err := F.f.WriteString(data)
		
	return err
}

func (F *FileStr) FileClear() {
	F.f.Truncate(0)
}

func (F *FileStr) AddDataInfoToFile(books []models.Books) {
	//AppendToFile("Batch number: " + strconv.Itoa(batch) + "\n")
	for _, book := range books {
		F.FileAppend("id: " + strconv.Itoa(int(book.ID)) + " | " + "author: " + book.Author + " | " + "title: " + book.Title + " | " + "publisher: " + book.Publisher + "\n")
	}
	F.FileAppend("There are " + strconv.Itoa(int(len(books))) + " rows in this batch\n")
}


func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewConnection(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASS"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatal("could not load the database")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	// ----------------------------------------- Reading data starts here --------------------------------------------------------

	ch := make(chan []models.Books)
	defer close(ch)
	//bCh := make(chan int)
	//done := make(chan bool)

	var wg sync.WaitGroup

	fmt.Println(time.Now())
	var results []models.Books

	var Fstr FileStr
	Fstr.FileOpen()

	go func() {
		for books := range ch {
			Fstr.AddDataInfoToFile(books)
		}
	}()

	result := db.FindInBatches(&results, batchSize, func(tx *gorm.DB, batch int) error {
		wg.Add(1)
		go func() { // anonymous goroutine for writing data to the file
			defer wg.Done()

			ch <- results
			//bCh <- batch

		}()

		fmt.Println(runtime.NumGoroutine())
		//tx.Save(&results)
		//fmt.Println(tx.RowsAffected)
		//fmt.Println(batch)

		if tx.Error != nil {
			fmt.Println(tx.Error)
		}

		return nil
	})

	fmt.Println("Function FindInBatches is done")
	
	//time.Sleep(30 * time.Second)

	fmt.Printf("Number of goroutines: %d", runtime.NumGoroutine())


	

	wg.Wait()

	Fstr.FileClose()
	fmt.Println(result.Error)
	fmt.Println(result.RowsAffected)

	fmt.Println(time.Now())

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}