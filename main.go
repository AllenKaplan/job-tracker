package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

func main() {
	fmt.Println("Starting Job Tracker")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	dbCfg := Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_NAME"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}

	db, err := NewDB(dbCfg)
	if err != nil {
		log.Fatal(err)
	}
	s := Server{db: db}

	router := gin.Default()

	router.GET("/ping", pingHandler)
	router.GET("/app/:id", s.getAppById)
	router.POST("/app", s.createApp)
	router.PUT("/app", s.updateAppById)
	router.DELETE("/app", s.deleteAppById)
	router.GET("/app", s.getApps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(fmt.Sprintf(":%s", port))
}

func NewDB(dbCfg Config) (*gorm.DB, error) {
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbCfg.Host, dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Port)
	db, err := gorm.Open(postgres.Open(dbURI))
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(Application{})
	return db, nil
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (s Server) getAppById(c *gin.Context) {
	app := Application{Uuid: c.Param("id")}

	s.db.Take(&app)

	c.JSON(http.StatusOK, app)
}

func (s Server) getApps(c *gin.Context) {
	var apps []Application
	result := s.db.Find(&apps)
	if result.Error != nil {
		log.Print(result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, apps)
}

func (s Server) createApp(c *gin.Context) {
	var app Application
	c.Bind(&app)
	app.LastUpdated = time.Now()
	app.Uuid = uuid.New().String()
	s.db.Create(&app)
	c.JSON(http.StatusOK, app)
}

func (s Server) updateAppById(c *gin.Context) {
	var app Application
	c.Bind(&app)
	app.LastUpdated = time.Now()
	s.db.Updates(&app)
	c.JSON(http.StatusOK, app)
}

func (s Server) deleteAppById(c *gin.Context) {
	type Request struct {
		Id string `json:"id,omitempty"`
	}
	var req Request
	c.Bind(&req)
	s.db.Delete(Application{Uuid: req.Id})
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("App %s deleted", req.Id)})
}
