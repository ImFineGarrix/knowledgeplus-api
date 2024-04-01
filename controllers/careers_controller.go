package controllers

import (
	"errors"
	"knowledgeplus/go-api/models"
	"knowledgeplus/go-api/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CareerRepo struct {
	Db      *gorm.DB
	UserDb  *gorm.DB
	AdminDb *gorm.DB
}

func NewCareerRepo(db *gorm.DB, userDb *gorm.DB, admiinDb *gorm.DB) *CareerRepo {
	db.AutoMigrate(&models.Career{})
	userDb.AutoMigrate(&models.Career{})
	admiinDb.AutoMigrate(&models.Career{})
	return &CareerRepo{Db: db, UserDb: userDb, AdminDb: admiinDb}
}

// get all Careers use with backoffice
// get Careers with pagination
func (repository *CareerRepo) GetCareers(c *gin.Context) {
	log.Default().Print(c.MustGet("userRole"))
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

// GetAllCareers retrieves all Career records from the database with filtering and pagination. user with frontend
func (repository *CareerRepo) GetAllCareersWithFilters(c *gin.Context) {
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

	search := c.Query("search")
	groupID, _ := strconv.ParseInt(c.Query("group"), 10, 64)

	careers, pagination, err = models.GetCareersWithFilters(repository.Db, page, limit, search, groupID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"careers":    careers,
		"pagination": pagination,
	})
}

// GetCareersByCourseId retrieves careers based on the provided CourseID with pagination.
func (repository *CareerRepo) GetCareersByCourseId(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CourseID"})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	careers, pagination, err := models.GetCareersByCourseId(repository.Db, courseID, page, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"careers":    careers,
		"pagination": pagination,
	})
}

// GetCareersBySkillId retrieves careers based on the provided SkillID with pagination.
func (repository *CareerRepo) GetCareersBySkillId(c *gin.Context) {
	skillID, err := strconv.Atoi(c.Param("skill_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SkillID"})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	careers, pagination, err := models.GetCareersBySkillId(repository.Db, skillID, page, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	// Bind the JSON request to the Career struct
	if err := c.ShouldBindJSON(&Career); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorMsg{
					Code:    http.StatusBadRequest,
					Field:   fe.Field(),
					Message: response.GetErrorMsg(fe),
				}
			}
			c.JSON(http.StatusBadRequest, out)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the name already exists in the database
	var existingCareer models.Career
	if err := repository.Db.Where("name = ?", Career.Name).First(&existingCareer).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// Create the Career record
	err := models.CreateCareer(repository.Db, &Career)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, Career)
}

// UpdateCareer updates a Career record by ID.
func (repository *CareerRepo) UpdateCareer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedCareer models.Career

	err := repository.Db.First(&updatedCareer, id).Error
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
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorMsg{
					Code:    http.StatusBadRequest,
					Field:   fe.Field(),
					Message: response.GetErrorMsg(fe),
				}
			}
			c.JSON(http.StatusBadRequest, out)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the name already exists in the database
	var existingCareer models.Career
	if err := repository.Db.Where("name = ? AND career_id != ?", updatedCareer.Name, id).First(&existingCareer).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
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

	// Delete associated records in groups_careers table
	err = repository.Db.Exec("DELETE FROM groups_careers WHERE career_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete associated records in careers_skills table
	err = repository.Db.Exec("DELETE FROM careers_skills_levels WHERE career_id = ?", id).Error
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

func (repository *CareerRepo) RecommendSkillsLevelsByCareer(c *gin.Context) {
	var Career models.CareerForRecommendSkillsLevels
	var careerBody models.Career

	// Bind the JSON request to the Career struct
	if err := c.ShouldBindJSON(&Career); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorMsg{
					Code:    http.StatusBadRequest,
					Field:   fe.Field(),
					Message: response.GetErrorMsg(fe),
				}
			}
			c.JSON(http.StatusBadRequest, out)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Find the career record
	err := repository.Db.Where("career_id = ?", Career.CurrentCareerID).First(&careerBody).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Call the function to get the differences
	result, rErr := models.RecommendSkillsLevelsByCareer(repository.Db, &Career)
	if rErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return differences in the response
	c.JSON(http.StatusOK, result)
}
