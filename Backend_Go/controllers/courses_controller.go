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

type CourseRepo struct {
	Db *gorm.DB
}

func NewCourseRepo() *CourseRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Courses{})
	return &CourseRepo{Db: db}
}

//get Courses
func (repository *CourseRepo) GetCourses(c *gin.Context) {
	var Course []models.Courses
	print(Course)
	err := models.GetCourses(repository.Db, &Course)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Course)
}

//get Course by id
func (repository *CourseRepo) GetCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var Course models.Courses
	err := models.GetCourseById(repository.Db, &Course, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Course)
}
