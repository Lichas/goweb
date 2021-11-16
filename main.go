package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("db open err.")
	}
	return db
}

type BlogTag struct {
	gorm.Model
	Name string
}

type Blog struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	TagID   uint
}
type Auth struct {
	gorm.Model
	Username string
	Password string
}

var db *gorm.DB

func initData(db *gorm.DB) {
	db.AutoMigrate(&BlogTag{})
	db.AutoMigrate(&Blog{})
	db.AutoMigrate(&Auth{})

	db.Create(&BlogTag{Name: "mytag"})
	db.Create(&Blog{Title: "First blog", Content: "my awesome hello text.", TagID: 1})
	db.Create(&Auth{Username: "Tom", Password: "123456**"})

	var t1 BlogTag
	db.First(&t1, 1)
	fmt.Println(t1.Name)
	var b1 Blog
	db.First(&b1)
	fmt.Println(b1.Content)
	var a1 Auth
	db.First(&a1)
	fmt.Println(a1.Password)
}

func main() {
	db = NewDB()
	initData(db)

	r := gin.Default()
	r.GET("/ping", Ping)     //when c *gin.Context pass in
	r.GET("/blog", GetBlogs) // curl -i -X POST 0.0.0.0:8080 -d '{"title":"testTitle","content":"good txt."}'
	r.POST("/addblog", AddBlog)
	r.Run()
}

// Ping ping detect
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"pingMessage": "Ping Success.",
	})
}

// GetBlogs get blogs from db return json
func GetBlogs(c *gin.Context) {
	var blogs []Blog
	if err := db.Find(&blogs).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, blogs)
	}
}

func AddBlog(c *gin.Context) {
	var blog Blog
	c.BindJSON(&blog)
	db.Create(&blog)
	c.JSON(200, blog)
}
