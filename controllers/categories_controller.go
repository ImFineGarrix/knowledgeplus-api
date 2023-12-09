package controllers

import (
	"errors"
	"fmt"
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

// CreateCategory creates a new Category record.
func (repository *CategoriesRepo) CreateCategory(c *gin.Context) {
	var Categories models.Category
	if err := c.ShouldBindJSON(&Categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := models.CreateCategory(repository.Db, &Categories)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, Categories)
}

// UpdateCategory updates a Category record by ID.
func (repository *CategoriesRepo) UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var existingCategory models.Category

	err := repository.Db.First(&existingCategory, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var updatedCategory models.UpdateCategoryModels
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		fmt.Println(updatedCategory.Name)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the fields you want to change
	if len(updatedCategory.Name) != 0 {
		existingCategory.Name = updatedCategory.Name
	}
	existingCategory.ImageUrl = updatedCategory.ImageUrl

	err = repository.Db.Save(&existingCategory).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, existingCategory)
}

// DeleteCategoryById deletes a Category record by ID.
func (repository *CategoriesRepo) DeleteCategoryById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var category models.Category

	err := repository.Db.Where("category_id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete associated records in categories_careers table
	err = repository.Db.Exec("DELETE FROM categories_careers WHERE category_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete the category record
	err = repository.Db.Delete(&category).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category and associated records deleted successfully"})
}
