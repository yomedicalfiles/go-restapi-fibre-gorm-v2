package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abiiranathan/gofibre-tuts/books"
	"github.com/abiiranathan/gofibre-tuts/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupRoutes(app *fiber.App) {
	bookRoutes := app.Group("/api/v1/books")

	bookRoutes.Get("/", books.GetBooks)
	bookRoutes.Post("/", books.CreateBook)
	bookRoutes.Get("/:id", books.GetBook)
	bookRoutes.Delete("/:id", books.DeleteBook)
	bookRoutes.Put("/:id", books.UpdateBook)
}

func init() {
	var err error
	godotenv.Load(".env")

	DSN := os.Getenv("DSN")
	DRIVER_NAME := os.Getenv("DB_DRIVER")

	if DSN == "" || DRIVER_NAME == "" {
		panic("Specify the driver name and data source name[DSN]")
	}

	database.DBConn, err = gorm.Open(mysql.New(mysql.Config{
		DriverName: os.Getenv("DRIVER_NAME"),
		DSN:        DSN,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		Logger:                                   logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		panic("Unable to connect to the database" + err.Error())
	}

	log.Printf("Connected to database...")
	database.DBConn.AutoMigrate(&books.Book{})
	log.Printf("Database migrated...")

}

func main() {
	app := fiber.New()

	setupRoutes(app)
	fmt.Println("Server started on port 3000")
	log.Fatal(app.Listen(":3000"))
}

/*
	sqlite3: database.DBConn, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
*/
