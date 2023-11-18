package controllers

import (
	"errors"
	"knowledgeplus/go-api/database"
	"knowledgeplus/go-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CareerRepo struct {
	Db *gorm.DB
}

func NewCareerRepo() *CareerRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Career{})
	return &CareerRepo{Db: db}
}

// get Careers
func (repository *CareerRepo) GetCareers(c *gin.Context) {
	var Career []models.Career
	print(Career)
	err := models.GetCareers(repository.Db, &Career)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Career)
}

// get Career by id
func (repository *CareerRepo) GetCareer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var Career models.Career
	err := models.GetCareerById(repository.Db, &Career, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Career)
}

// CreateCareer creates a new Career record.
func (repository *CareerRepo) CreateCareer(c *gin.Context) {
	var Career models.Career
	// var CategoriesID = Career.Categories[0].CategoryID
	if err := c.ShouldBindJSON(&Career); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(&Career)
	err := models.CreateCareer(repository.Db, &Career)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// fmt.Println(Career.CareerID)

	c.JSON(http.StatusCreated, Career)
}
