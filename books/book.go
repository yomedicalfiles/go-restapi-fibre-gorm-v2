package books

import (
	"time"

	"github.com/abiiranathan/gofibre-tuts/database"
	"github.com/gofiber/fiber/v2"
)

type Book struct {
	ID        uint      `json:"id" gorm:"primarykey;autoIncrement"`
	Title     string    `json:"title" gorm:"not null; size:100"`
	Author    string    `json:"author" gorm:"not null; size:100"`
	Rating    int       `json:"rating" gorm:"not null; check:rating > 0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func GetBooks(c *fiber.Ctx) error {
	db := database.DBConn

	var books []Book
	db.Find(&books)

	return c.JSON(books)

}

func CreateBook(c *fiber.Ctx) error {
	db := database.DBConn

	book := new(Book)
	errorResp := new(ErrorResponse)

	if err := c.BodyParser(book); err != nil {
		c.Status(400).Send([]byte(err.Error()))
	}

	if book.Title == "" {
		errorResp.Error = "Book title is required!"
		c.Status(400).JSON(errorResp)
		return nil
	}

	if book.Author == "" {
		errorResp.Error = "Book author is required!"
		c.Status(400).JSON(errorResp)
		return nil
	}

	if book.Rating == 0 {
		errorResp.Error = "Book rating should be greater than 0 !"
		c.Status(400).JSON(errorResp)
		return nil
	}

	// Make sure no same book with title & author

	if result := db.Limit(1).Where("title = ? AND author = ?", book.Title, book.Author).Find(&Book{}); result.RowsAffected > 0 {
		errorResp.Error = "Book with this title and author already exists"
		c.Status(400).JSON(errorResp)
		return nil
	} else {
		db.Create(&book)
		return c.JSON(book)

	}

}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DBConn
	var book Book

	db.Find(&book, id)
	return c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DBConn
	var book Book

	db.First(&book, id)

	if book.ID == 0 {
		c.Status(404).JSON(nil)
		return nil
	}

	db.Delete(&book)

	return c.Status(204).JSON(nil)
}

func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DBConn
	var book Book

	db.First(&book, id)

	if book.ID == 0 {
		c.Status(404).JSON(nil)
		return nil
	}

	var newBook = new(Book)

	// Parse request body
	if err := c.BodyParser(newBook); err != nil {
		c.Status(400).Send([]byte(err.Error()))
	}

	if newBook.Title != "" {
		book.Title = newBook.Title
	}

	if newBook.Author != "" {
		book.Author = newBook.Author
	}

	if newBook.Rating > 0 {
		book.Rating = newBook.Rating
	}

	// db.Select("title", "author", "rating", "updated_at").Updates(&book)
	db.Updates(&book)

	return c.Status(200).JSON(book)
}
