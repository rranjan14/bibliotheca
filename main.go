package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	ID     uint
	Title  string
	Author string
}

func connectDatabase(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("database", db)
	}
}

func setupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Book{},
	)
	if err != nil {
		return fmt.Errorf("Error migrating database: %s", err)
	}
	return nil
}

func defaultHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "default.html", gin.H{})
}

func setupRouter(r *gin.Engine, db *gorm.DB) {
	r.LoadHTMLGlob("templates/**/*.html")
	r.Use(connectDatabase(db))
	r.GET("/books", bookIndexHandler)
	r.GET("/", defaultHandler, func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/books")
	})
}

func bookIndexHandler(c *gin.Context) {
	db := c.Value("database").(*gorm.DB)
	books := []Book{}
	if err := db.Find(&books).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "books/index.html", gin.H{"books": books})
}

func main() {
	db, err := gorm.Open(sqlite.Open("bibliotheca.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	err = setupDatabase(db)

	if err != nil {
		log.Fatalf("Database setup error: %s", err)
	}

	r := gin.Default()

	setupRouter(r, db)

	err = r.Run(":8000")

	if err != nil {
		log.Fatalf("gin Run error: %s", err)
	}
}
