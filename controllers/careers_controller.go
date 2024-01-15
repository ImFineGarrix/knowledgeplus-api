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

// get Careers with pagination
func (repository *CareerRepo) GetCareers(c *gin.Context) {
	var (
		careers    []models.Career
		pagination models.Pagination
		err        error
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	careers, pagination, err = models.GetCareers(repository.Db, page, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"careers":    careers,
		"pagination": pagination,
	})
}

// get Careers with categories and pagination
func (repository *CareerRepo) GetCareersWithHaveCategories(c *gin.Context) {
	var (
		careerIDs  []uint
		careers    []models.Career
		pagination models.Pagination
		err        error
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	// Query distinct career IDs
	err = repository.Db.
		Model(&models.Career{}).
		Select("DISTINCT career_id").
		Limit(limit).
		Offset((page-1)*limit).
		Pluck("career_id", &careerIDs).
		Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Query complete careers based on distinct IDs
	err = repository.Db.
		Preload("Categories").
		Preload("Skills").
		Where("career_id IN ?", careerIDs).
		Find(&careers).
		Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Count total careers (ignoring pagination)
	var totalCareers int64
	err = repository.Db.Model(&models.Career{}).Distinct("career_id").Count(&totalCareers).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert totalCareers to int
	totalCareersInt := int(totalCareers)

	// Create pagination information
	pagination = models.Pagination{
		Page:  page,
		Total: totalCareersInt,
		Limit: limit,
	}

	c.JSON(http.StatusOK, gin.H{
		"careers":    careers,
		"pagination": pagination,
	})
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

	// fmt.Println(Career.Categories[0].CategoryID)
	// fmt.Println(&Career)
	err := models.CreateCareer(repository.Db, &Career)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// fmt.Println(Career.CareerID)

	c.JSON(http.StatusCreated, Career)
}

// UpdateCareer updates a Career record by ID.
func (repository *CareerRepo) UpdateCareer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedCareer models.Career

	// Check if the career record exists
	err := models.GetCareerById(repository.Db, &updatedCareer, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Bind the updated data from the request
	if err := c.ShouldBindJSON(&updatedCareer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the career record
	err = models.UpdateCareer(repository.Db, &updatedCareer)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedCareer)
}

// DeleteCareer deletes a Career record by ID.
func (repository *CareerRepo) DeleteCareer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var career models.Career

	// Find the career record
	err := repository.Db.Where("career_id = ?", id).First(&career).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete associated records in categories_careers table
	err = repository.Db.Exec("DELETE FROM categories_careers WHERE career_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete associated records in careers_skills table
	err = repository.Db.Exec("DELETE FROM careers_skills WHERE career_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete the career record
	err = repository.Db.Delete(&career).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Career and associated records deleted successfully"})
}
