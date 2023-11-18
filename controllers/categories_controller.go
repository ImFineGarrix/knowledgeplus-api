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

type CategoriesRepo struct {
	Db *gorm.DB
}

func NewCategoriesRepo() *CategoriesRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Category{})
	return &CategoriesRepo{Db: db}
}

// get Categoriess
func (repository *CategoriesRepo) GetCategories(c *gin.Context) {
	var Categories []models.Category
	print(Categories)
	err := models.GetCategories(repository.Db, &Categories)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Categories)
}

// get Categories by id
func (repository *CategoriesRepo) GetCategoryById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var Categories models.Category
	err := models.GetCategoryById(repository.Db, &Categories, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Categories)
}

// CreateCareer creates a new Career record.
func (repository *CategoriesRepo) CreateCategory(c *gin.Context) {
	var Categories models.Category
	// var CategoriesID = Career.Categories[0].CategoryID
	if err := c.ShouldBindJSON(&Categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(&Career)
	err := models.CreateCategory(repository.Db, &Categories)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// fmt.Println(Career.CareerID)

	c.JSON(http.StatusCreated, Categories)
}
