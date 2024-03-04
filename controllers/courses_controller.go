package controllers

import (
	"errors"
	"knowledgeplus/go-api/database"
	"knowledgeplus/go-api/models"
	"knowledgeplus/go-api/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CourseRepo struct {
	Db *gorm.DB
}

func NewCourseRepo() *CourseRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Course{}, &models.Organizations{})
	return &CourseRepo{Db: db}
}

// GetCourses retrieves all Course records from the database.
func (repository *CourseRepo) GetCourses(c *gin.Context) {
	var (
		courses    []models.Course
		pagination models.Pagination
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	pagination, err = models.GetCourses(repository.Db, page, limit, &courses)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses":    courses,
		"pagination": pagination,
	})
}

// GetAllSkillsWithFilter retrieves all skill records from the database with filtering and pagination.
func (repository *CourseRepo) GetAllSkillsWithFilter(c *gin.Context) {
	var (
		courses    []models.Course
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

	courses, pagination, err = models.GetAllCoursesWithFilter(repository.Db, page, limit, search)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses":    courses,
		"pagination": pagination,
	})
}

// GetCourseById retrieves a Course by its ID from the database.
func (repository *CourseRepo) GetCourseById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var course models.Course
	err := models.GetCourseById(repository.Db, &course, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, course)
}

// GetCoursesBySkillId retrieves courses based on the provided SkillID with pagination.
func (repository *CourseRepo) GetCoursesBySkillId(c *gin.Context) {
	skillID, err := strconv.Atoi(c.Param("skill_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SkillID"})
		return
	}

	var (
		courses    []models.Course
		pagination models.Pagination
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	pagination, err = models.GetCoursesBySkillId(repository.Db, skillID, page, limit, &courses)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses":    courses,
		"pagination": pagination,
	})
}

// GetCoursesByCareerId retrieves courses based on the provided CareerID with pagination.
func (repository *CourseRepo) GetCoursesByCareerId(c *gin.Context) {
	careerID, err := strconv.Atoi(c.Param("career_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CareerID"})
		return
	}

	var (
		courses    []models.CourseWithoutSkillLevels
		pagination models.Pagination
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	pagination, err = models.GetCoursesByCareerId(repository.Db, careerID, page, limit, &courses)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses":    courses,
		"pagination": pagination,
	})
}

// CreateCourse creates a new Course record.
func (repository *CourseRepo) CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
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
			c.JSON(http.StatusCreated, out)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the name already exists in the database
	var existingCourse models.Course
	if err := repository.Db.Where("name = ?", course.Name).First(&existingCourse).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	err := models.CreateCourse(repository.Db, &course)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// UpdateCourse updates a Course record by ID.
func (repository *CourseRepo) UpdateCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedCourse models.Course

	// Check if the course record exists
	err := models.GetCourseById(repository.Db, &updatedCourse, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Bind the updated data from the request
	if err := c.ShouldBindJSON(&updatedCourse); err != nil {
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
			c.JSON(http.StatusCreated, out)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the name already exists in the database
	var existingCourse models.Course
	if err := repository.Db.Where("name = ? AND course_id != ?", updatedCourse.Name, updatedCourse.CourseID).First(&existingCourse).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// Update the course record
	err = models.UpdateCourse(repository.Db, &updatedCourse)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedCourse)
}

// DeleteCourseById deletes a Course record by ID.
func (repository *CourseRepo) DeleteCourseById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var course models.Course

	err := repository.Db.Where("course_id = ?", id).
		Preload("Organization").
		Preload("SkillsLevels").
		First(&course).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete associated records in skills_levels table
	err = repository.Db.Exec("DELETE FROM skills_levels WHERE course_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete the course record
	err = repository.Db.Delete(&course).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course and associated records deleted successfully"})
}
