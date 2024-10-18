package main

import (
	"database/sql"
	"log"

	"os"
	"test/config"
	"test/run"

	"test/internal/models"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found here")
	}
}

func generateData() {
	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	gofakeit.Seed(0)
	pword := models.Password{}
	pword.Set("123456789")
	for i := 0; i <= 100; i++ {
		user := &models.User{
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: pword,
			Deleted:  false,
		}
		query := `
        INSERT INTO users (name, email, password_hash, deleted) 
        VALUES ($1, $2, $3, $4)
		RETURNING id, version
       `

		args := []any{user.Name, user.Email, user.Password.Hash, user.Deleted}

		err := db.QueryRow(query, args...).Scan(&user.ID, &user.Version)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i <= 10; i++ {
		author := &models.Author{
			Name: gofakeit.Name(),
		}
		query := `
        INSERT INTO authors (name) 
        VALUES ($1)
		RETURNING id
        `
		args := []any{author.Name}

		err := db.QueryRow(query, args...).Scan(&author.ID)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i <= 100; i++ {
		book := &models.Book{
			Title: gofakeit.BookTitle(),
			Year:  gofakeit.Date().Year(),
			Author: &models.Author{
				ID: int64(gofakeit.Number(1, 10)),
			},
		}
		query := `
	    INSERT INTO books (year, title, author_id)
	    VALUES ($1, $2, $3)
	    RETURNING id
	    `
		args := []any{book.Year, book.Title, book.Author.ID}

		err := db.QueryRow(query, args...).Scan(&book.ID)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {

	generateData()

	cfg := config.NewConfig(config.WithPort(8080), config.WithDBname("postgres"), config.WithDSN(os.Getenv("DB_DSN")))

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	app := run.NewApp(cfg, logger)

	app.Run()
	err = app.Serve()
	logger.Error(err.Error())
	os.Exit(1)

}
